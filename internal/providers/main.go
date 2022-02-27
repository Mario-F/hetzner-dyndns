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

func captureIPv6(text string) (string, error) {
	// used modified regex from: https://gist.github.com/syzdek/6086792
	r, err := regexp.Compile(`(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`)
	if err != nil {
		return "", err
	}
	return r.FindString(text), nil
}
