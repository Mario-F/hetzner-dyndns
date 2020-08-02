package logger

import (
	"log"
)

var (
	// DebugMode logs debug info to console
	DebugMode bool = false
	// QuietMode surpresses all log output except debug
	QuietMode bool = false
)

// Debugf logs formatted debug messages
func Debugf(msg string, vars ...interface{}) {
	if DebugMode {
		log.Printf(msg, vars...)
	}
}

// Infof logs formatted info messages
func Infof(msg string, vars ...interface{}) {
	if !QuietMode {
		log.Printf(msg, vars...)
	}
}
