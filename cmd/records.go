package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Mario-F/hetzner-dyndns/internal/hetzner"
	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
	"github.com/spf13/cobra"
)

var domain string

var recordsCmd = &cobra.Command{
	Use:   "records",
	Short: "Retrieving all records with id from Hetzner DNS",
	Run: func(cmd *cobra.Command, args []string) {
		var resRecords []hetzner.Record

		hetzner.SetToken(token)
		resRecords = hetzner.GetRecords(network.IPVersion(ipVersion))
		logger.Debugf("Results from GetRecords: %d", len(resRecords))

		// Process result to find order and size data
		var (
			longestDomain int            = 8
			longestType   int            = 4
			longestValue  int            = 8
			domainSlice   []string       = make([]string, 0, len(resRecords))
			domainMap     map[string]int = make(map[string]int, len(resRecords))
		)
		for i, v := range resRecords {
			if domain != "" {
				if v.Fullname != domain {
					continue
				}
			}
			if len(v.Fullname) > longestDomain {
				longestDomain = len(v.Fullname)
			}
			if len(v.Type) > longestType {
				longestType = len(v.Type)
			}
			if len(v.Value) > longestValue {
				longestValue = len(v.Value)
			}
			domainMap[v.Fullname] = i
			domainSlice = append(domainSlice, v.Fullname)
		}
		sort.Strings(domainSlice)

		if len(domainSlice) == 0 {
			fmt.Println("No records found")
			return
		}

		fmt.Println()
		printRecordLine("Domain", "Type", "Value", "ID", longestDomain, longestType, longestValue)
		for _, k := range domainSlice {
			record := resRecords[domainMap[k]]
			printRecordLine(record.Fullname, record.Type, record.Value, record.ID, longestDomain, longestType, longestValue)
		}
		fmt.Println()
	},
}

func printRecordLine(rDomain string, rType string, rValue string, rID string, longestDomain int, longestType int, longestValue int) {
	fillDomain := longestDomain - len(rDomain)
	fillType := longestType - len(rType)
	fillValue := longestValue - len(rValue)
	fmt.Println(rDomain, strings.Repeat(" ", fillDomain), rType, strings.Repeat(" ", fillType), rValue, strings.Repeat(" ", fillValue), rID)
}

func init() {
	rootCmd.AddCommand(recordsCmd)

	recordsCmd.Flags().StringVar(&token, "token", "", "The hetzner token to access DNS API")
	recordsCmd.Flags().StringVar(&domain, "domain", "", "Get exactly this domain")
	err := recordsCmd.MarkFlagRequired("token")
	if err != nil {
		fmt.Printf("An error occurred: %v+", err)
	}
}
