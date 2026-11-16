package providers

import (
	"io/ioutil"
	"net/http"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
)

func icanhazip4GetIP() (string, error) {
	logger.Debugf("Start GetIP with icanhazip4")

	resp, err := http.Get("https://ipv4.icanhazip.com")
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
	logger.Debugf("Found IP with icanhazip4: %+v\n", ip)
	return ip, nil
}

var icanhazip4Provider Provider = Provider{
	GetIP:        icanhazip4GetIP,
	Version:      network.IPv4,
	ProviderName: "icanhazip4",
}

func init() {
	ProviderList = append(ProviderList, icanhazip4Provider)
}
