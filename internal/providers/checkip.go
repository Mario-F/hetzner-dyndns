package providers

import (
	"io/ioutil"
	"net/http"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
)

func checkIPGetIP() (string, error) {
	logger.Debugf("Start GetIP with CheckIP")

	resp, err := http.Get("http://checkip.dyndns.org")
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
	logger.Debugf("Found IP with CheckIP: %+v\n", ip)
	return ip, nil
}

var checkIPProvider Provider = Provider{
	GetIP:        checkIPGetIP,
	Version:      network.IPv4,
	ProviderName: "CheckIP",
}

func init() {
	ProviderList = append(ProviderList, checkIPProvider)
}
