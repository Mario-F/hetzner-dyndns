package providers

import (
	"io/ioutil"
	"net/http"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
)

func whatismyipGetIP() (string, error) {
	logger.Debugf("Start GetIP with whatismyip")

	resp, err := http.Get("https://www.whatismyip-address.com/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errResponseNotOK
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip, err := captureIPv4(string(body))
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", errIPNotFound
	}
	logger.Debugf("Found IP with whatismyip: %+v\n", ip)
	return ip, nil
}

var whatismyipProvider Provider = Provider{
	GetIP:        whatismyipGetIP,
	Version:      network.IPv4,
	ProviderName: "whatismyip",
}

func init() {
	ProviderList = append(ProviderList, whatismyipProvider)
}
