package comms

import "google.golang.org/grpc"

type MeerkatNode struct {
	Clients     []string
	ClientsConn []*grpc.ClientConn
	Data        string
	Address     string
	Port        int
}
