package main

import (
	"log"
	comms "meerkat/pkg/comms/meerkat"

	// "meerkat/pkg/comms"

	"github.com/spf13/cobra"
)

func main() {
	var paddr string // port address
	var saddr string // server address to connect to
	var server bool  // used by node to initiate the connections

	// var wg sync.WaitGroup

	rootCmd := &cobra.Command{
		Use:   "meerkat",
		Short: "connection client",
		Run: func(cmd *cobra.Command, args []string) {


			if paddr != "" {
				go comms.StartListener(paddr)

				if !server {
					// this is for connecting to a node who
					// is in the network pool
					if saddr == "" {
						log.Fatalln("Didnt provide a server node to connect to")
					} else {
						comms.ConnectToServer(saddr)
					}
				}
			} else {
				log.Fatalln("Didnt provide a port to open connection")
			}

			for {}


		},
	}

	rootCmd.Flags().StringVarP(&paddr, "port", "p", "", "port address")
	rootCmd.Flags().StringVarP(&saddr, "saddr", "c", "", "server address to connect to")
	rootCmd.Flags().BoolVarP(&server, "server", "s", false, "Used to initiate the connection")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
