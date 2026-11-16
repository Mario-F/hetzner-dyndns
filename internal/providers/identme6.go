package providers

import (
	"io/ioutil"
	"net/http"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
)

func identme6GetIP() (string, error) {
	logger.Debugf("Start GetIP with identme6")

	resp, err := http.Get("https://v6.ident.me")
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
	logger.Debugf("Found IP with identme6: %+v\n", ip)
	return ip, nil
}

var identme6Provider Provider = Provider{
	GetIP:        identme6GetIP,
	Version:      network.IPv6,
	ProviderName: "identme6",
}

func init() {
	ProviderList = append(ProviderList, identme6Provider)
}
