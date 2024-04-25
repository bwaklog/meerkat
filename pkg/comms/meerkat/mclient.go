package comms

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	pb "meerkat/pkg/comms/meerkat_protocol"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToServer(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	return conn
}

// REPLACE: above connecttoserver function
// JoinNetworkPool will be handling connection to the server
func JoinNetworkPool(node *MeerkatNode, addr string) {
	// TODO
	// Send over a PoolJoinRequest to server with the nodes empty client list
	// and address and port

	// address := fmt.Sprintf("localhost:%s", addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	// Send over a PoolJoinRequest to the server
	c := pb.NewMeerkatGuideClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	poolJoinRequestTemp := &pb.PoolJoinRequest{
		Address: node.Address,
		// Port:       int32(node.Port),
		ClientList: node.Clients,
	}

	// r, err := c.JoinPoolProtocol(ctx, poolJoinRequestTemp)
	// if err != nil {
	// 	log.Fatalf("Not able to send a join pool request to server")
	// }

	// node.Clients = r.GetClientList()

	stream, err := c.JoinPoolProtocol(ctx, poolJoinRequestTemp)
	if err != nil {
		log.Fatalf("Not able to send a join pool request to server")
	}

	// node.mutex.Lock()

	for {
		// continuously recieve data from server
		response, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}

		if err == io.EOF {
			break
		}

		switch response.GetResponse().(type) {
		case *pb.PoolJoinResponse_ClientList:
			node.Clients = response.GetClientList().ClientList

		case *pb.PoolJoinResponse_DirData:
			err := os.MkdirAll(node.NodeData.BaseDir + "/" + response.GetDirData().Path, 0755)
			if err != nil {
				log.Fatalf("Unable to create directory: %v", err)
			}
			log.Println("Created directories for: ", response.GetDirData().Path)

		case *pb.PoolJoinResponse_FileData:
			// TODO
			fileData := response.GetFileData()
			filePath := fileData.GetPath()
			fileBytes := fileData.GetData()

			// write to file system

			file, err := os.Create(node.NodeData.BaseDir + "/" + filePath)
			if err != nil {
				if err == fs.ErrNotExist {
					log.Println("Non existing directory")
				} else {
					log.Fatalf("Not able to write to file system: %v", err)
				}
			}

			file.Write(fileBytes)
			node.NodeData.FileTrack[filePath] = time.Now()
			file.Close()


			// log.Printf("Received %s from node %s", filePath, conn.Target())
			node.NodeData.LoadFileSystem(node.NodeData.BaseDir)
		}

	}
	// node.mutex.Unlock()

	HandshakeClients(node)

	node.ClientsConn = append(node.ClientsConn, conn)
	node.Clients = append(node.Clients, addr)
}

func HandshakeClients(node *MeerkatNode) {
	for _, connAddr := range node.Clients {

		conn, err := grpc.Dial(connAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}

		c := pb.NewMeerkatGuideClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		handshakeRequest := &pb.PoolHandshakesRequest{
			Address: node.Address,
			// Port:       int32(node.Port),
			ClientList: node.Clients,
		}

		log.Printf("Attempting to handshake with clinent %s", conn.Target())
		r, err := c.HandshakePoolProtocol(ctx, handshakeRequest)
		if err != nil {
			log.Fatalf("Couldn't send a handshake with client %s", conn.Target())
		}

		if r.GetSuccess() {
			log.Printf("Handshake with client %s successful", conn.Target())
			node.ClientsConn = append(node.ClientsConn, conn)
		} else {
			log.Printf("Handshake with client %s failed", conn.Target())
		}
	}
}

func HandleDisconnect(node *MeerkatNode) bool {
	log.Println("Clients to disconnect from", node.Clients)
	for _, conn := range node.ClientsConn {
		// send a payload informing the server that the client is disconnecting
		// and close the connection

		log.Println("Disconnecting from pool", conn.Target())

		c := pb.NewMeerkatGuideClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		PoolDisconnectRequestTemp := &pb.PoolDisconnectRequest{
			Address: node.Address,
			// Port:    int32(node.Port),
		}

		r, err := c.DisconnectPoolProtocol(ctx, PoolDisconnectRequestTemp)
		if err != nil {
			if err == grpc.ErrServerStopped {
				log.Println("Server has stopped")
			} else {
				log.Fatalf("Could not send disconnect request to server: %v", err)
			}
		}

		if r.GetSuccess() {
			log.Printf("Successfully disconnected from the pool")
			conn.Close()
		}
	}
	log.Printf("Disconnected from all nodes in the pool")
	return true
}

func (n *MeerkatNode)HandleBroadcastChanges(path string, data []byte) {
	log.Printf("Sending data to : %v", n.Clients)
	for _, conn := range n.ClientsConn {
		c := pb.NewMeerkatGuideClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		stream, err :=c.DataModProtocol(ctx)
		if err != nil {
			log.Fatalf("Not able to send a join pool request to server")
		}

		response := &pb.DataModRequest{
			Response: &pb.DataModRequest_FileData{
				FileData: &pb.FileData{
					Path: path,
					Data: data,
				},
			},
		}

		if err := stream.Send(response); err != nil {
			log.Fatalf("Error from reciever: %v", err)
		}

		reply, err :=stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Error from reciever: %v", err)
		}

		if reply.Success {
			log.Printf("Sent change in %s to %s", path, conn.Target())
		} else {
			log.Fatalf("Failed to send %s change to %s", path, conn.Target())
		}

	}
}

func (n *MeerkatNode) FileTracker() {
	fileSys := n.NodeData.FileSystem
	for {
		n.mutex.Lock()
		fs.WalkDir(fileSys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.Type().IsRegular() && !strings.Contains(path, ".DS_Store"){
				fileInfo, err := d.Info()
				if err != nil {
					log.Fatal(err)
					return err
				}

				// if file is not in the file tracker, add it
				if _, ok := n.NodeData.FileTrack[path]; !ok {
					n.NodeData.FileTrack[path] = fileInfo.ModTime()
					// broadcast the file to all nodes
					buff, err := fs.ReadFile(n.NodeData.FileSystem, path)
					if err != nil {
						log.Fatal(err)
					}

					n.HandleBroadcastChanges(path, buff)
				}

				if modTime, ok := n.NodeData.FileTrack[path]; ok {
					if modTime != fileInfo.ModTime() {

						buff, err := fs.ReadFile(n.NodeData.FileSystem, path)
						if err != nil {
							return err
						}
						n.HandleBroadcastChanges(path, buff)
						n.NodeData.FileTrack[path] = fileInfo.ModTime()

					}
				}
			}

			return nil
		})
		n.mutex.Unlock()
	}
}


func HandleEcho(input string, node *MeerkatNode) {

	for _, conn := range node.ClientsConn {
		c := pb.NewMeerkatGuideClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.EchoText(ctx, &pb.EchoRequest{Message: input})
		if err != nil {
			log.Fatalf("Could not send echo to server: %v", err)
		}

		log.Printf("%s", r.GetMessage())
	}
}
