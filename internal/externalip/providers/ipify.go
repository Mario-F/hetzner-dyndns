package providers

import (
	"io/ioutil"
	"log"
	"net/http"
)

func ipifyGetIP() (string, error) {
	log.Println("Start GetIP with ipify")

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

	ip, err := captureIP(string(body))
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", errIPNotFound
	}
	log.Printf("Found IP wihth ipify: %+v\n", ip)
	return ip, nil
}

var ipifyProvider Provider = Provider{
	GetIP:        ipifyGetIP,
	ProviderName: "Ipify",
}
