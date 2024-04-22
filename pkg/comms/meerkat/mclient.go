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

	address := fmt.Sprintf("localhost:%s", addr)

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	// Send over a PoolJoinRequest to the server
	c := pb.NewMeerkatGuideClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	poolJoinRequestTemp := &pb.PoolJoinRequest{
		Address:    node.Address,
		Port:       int32(node.Port),
		ClientList: node.Clients,
	}

	r, err := c.JoinPoolProtocol(ctx, poolJoinRequestTemp)
	if err != nil {
		log.Fatalf("Not able to send a join pool request to server")
	}

	node.Clients = r.GetClientList()
	node.Data = r.GetData()

	if len(node.Clients) > 0 {
		HandshakeClients(node)
	}

	node.ClientsConn = append(node.ClientsConn, conn)
	node.Clients = append(node.Clients, address)
}

func HandshakeClients(node *MeerkatNode) {
	for _, conn := range node.ClientsConn {
		c := pb.NewMeerkatGuideClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		handshakeRequest := &pb.PoolHandshakesRequest{
			Address:    node.Address,
			Port:       int32(node.Port),
			ClientList: node.Clients,
		}

		r, err := c.HandshakePoolProtocol(ctx, handshakeRequest)
		if err != nil {
			log.Fatalf("Couldn't send a handshake with client %s", conn.Target())
		}

		if r.GetSuccess() {
			log.Printf("Handshake with client %s successful", conn.Target())
			node.Clients = append(node.Clients, conn.Target())
			node.ClientsConn = append(node.ClientsConn, conn)
		} else {
			log.Printf("Handshake with client %s failed", conn.Target())
		}
	}
}

func HandleDisconnect(node *MeerkatNode) bool {
	for _, conn := range node.ClientsConn {
		// send a payload informing the server that the client is disconnecting
		// and close the connection

		c := pb.NewMeerkatGuideClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		PoolDisconnectRequestTemp := &pb.PoolDisconnectRequest{
			Address: node.Address,
			Port:    int32(node.Port),
		}

		r, err := c.DisconnectPoolProtocol(ctx, PoolDisconnectRequestTemp)
		if err != nil {
			log.Fatalf("Could not send disconnect request to server: %v", err)
		}

		if r.GetSuccess() {
			log.Printf("Successfully disconnected from the pool")
			conn.Close()
			return true
		}
	}
	log.Printf("Failed to disconnect from the pool")
	return false
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
