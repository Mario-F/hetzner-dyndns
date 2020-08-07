package providers

import (
	"errors"
	"regexp"
)

var errIPNotFound error = errors.New("IP Not Found")

// Provider holds a external ip provider
type Provider struct {
	GetIP        func() (string, error)
	ProviderName string
}

// ProviderList has all created providers
var ProviderList []Provider = []Provider{
	testProvider,
	ipifyProvider,
}

var testProvider Provider = Provider{
	GetIP: func() (string, error) {
		return "TestIP", nil
	},
	ProviderName: "TestProvider",
}

func captureIP(text string) (string, error) {
	r, err := regexp.Compile(`(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)
	if err != nil {
		return "", err
	}
	return r.FindString(text), nil
}
