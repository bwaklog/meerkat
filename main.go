package main

import (
	"fmt"
	"log"
	comms "meerkat/pkg/comms/meerkat"
	"meerkat/pkg/data"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	var paddr string    // port address
	var saddr string    // server address to connect to
	var server bool     // used by node to initiate the connections
	var dirwatch string // the path to the directory to watch

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// var wg sync.WaitGroup

	// if dirwatch == "" {
	// 	log.Fatalln("Didn't provide a directory to watch")
	// }

	// load file system and fill the NodeData
	// with the file system data

	rootCmd := &cobra.Command{
		Use:   "meerkat",
		Short: "connection client",
		Run: func(cmd *cobra.Command, args []string) {

			// var conn *grpc.ClientConn
			if dirwatch == "" {
				log.Fatalln("Didn't provide a directory to watch")
			}

			node := comms.MeerkatNode{
				Clients: []string{},
				NodeData: data.NodeData{
					FileSystem: os.DirFS(dirwatch),
					BaseDir:    dirwatch,
					FileTrackMap: data.FileTrackMap{
						FileTrack: make(map[string]time.Time),
					},
				},
			}

			node.NodeData.LoadFileSystem(dirwatch)
			log.Println("Loaded file system: ", dirwatch)

			// print all files available
			// for k, v := range node.NodeData.FileTrack {
			// 	log.Printf("File: %s, ModTime: %s", k, v)
			// }

			if paddr != "" {

				node.Address = paddr
				// node.Port, _ = strconv.Atoi(paddr)

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

			// run the evaluation throuhg in the background
			// go func() {
			// 	for {
			// 		log.Printf("Evaluating %d clients", len(node.Clients))
			// 		for _, conn := range node.Clients {
			// 			log.Println("🖥️ CLIENT: ", conn)
			// 		}
			// 		time.Sleep(5 * time.Second)
			// 	}
			// }()

			// go node.FileTracker()

			// graceful exit
			go func() {
				sigChan := make(chan os.Signal, 1)
				signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
				<-sigChan

				log.Println("Shutting down...")
				// exit
				comms.HandleDisconnect(&node)
				os.Exit(0)
			}()

			for {
				var input string
				// fmt.Print("MSG: ")
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
	rootCmd.Flags().StringVarP(&dirwatch, "dir", "d", "", "The directory to watch/share among clients")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
