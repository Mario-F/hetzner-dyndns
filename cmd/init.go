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

func initCommands() {
	flag.BoolVar(&logger.DebugMode, "debug", false, "Turn on debug messages.")
	flag.BoolVar(&logger.QuietMode, "quiet", false, "Surpress logmessages.")
	flag.StringVar(&token, "token", "", "Your Hetzner DNS token.")
	flag.StringVar(&record, "record", "", "Your Hetzner DNS token.")

	flag.Parse()
}

// Execute is the main entrypoint to the application
func Execute() {
	initCommands()

	command := flag.Arg(0)
	logger.Debugf("Executing with token '%v' and command '%v'", token, command)

	// Token is required
	if token == "" {
		fmt.Println("The token flag is required!")
		usage()
	}

	switch command {
	case "records":
		logger.Infof("Retrieving all records with id from Hetzer DNS")
		cmdRecords(token)
	case "update":
		logger.Infof("Updating Hetzner DNS with external IP")
	default:
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

Flags:
`)
	flag.PrintDefaults()
	os.Exit(1)
}
