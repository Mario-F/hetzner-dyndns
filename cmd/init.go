package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
)

var (
	token  string
	record string
)

var Version = "development"

func initCommands() {
	flag.BoolVar(&logger.DebugMode, "debug", false, "Turn on debug messages.")
	flag.BoolVar(&logger.QuietMode, "quiet", false, "Suppress logmessages.")
	flag.StringVar(&token, "token", "", "Your Hetzner DNS token.")
	flag.StringVar(&record, "record", "", "Your Hetzner DNS token.")

	flag.Parse()
}

// Execute is the main entrypoint to the application
func Execute() {
	initCommands()

	command := flag.Arg(0)
	logger.Debugf("Executing with token '%v' and command '%v'", token, command)

	switch command {
	case "records":
		logger.Infof("Retrieving all records with id from Hetzer DNS")
		checkToken()
		cmdRecords(token)
	case "myip":
		logger.Infof("Obtain my external IP")
		cmdMyIP()
	case "update":
		logger.Infof("Updating Hetzner DNS with external IP")
		checkToken()
		checkRecord()
		cmdUpdate(token, record)
	case "version":
		fmt.Printf("%s", Version)
	default:
		usage()
	}
}

func checkToken() {
	if token == "" {
		fmt.Println("The token flag is required!")
		usage()
	}
}

func checkRecord() {
	if record == "" {
		fmt.Println("The record id is required!")
		usage()
	}
}

func usage() {
	fmt.Print(`
Updates Hetzner DNS with external IP

Usage:
	hetner-dyndns [command]

Available Commands:
	records		Get record with id from Hetzner DNS
	update		Update record with external IP
	myip		Acquire your external IP and output it
	update		Update your external IP to Hetzner DNS
	version		Print version of hetzner-dyndns

Flags:
`)
	flag.PrintDefaults()
	os.Exit(1)
}
