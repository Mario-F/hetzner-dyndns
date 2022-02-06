package cmd

import (
	"fmt"

	"github.com/Mario-F/hetzner-dyndns/internal/externalip"
	"github.com/spf13/cobra"
)

var myipCmd = &cobra.Command{
	Use:   "myip",
	Short: "Acquire your external IP and output it",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := externalip.GetExternalIP()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Your external IP is: %v\n", res)
	},
}

func init() {
	rootCmd.AddCommand(myipCmd)
	// TODO: Add flag for output formats
}
