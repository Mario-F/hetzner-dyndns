package cmd

import (
	"sync"

	"github.com/Mario-F/hetzner-dyndns/internal/externalip"
	"github.com/Mario-F/hetzner-dyndns/internal/hetzner"
	"github.com/Mario-F/hetzner-dyndns/internal/logger"
)

func cmdUpdate(token string, recordID string) {
	var (
		wg     sync.WaitGroup = sync.WaitGroup{}
		record hetzner.Record
		myip   string
	)
	hetzner.SetToken(token)

	logger.Infof("Requesting external IP and Hetzner DNS Record")
	wg.Add(2)
	go func() {
		record = hetzner.GetRecord(recordID)
		wg.Done()
	}()
	go func() {
		var err error
		myip, err = externalip.GetExternalIP()
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()
	wg.Wait()

	logger.Infof("Domain: %+v", record.Fullname)
	logger.Infof("Old - IP: %+v", record.Value)
	logger.Infof("New - IP: %+v", myip)

	if record.Value == myip {
		logger.Infof("No IP diff detected, update not necessary.")
		return
	}

	record.Value = myip
	err := record.Update()
	if err != nil {
		panic(err)
	}
	logger.Infof("Record updated successfully")
}
