package externalip

import (
	"log"

	"github.com/Mario-F/hetzner-dyndns/internal/externalip/providers"
)

// GetExternalIP gets the actual external IP by different Provides
func GetExternalIP() (string, error) {
	for _, p := range providers.ProviderList {
		ip, err := p.GetIP()
		if err == nil {
			log.Printf("Name: %+v, IP: %+v", p.ProviderName, ip)
		}
	}

	return "test", nil
}
