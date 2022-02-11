package cmd

import (
	"fmt"
	"sync"

	"github.com/Mario-F/hetzner-dyndns/internal/externalip"
	"github.com/Mario-F/hetzner-dyndns/internal/hetzner"
	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update record with external IP",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			wg     sync.WaitGroup = sync.WaitGroup{}
			record hetzner.Record
			myip   externalip.ExternalIP
		)
		hetzner.SetToken(token)

		logger.Infof("Requesting external IP and Hetzner DNS Record")
		wg.Add(2)
		go func() {
			record = hetzner.GetRecord(recordInput)
			wg.Done()
		}()
		go func() {
			var err error
			myip, err = externalip.GetExternalIP(externalip.IPv4)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}()
		wg.Wait()

		logger.Infof("Domain: %+v", record.Fullname)
		logger.Infof("Old - IP: %+v", record.Value)
		logger.Infof("New - IP: %+v", myip.IP)

		if record.Value == myip.IP {
			logger.Infof("No IP diff detected, update not necessary.")
			return
		}

		record.Value = myip.IP
		err := record.Update()
		if err != nil {
			panic(err)
		}
		logger.Infof("Record updated successfully")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&token, "token", "", "The hetzner token to access DNS API")
	err := updateCmd.MarkFlagRequired("token")
	if err != nil {
		fmt.Printf("An error occurred: %v+", err)
	}
	updateCmd.Flags().StringVar(&recordInput, "record", "", "The DNS record to update with the external IP")
	err = updateCmd.MarkFlagRequired("record")
	if err != nil {
		fmt.Printf("An error occurred: %v+", err)
	}
}
