package providers

import (
	"io/ioutil"
	"net/http"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
)

func ipifyGetIP() (string, error) {
	logger.Debugf("Start GetIP with ipify")

	resp, err := http.Get("https://api.ipify.org")
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
	logger.Debugf("Found IP with ipify: %+v\n", ip)
	return ip, nil
}

var ipifyProvider Provider = Provider{
	GetIP:        ipifyGetIP,
	Version:      network.IPv4,
	ProviderName: "Ipify",
}

func init() {
	ProviderList = append(ProviderList, ipifyProvider)
}
