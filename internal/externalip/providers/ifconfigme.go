package providers

import (
	"io/ioutil"
	"log"
	"net/http"
)

func ifconfigMEGetIP() (string, error) {
	log.Println("Start GetIP with ifconfigME")

	resp, err := http.Get("http://ifconfig.me")
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
	log.Printf("Found IP wihth ifconfigME: %+v\n", ip)
	return ip, nil
}

var ifconfigMEProvider Provider = Provider{
	GetIP:        ifconfigMEGetIP,
	ProviderName: "IfconfigME",
}
