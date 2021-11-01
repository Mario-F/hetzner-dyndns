package externalip

import (
	"errors"
	"math/rand"
	"time"

	"github.com/Mario-F/hetzner-dyndns/internal/externalip/providers"
	"github.com/Mario-F/hetzner-dyndns/internal/logger"
)

// GetExternalIP gets the actual external IP by different Provides
func GetExternalIP() (string, error) {
	var pList []providers.Provider = providers.ProviderList

	// shuffle ProviderList
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(pList), func(i, j int) { pList[i], pList[j] = pList[j], pList[i] })

	var found string
	for _, p := range pList {
		ip, err := p.GetIP()
		if err != nil {
			logger.Debugf("Provider %+v failed", p.ProviderName)
			continue
		}
		logger.Debugf("Name: %+v, IP: %+v", p.ProviderName, ip)
		// found is equal ip means that is the second confirmation
		if found == ip {
			logger.Debugf("Second confirmation for IP: %+v", ip)
			return ip, nil
		}
		found = ip
	}

	return "", errors.New("At least 2 providers doesnt confirm external IP")
}
