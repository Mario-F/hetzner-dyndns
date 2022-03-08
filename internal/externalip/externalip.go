package externalip

import (
	"errors"
	"math/rand"
	"time"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
	"github.com/Mario-F/hetzner-dyndns/internal/providers"
)

type ExternalIP struct {
	IP      string            `json:"ip"`
	Version network.IPVersion `json:"version"`
}

// GetExternalIP gets the actual external IP by different Provides
func GetExternalIP(version network.IPVersion) (ExternalIP, error) {
	var pList []providers.Provider

	for _, p := range providers.ProviderList {
		if p.Version == version {
			pList = append(pList, p)
		}
	}
	logger.Debugf("Found %d providers for ip version %s", len(pList), string(version))

	result := ExternalIP{}
	result.Version = version

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
			result.IP = ip
			return result, nil
		}
		found = ip
	}

	return result, errors.New("at least 2 providers doesnt confirm external IP")
}
