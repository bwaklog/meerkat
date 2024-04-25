package comms

import (
	"meerkat/pkg/data"
	"sync"

	"google.golang.org/grpc"
)

type MeerkatNode struct {
	Clients     []string
	ClientsConn []*grpc.ClientConn
	Address     string
	// Port        int
	NodeData data.NodeData

	mutex sync.Mutex
}
