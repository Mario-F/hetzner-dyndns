package providers

import (
	"errors"
	"regexp"

	"github.com/Mario-F/hetzner-dyndns/internal/network"
)

var errIPNotFound error = errors.New("IP Not Found")
var errResponseNotOK error = errors.New("HTTP Response is not OK")

// Provider holds a external ip provider
type Provider struct {
	GetIP        func() (string, error)
	Version      network.IPVersion
	ProviderName string
}

// ProviderList has all created providers
var ProviderList []Provider = []Provider{}

func captureIPv4(text string) (string, error) {
	r, err := regexp.Compile(`(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)
	if err != nil {
		return "", err
	}
	return r.FindString(text), nil
}
