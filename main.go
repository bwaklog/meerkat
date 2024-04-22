package main

import (
	"fmt"
	"log"
	comms "meerkat/pkg/comms/meerkat"
	"strconv"
	"time"

	// "meerkat/pkg/comms"

	"github.com/spf13/cobra"
)

func main() {
	var paddr string // port address
	var saddr string // server address to connect to
	var server bool  // used by node to initiate the connections

	// var wg sync.WaitGroup

	node := comms.MeerkatNode{
		Clients: []string{},
		Data:    "",
	}

	rootCmd := &cobra.Command{
		Use:   "meerkat",
		Short: "connection client",
		Run: func(cmd *cobra.Command, args []string) {

			// var conn *grpc.ClientConn

			if paddr != "" {

				node.Address = "localhost"
				node.Port, _ = strconv.Atoi(paddr)

				// this actively listens in a separate go routine
				// for any incoming connections
				go comms.StartListener(paddr, &node)

				if !server {
					// this is for connecting to a node who
					// is in the network pool
					if saddr == "" {
						log.Fatalln("Didnt provide a server node to connect to")
					} else {
						// conn = comms.ConnectToServer(saddr)
						// this is where we joing the network pool
						comms.JoinNetworkPool(&node, saddr)
						// function adds the client to its client list
					}
				}

			} else {
				log.Fatalln("Didnt provide a port to open connection")
			}

			go func() {
				log.Printf("Evaluating %d clients", len(node.Clients))
				for _, conn := range node.Clients {
					log.Println("🖥️ CLIENT: ", conn)
				}
				time.Sleep(5 * time.Second)
			}()

			for {
				var input string
				fmt.Print("MSG: ")
				fmt.Scanln(&input)

				if input == "exit" {
					if comms.HandleDisconnect(&node) {
						break
					} else {
						log.Println("Failed to disconnect")
					}
				}
			}

		},
	}

	rootCmd.Flags().StringVarP(&paddr, "port", "p", "", "port address")
	rootCmd.Flags().StringVarP(&saddr, "saddr", "c", "", "server address to connect to")
	rootCmd.Flags().BoolVarP(&server, "server", "s", false, "Used to initiate the connection")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
