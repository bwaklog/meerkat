package comms

import (
	"context"
	"io"
	"io/fs"
	"log"
	"net"
	"os"

	pb "meerkat/pkg/comms/meerkat_protocol"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MeerkatServer struct {
	pb.UnimplementedMeerkatGuideServer

	// mu sync.Mutex
	Node *MeerkatNode
}

// reply back to client wiht :: "Hello from Meerkat Server"
func (s *MeerkatServer) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Recieved: %s", request.GetName())
	return &pb.HelloResponse{Message: "Hello from Meerkat Server " + request.GetName()}, nil
}

func (s *MeerkatServer) EchoText(ctx context.Context, request *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Printf("Recieved: %s", request.GetMessage())
	return &pb.EchoResponse{Message: "ECHO: " + request.GetMessage()}, nil
}

func (s *MeerkatServer) JoinPoolProtocol(request *pb.PoolJoinRequest, stream pb.MeerkatGuide_JoinPoolProtocolServer) error {
	log.Printf("Node %s attempting to connect to pool", request.GetAddress())

	// exchanging credentials
	NodeClientList := s.Node.Clients

	s.Node.Clients = append(s.Node.Clients, request.GetAddress())

	nodeConn := DialNewNode(request.GetAddress())

	s.Node.ClientsConn = append(s.Node.ClientsConn, nodeConn)

	response := &pb.PoolJoinResponse{
		Response: &pb.PoolJoinResponse_ClientList{
			ClientList: &pb.ClientList{
				ClientList: NodeClientList,
			},
		},
	}

	if err := stream.Send(response); err != nil {
		return err
	}

	fileSystem := s.Node.NodeData.FileSystem

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {

		if d.IsDir() && path != "." {
			response = &pb.PoolJoinResponse{
				Response: &pb.PoolJoinResponse_DirData{
					DirData: &pb.DirData{
						Path: path,
					},
				},
			}

			if err = stream.Send(response); err != nil {
				return err
			}

		}

		if d.Type().IsRegular() {
			buff, err := fs.ReadFile(s.Node.NodeData.FileSystem, path)
			if err != nil {
				return err
			}

			response = &pb.PoolJoinResponse{
				Response: &pb.PoolJoinResponse_FileData{
					FileData: &pb.FileData{
						Path: path,
						Data: buff,
					},
				},
			}

			if err = stream.Send(response); err != nil {
				return err
			}

		}

		return nil
	})

	// send file to client with meta data such as path

	// log.Printf("Sent %s from node %s",, request.GetAddress())

	return nil

}

func (s *MeerkatServer) HandshakePoolProtocol(ctx context.Context, request *pb.PoolHandshakesRequest) (*pb.PoolHandshakeResponse, error) {
	// this is when a node already connected to some other node in the
	// network pool is initiating a connection between nodes part
	// of the client list
	// address := fmt.Sprintf("%s:%d", request.GetAddress(), request.GetPort())

	log.Printf("Handshake with node %s", request.GetAddress())

	// adds node to the list of clients
	s.Node.Clients = append(s.Node.Clients, request.GetAddress())
	conn := DialNewNode(request.GetAddress())
	s.Node.ClientsConn = append(s.Node.ClientsConn, conn)

	return &pb.PoolHandshakeResponse{Success: true}, nil
}

