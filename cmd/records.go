package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"

	"github.com/Mario-F/hetzner-dyndns/internal/hetzner"
)

func cmdRecords(token string) {
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
		if len(v.Fullname) > longestDomain {
			longestDomain = len(v.Fullname)
		}
		domainMap[v.Fullname] = i
		domainSlice = append(domainSlice, v.Fullname)
	}
	sort.Strings(domainSlice)

	// Create an human readable result table
	domainCellLength := longestDomain + 3
	fmt.Println()
	printRecordLine("Domain", "ID", domainCellLength)
	for _, k := range domainSlice {
		record := resRecords[domainMap[k]]
		printRecordLine(record.Fullname, record.ID, domainCellLength)
	}
	fmt.Println()
}

func printRecordLine(lineFirst string, lineSecond string, lengthFirst int) {
	fillSpace := lengthFirst - len(lineFirst)
	fmt.Println(lineFirst, strings.Repeat(" ", fillSpace), lineSecond)
}
