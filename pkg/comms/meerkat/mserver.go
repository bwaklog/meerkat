package comms

import (
	"context"
	"fmt"
	"log"
	"net"

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

func (s *MeerkatServer) JoinPoolProtocol(ctx context.Context, request *pb.PoolJoinRequest) (*pb.PoolJoinResponse, error) {
	// log.Printf("Node %s:%d attempting to connect to pool", request.GetAddress(), request.GetPort())
	address := fmt.Sprintf("%s:%d", request.GetAddress(), request.GetPort())

	// exchanging credentials
	NodeClientList := s.Node.Clients

	s.Node.Clients = append(s.Node.Clients, address)

	nodeConn := DialNewNode(address)

	s.Node.ClientsConn = append(s.Node.ClientsConn, nodeConn)

	// for _, connAddr := range s.Node.Clients {
	// 	fmt.Println("CLIENT: ", connAddr)
	// }

	return &pb.PoolJoinResponse{ClientList: NodeClientList, Data: s.Node.Data}, nil
}

func (s *MeerkatServer) HandshakePoolProtocol(ctx context.Context, request *pb.PoolHandshakesRequest) (*pb.PoolHandshakeResponse, error) {
	// this is when a node already connected to some other node in the
	// network pool is initiating a connection between nodes part
	// of the client list
	address := fmt.Sprintf("%s:%d", request.GetAddress(), request.GetPort())

	log.Printf("Handshake with node %s", address)

	// adds node to the list of clients
	s.Node.Clients = append(s.Node.Clients, address)
	conn := DialNewNode(address)
	s.Node.ClientsConn = append(s.Node.ClientsConn, conn)

	return &pb.PoolHandshakeResponse{Success: true}, nil
}

// func (s *MeerkatServer) DisconnectPoolProtocol(ctx context.Context, request )
func (s *MeerkatServer) DisconnectPoolProtocol(ctx context.Context, request *pb.PoolDisconnectRequest) (*pb.PoolDisconnectResponse, error) {
	address := fmt.Sprintf("%s:%d", request.GetAddress(), request.GetPort())
	log.Printf("Node %s is disconnecting ", address)

	// remove the node from the list of clients
	for i, client := range s.Node.Clients {
		if client == address {
			s.Node.Clients = append(s.Node.Clients[:i], s.Node.Clients[i+1:]...)
			break
		}
	}

	// close the connection
	for i, conn := range s.Node.ClientsConn {
		if conn.Target() == address {
			s.Node.ClientsConn = append(s.Node.ClientsConn[:i], s.Node.ClientsConn[i+1:]...)
			conn.Close()
			break
		}
	}

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
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", addr))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	mkServer := &MeerkatServer{}
	mkServer.Node = mkNode

	pb.RegisterMeerkatGuideServer(s, mkServer)

	log.Printf("Server listening at localhost:%s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
