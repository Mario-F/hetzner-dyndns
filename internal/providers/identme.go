package providers

import (
	"io/ioutil"
	"net/http"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
)

func identme4GetIP() (string, error) {
	logger.Debugf("Start GetIP with identme4")

	resp, err := http.Get("https://v4.ident.me")
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
	logger.Debugf("Found IP with identme4: %+v\n", ip)
	return ip, nil
}

var identme4Provider Provider = Provider{
	GetIP:        identme4GetIP,
	Version:      network.IPv4,
	ProviderName: "identme4",
}

func init() {
	ProviderList = append(ProviderList, identme4Provider)
}
