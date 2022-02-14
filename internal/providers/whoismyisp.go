package providers

import (
	"io/ioutil"
	"net/http"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
)

func whoismyispGetIP() (string, error) {
	logger.Debugf("Start GetIP with whoismyisp")

	resp, err := http.Get("https://www.whoismyisp.org/")
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
	logger.Debugf("Found IP wihth whoismyisp: %+v\n", ip)
	return ip, nil
}

var whoismyispProvider Provider = Provider{
	GetIP:        whoismyispGetIP,
	Version:      network.IPv4,
	ProviderName: "whoismyisp",
}

func init() {
	ProviderList = append(ProviderList, whoismyispProvider)
}
