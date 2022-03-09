package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Mario-F/hetzner-dyndns/internal/externalip"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
	"github.com/spf13/cobra"
)

var OutputModes = []string{
	"text",
	"short",
	"json",
}

var (
	outputMode string
	ipVersion  string
)

var myipCmd = &cobra.Command{
	Use:   "myip",
	Short: "Acquire your external IP and output it",
	Run: func(cmd *cobra.Command, args []string) {
		var res externalip.ExternalIP
		var err error

		switch ipVersion {
		case string(network.IPv4):
			res, err = externalip.GetExternalIP(network.IPv4)
		case string(network.IPv6):
			res, err = externalip.GetExternalIP(network.IPv6)
		}
		if err != nil {
			panic(err)
		}

		switch outputMode {
		case "text":
			fmt.Printf("Your external IP is: %s\n", res.IP)
		case "short":
			fmt.Printf("%s", res.IP)
		case "json":
			jsonRes, err := json.Marshal(res)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s", jsonRes)
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
