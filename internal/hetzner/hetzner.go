package hetzner

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
)

const (
	hetznerURL string = "https://dns.hetzner.com/api/v1/"
)

type responses struct {
	Records []Record `json:"records"`
	Record  Record   `json:"record"`
	Zones   []Zone   `json:"zones"`
}

// Record is a hetzner record entry from the api
type Record struct {
	Type     string `json:"type"`
	ID       string `json:"id"`
	ZoneID   string `json:"zone_id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	TTL      int    `json:"ttl"`
	Fullname string
	Zone     Zone
}

// Zone is a hetzner zone record entry from the api
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func getRequest(token string, uri string) responses {
	client := &http.Client{}

	req, err := http.NewRequest("GET", hetznerURL+uri, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Auth-API-Token", token)
	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		panic(parseFormErr)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic(errors.New(resp.Status))
	}

	var respBody responses = responses{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&respBody)
	if decodeErr != nil {
		panic(decodeErr)
	}
	return respBody
}

// GetRecords get all dns records merged with its zone information
func GetRecords(token string) []Record {
	var (
		wg         sync.WaitGroup = sync.WaitGroup{}
		resZones   []Zone
		resRecords []Record
	)

	logger.Infof("Requesting Records and Zone info in parallel")
	wg.Add(2)
	go func() {
		resZones = getRequest(token, "zones").Zones
		logger.Debugf("Zones request finished with %d results", len(resZones))
		wg.Done()
	}()
	go func() {
		resRecords = getRequest(token, "records").Records
		logger.Debugf("Records request finished with %d results", len(resRecords))
		wg.Done()
	}()
	wg.Wait()

	aRecords := []Record{}
	for _, r := range resRecords {
		if r.Type == "A" {
			for _, z := range resZones {
				if z.ID == r.ZoneID {
					r.Zone = z
					r.Fullname = r.Name + "." + z.Name
					break
				}
			}
			aRecords = append(aRecords, r)
		}
	}

	return aRecords
}
