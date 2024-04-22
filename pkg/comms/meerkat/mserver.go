package comms

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "meerkat/pkg/comms/meerkat_protocol"

	"google.golang.org/grpc"
)

type MeerkatServer struct {
	pb.UnimplementedMeerkatGuideServer

	mu sync.Mutex
}

// reply back to client wiht :: "Hello from Meerkat Server"
func (s *MeerkatServer) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Recieved: %s", request.GetName())
	return &pb.HelloResponse{Message: "Hello from Meerkat Server " + request.GetName()}, nil
}

func StartListener(addr string) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", addr))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterMeerkatGuideServer(s, &MeerkatServer{})
	log.Printf("Server listening at localhost:%s", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
