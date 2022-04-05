package hetzner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
)

const (
	hetznerURL string = "https://dns.hetzner.com/api/v1/"
)

var (
	apiToken string
)

type responses struct {
	Records []Record `json:"records"`
	Record  Record   `json:"record"`
	Zones   []Zone   `json:"zones"`
	Zone    Zone     `json:"zone"`
}

// Record is a hetzner record entry from the api
type Record struct {
	Type     string `json:"type"`
	ID       string `json:"id"`
	ZoneID   string `json:"zone_id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	TTL      int    `json:"ttl,omitempty"`
	Fullname string `json:"-"`
	Zone     Zone   `json:"-"`
}

// Zone is a hetzner zone record entry from the api
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SetToken sets the api token for all following api calls
func SetToken(token string) {
	apiToken = token
}

// Update updates the record on hetzner api
func (r Record) Update() error {
	body, err := json.Marshal(r)
	if err != nil {
		return err
	}
	logger.Debugf("Update record will send following json:\n%+v", string(body))

	client := &http.Client{}
	req, err := http.NewRequest("PUT", hetznerURL+"records/"+r.ID, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-API-Token", apiToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Failed to update record, statuscode: " + fmt.Sprint(resp.StatusCode))
	}

	return nil
}

func getRequest(uri string) responses {
	client := &http.Client{}

	req, err := http.NewRequest("GET", hetznerURL+uri, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Auth-API-Token", apiToken)
	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		panic(parseFormErr)
	}

	var resp *http.Response
	finshed := false
	for i := 0; i < 3 && !finshed; i++ {
		// TODO: Retry on timeout because some times hetzner is not ready

		resp, err = client.Do(req)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode == 404 {
			logger.Infof("Hetzner API returned 404, retry.")
			time.Sleep(2 * time.Second)
			continue
		}
		finshed = true
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
func GetRecords(ipVersion network.IPVersion) []Record {
	var (
		wg         sync.WaitGroup = sync.WaitGroup{}
		resZones   []Zone
		resRecords []Record
	)

	logger.Infof("Requesting Records and Zone info in parallel")
	wg.Add(2)
	go func() {
		resZones = getRequest("zones").Zones
		logger.Debugf("Zones request finished with %d results", len(resZones))
		wg.Done()
	}()
	go func() {
		resRecords = getRequest("records").Records
		logger.Debugf("Records request finished with %d results", len(resRecords))
		wg.Done()
	}()
	wg.Wait()

	var recordType string
	switch ipVersion {
	case network.IPv4:
		recordType = "A"
	case network.IPv6:
		recordType = "AAAA"
	}

	aRecords := []Record{}
	for _, r := range resRecords {
		if r.Type == recordType {
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

// GetRecord get one record by id (includes zone)
func GetRecord(recordID string) Record {
	resRecord := getRequest("records/" + recordID).Record
	resRecord.Zone = getRequest("zones/" + resRecord.ZoneID).Zone
	resRecord.Fullname = resRecord.Name + "." + resRecord.Zone.Name

	return resRecord
}
