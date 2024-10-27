package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/b31ngd3v/cat-forwarding/internal/client"
	"github.com/spf13/cobra"
)

const (
	shortDesc = "Open the purr-tals for cats to connect worldwide!"
	longDesc  = `Cat Forwarding helps expose a port on your local machine to the outside world,
like opening a secret tunnel for cats to chat globally. Just provide a valid 
port (1-65535) to let your cat's meows travel beyond the home network.`
)

func main() {
	var port int

	var rootCmd = &cobra.Command{
		Use:   "cat-forwarding [port]",
		Short: shortDesc,
		Long:  longDesc,
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			p, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid port number: %s", args[0])
			}
			if p <= 0 || p > 65535 {
				return fmt.Errorf("port number must be between 1 and 65535")
			}
			port = p
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			client.Run(port)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