// func (s *MeerkatServer) DataModProtocol(stream pb.Meerk)
func (s *MeerkatServer) DataModProtocol(stream pb.MeerkatGuide_DataModProtocolServer) error {
	log.Printf("Initiating data mod protocol")

	s.Node.mutex.Lock()

	for {
		request, err := stream.Recv()

		if err == io.EOF {
			log.Printf("Recieved EOF")
			response := &pb.DataModResponse{
				Success: true,
			}

			// s.Node.NodeData.LoadFileSystem(s.Node.NodeData.BaseDir)
			s.Node.NodeData.DirSnapshot()

			log.Println("Unlocking the node mutex")
			s.Node.mutex.Unlock()
			return stream.SendAndClose(response)
		}

		if err != nil {
			log.Printf("Error in request: %v", err)
			return err
		}

		switch request.Response.(type) {

		case *pb.DataModRequest_DirData:
			log.Printf("Recieving data from peer: %s", request.GetFileData().Path)
			dirData := request.GetDirData()
			// if the directory doesnt exist, create it
			if _, err := fs.Stat(s.Node.NodeData.FileSystem, dirData.Path); err != nil {
				err := os.MkdirAll(dirData.Path, 0755)
				if err != nil {
					log.Fatalf("Failed to create directory: %v", err)
				}
			}

		case *pb.DataModRequest_FileData:

			log.Printf("Mod request : %v", request.GetFileData().Action)

			switch request.GetFileData().Action {
			case pb.Action_ADD:
				log.Println("Add/Update file operation initiated")
				fileData := request.GetFileData()
				filePath := fileData.GetPath()
				fileBytes := fileData.GetData()

				file, err := os.Create(s.Node.NodeData.BaseDir + "/" + filePath)
				if err != nil {
					if err == fs.ErrNotExist {
						log.Println("Not existing directory")
					} else {
						log.Fatalf("Not able to write to file system: %v", err)
					}
				}

				file.Write(fileBytes)

				s.Node.NodeData.DiskSnapshot.Lock.Lock()

				// create a fs.DirEntry of the received file
				// fill in the old stats from the current file

				fileDirInfo, err := fs.Stat(s.Node.NodeData.FileSystem, filePath)
				if err != nil {
					log.Fatalf("Error in getting file info: %v", err)
				}

				// update the file in the file system
				s.Node.NodeData.DiskSnapshot.File[filePath] = fileDirInfo

				// s.Node.NodeData.FileTrackMap.Lock.Unlock()
				s.Node.NodeData.DiskSnapshot.Lock.Unlock()

				log.Println("Updated file written to disk")

				file.Close()

			case pb.Action_DELETE:
				fileData := request.GetFileData()
				filePath := fileData.GetPath()

				os.Remove(s.Node.NodeData.BaseDir + "/" + filePath)
				if err != nil {
					log.Fatal(err)
					log.Fatalf("Unable to delete file %s", filePath)
				}
				log.Printf("Deleted file %s", filePath)

				s.Node.NodeData.DiskSnapshot.Lock.Lock()
				delete(s.Node.NodeData.DiskSnapshot.File, filePath)
				s.Node.NodeData.DiskSnapshot.Lock.Unlock()

				log.Printf("Deleted %s from disk snapshot", filePath)
			}

		}
	}
}

// func (s *MeerkatServer) DisconnectPoolProtocol(ctx context.Context, request )
func (s *MeerkatServer) DisconnectPoolProtocol(ctx context.Context, request *pb.PoolDisconnectRequest) (*pb.PoolDisconnectResponse, error) {
	// address := fmt.Sprintf("%s:%d", request.GetAddress(), request.GetPort())
	log.Printf("Node %s is disconnecting ", request.GetAddress())

	// remove the node from the list of clients
	s.Node.mutex.Lock()
	for i, client := range s.Node.Clients {
		if client == request.GetAddress() {
			s.Node.Clients = append(s.Node.Clients[:i], s.Node.Clients[i+1:]...)
			break
		}
	}

	// close the connection
	for i, conn := range s.Node.ClientsConn {
		if conn.Target() == request.GetAddress() {
			s.Node.ClientsConn = append(s.Node.ClientsConn[:i], s.Node.ClientsConn[i+1:]...)
			conn.Close()
			break
		}
	}
	s.Node.mutex.Unlock()

	return &pb.PoolDisconnectResponse{Success: true}, nil
}

// dial a connection to a node entering a pool
func DialNewNode(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	log.Printf("Dialed %s and connected", addr)

	return conn
}

func StartListener(addr string, mkNode *MeerkatNode) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	mkServer := &MeerkatServer{}
	mkServer.Node = mkNode

	pb.RegisterMeerkatGuideServer(s, mkServer)

	log.Printf("Server listening at %s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
