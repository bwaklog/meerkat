package comms

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "meerkat/pkg/comms/meerkat_protocol"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToServer(addr string) {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewMeerkatGuideClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Meerkat Client"})
	// if err != nil {
	// 	log.Fatalf("Could not greet: %v", err)
	// }
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Meerkat Client"})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}

	log.Printf("Server: %s", r.GetMessage())
}
