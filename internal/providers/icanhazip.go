package providers

import (
	"io/ioutil"
	"net/http"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
)

func icanhazipGetIP() (string, error) {
	logger.Debugf("Start GetIP with icanhazip")

	resp, err := http.Get("https://ipv6.icanhazip.com")
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

	ip, err := captureIPv6(string(body))
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", errIPNotFound
	}
	logger.Debugf("Found IP with icanhazip: %+v\n", ip)
	return ip, nil
}

var icanhazipProvider Provider = Provider{
	GetIP:        icanhazipGetIP,
	Version:      network.IPv6,
	ProviderName: "icanhazip",
}

func init() {
	ProviderList = append(ProviderList, icanhazipProvider)
}
