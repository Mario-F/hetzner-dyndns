package providers

import (
	"log"
)

func ipifyGetIP() (string, error) {
	log.Println("Start GetIP with ipify")

	// TODO real call

	ip, err := captureIP("{'ip':'20.122.241.150'}")
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
