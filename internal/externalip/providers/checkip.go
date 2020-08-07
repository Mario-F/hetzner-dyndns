package providers

import (
	"io/ioutil"
	"log"
	"net/http"
)

func checkIPGetIP() (string, error) {
	log.Println("Start GetIP with CheckIP")

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

	ip, err := captureIP(string(body))
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", errIPNotFound
	}
	log.Printf("Found IP wihth CheckIP: %+v\n", ip)
	return ip, nil
}

var checkIPProvider Provider = Provider{
	GetIP:        checkIPGetIP,
	ProviderName: "CheckIP",
}
