package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Mario-F/hetzner-dyndns/internal/hetzner"
	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/spf13/cobra"
)

var domain string

var recordsCmd = &cobra.Command{
	Use:   "records",
	Short: "Retrieving all records with id from Hetzner DNS",
	Run: func(cmd *cobra.Command, args []string) {
		var resRecords []hetzner.Record

		hetzner.SetToken(token)
		resRecords = hetzner.GetRecords()
		logger.Debugf("Results from GetRecords: %d", len(resRecords))

		// Process result to find order and size data
		var (
			longestDomain int            = 10
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
			domainMap[v.Fullname] = i
			domainSlice = append(domainSlice, v.Fullname)
		}
		sort.Strings(domainSlice)

		if len(domainSlice) == 0 {
			fmt.Println("No records found")
			return
		}

		// Create an human readable result table
		domainCellLength := longestDomain + 3
		fmt.Println()
		printRecordLine("Domain", "ID", domainCellLength)
		for _, k := range domainSlice {
			record := resRecords[domainMap[k]]
			printRecordLine(record.Fullname, record.ID, domainCellLength)
		}
		fmt.Println()
	},
}

func printRecordLine(lineFirst string, lineSecond string, lengthFirst int) {
	fillSpace := lengthFirst - len(lineFirst)
	fmt.Println(lineFirst, strings.Repeat(" ", fillSpace), lineSecond)
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
