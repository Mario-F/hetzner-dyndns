package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of hetzner-dyndns",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
