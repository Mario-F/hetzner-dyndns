package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Mario-F/hetzner-dyndns/internal/externalip"
	"github.com/spf13/cobra"
)

var OutputModes = []string{
	"text",
	"short",
}

var outputMode string

var myipCmd = &cobra.Command{
	Use:   "myip",
	Short: "Acquire your external IP and output it",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := externalip.GetExternalIP()
		if err != nil {
			panic(err)
		}

		switch outputMode {
		case "text":
			fmt.Printf("Your external IP is: %s\n", res)
		case "short":
			fmt.Printf("%s", res)
		default:
			err := fmt.Errorf("Output mode is not valid, use one of: %s", strings.Join(OutputModes, ", "))
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(myipCmd)

	myipCmd.Flags().StringVarP(&outputMode, "output", "o", "text", fmt.Sprintf("How the result should be formatted.\nAllowed values: %s", strings.Join(OutputModes, ", ")))
}
