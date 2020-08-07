package providers

import (
	"log"
)

func ipifyGetIP() (string, error) {
	log.Println("Start GetIP with ipify")

	// TODO real call

	ip, err := captureIP("{'ip':'.122.241.150'}")
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", errIPNotFound
	}
	return ip, nil
}

var ipifyProvider Provider = Provider{
	GetIP:        ipifyGetIP,
	ProviderName: "Ipify",
}
