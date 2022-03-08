package providers

import (
	"io/ioutil"
	"net/http"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
)

func ip6meGetIP() (string, error) {
	logger.Debugf("Start GetIP with ip6me")

	resp, err := http.Get("http://ip6only.me/api/")
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
	logger.Debugf("Found IP with ip6me: %+v\n", ip)
	return ip, nil
}

var ip6meProvider Provider = Provider{
	GetIP:        ip6meGetIP,
	Version:      network.IPv6,
	ProviderName: "ip6me",
}

func init() {
	ProviderList = append(ProviderList, ip6meProvider)
}
