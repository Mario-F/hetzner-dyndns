package providers

import (
	"io/ioutil"
	"net/http"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
)

func piliappGetIP() (string, error) {
	logger.Debugf("Start GetIP with piliapp")

	resp, err := http.Get("https://de.piliapp.com/what-is-my/ipv6/")
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
	logger.Debugf("Found IP with piliapp: %+v\n", ip)
	return ip, nil
}

var piliappProvider Provider = Provider{
	GetIP:        piliappGetIP,
	Version:      network.IPv6,
	ProviderName: "piliapp",
}

func init() {
	ProviderList = append(ProviderList, piliappProvider)
}
