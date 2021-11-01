package providers

import (
	"errors"
	"regexp"
)

var errIPNotFound error = errors.New("IP Not Found")
var errResponseNotOK error = errors.New("HTTP Response is not OK")

// Provider holds a external ip provider
type Provider struct {
	GetIP        func() (string, error)
	ProviderName string
}

// ProviderList has all created providers
var ProviderList []Provider = []Provider{
	checkIPProvider,
	ifconfigMEProvider,
	ipifyProvider,
	whatismyipProvider,
	whoismyispProvider,
}

func captureIP(text string) (string, error) {
	r, err := regexp.Compile(`(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)
	if err != nil {
		return "", err
	}
	return r.FindString(text), nil
}
