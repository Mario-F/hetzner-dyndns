package hetzner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Mario-F/hetzner-dyndns/internal/logger"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
)

const (
	hetznerURL string = "https://api.hetzner.cloud/v1/"
)

var (
	apiToken string
)

type pagination struct {
	NextPage *int `json:"next_page"`
}

type meta struct {
	Pagination pagination `json:"pagination"`
}

type rrSetRecord struct {
	Value string `json:"value"`
}

type rrSet struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Type    string        `json:"type"`
	Records []rrSetRecord `json:"records"`
}

type zonesResponse struct {
	Zones []Zone `json:"zones"`
	Meta  meta   `json:"meta"`
}

type rrSetsResponse struct {
	RRSets []rrSet `json:"rrsets"`
	Meta   meta    `json:"meta"`
}

type rrSetResponse struct {
	RRSet rrSet `json:"rrset"`
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
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// SetToken sets the api token for all following api calls
func SetToken(token string) {
	apiToken = token
}

// Update updates the record on hetzner api
func (r Record) Update() error {
	zoneID, rrName, rrType, parseErr := parseRecordID(r.ID)
	if parseErr != nil {
		if r.Zone.ID != 0 && r.Name != "" && r.Type != "" {
			zoneID = strconv.FormatInt(r.Zone.ID, 10)
			rrName = r.Name
			rrType = r.Type
		} else {
			return parseErr
		}
	}

	payload := map[string][]rrSetRecord{
		"records": {
			{Value: r.Value},
		},
	}

	apiPath := fmt.Sprintf(
		"zones/%s/rrsets/%s/%s/actions/set_records",
		url.PathEscape(zoneID),
		url.PathEscape(rrName),
		url.PathEscape(rrType),
	)

	logger.Debugf("Update rrset %s/%s/%s with value %s", zoneID, rrName, rrType, r.Value)
	if err := requestJSON(http.MethodPost, apiPath, payload, nil, apiToken, http.StatusOK, http.StatusCreated, http.StatusAccepted); err != nil {
		return err
	}

	return nil
}

func CheckToken(token string) error {
	if err := requestJSON(http.MethodGet, "zones?page=1&per_page=1", nil, &zonesResponse{}, token, http.StatusOK); err != nil {
		return fmt.Errorf("failed to check token: %w", err)
	}

	return nil
}

// GetRecords get all dns records merged with its zone information
func GetRecords(ipVersion network.IPVersion) []Record {
	var recordType string
	switch ipVersion {
	case network.IPv4:
		recordType = "A"
	case network.IPv6:
		recordType = "AAAA"
	}

	zones := mustListZones()
	logger.Infof("Requesting RRSets for %d zones", len(zones))

	records := make([]Record, 0)
	for _, zone := range zones {
		zoneID := strconv.FormatInt(zone.ID, 10)
		rrsets := mustListRRSets(zoneID, recordType)
		for _, s := range rrsets {
			if len(s.Records) != 1 {
				logger.Infof("Skip rrset %s in zone %s because it has %d records", s.ID, zone.Name, len(s.Records))
				continue
			}

			records = append(records, Record{
				Type:     s.Type,
				ID:       buildRecordID(zoneID, s.Name, s.Type),
				ZoneID:   zoneID,
				Name:     s.Name,
				Value:    s.Records[0].Value,
				Fullname: fullDomain(s.Name, zone.Name),
				Zone:     zone,
			})
		}
	}

	return records
}

// GetRecord get one record by id (includes zone)
func GetRecord(recordID string) Record {
	zoneID, rrName, rrType, parseErr := parseRecordID(recordID)
	if parseErr == nil {
		zone := mustGetZone(zoneID)
		resolvedZoneID := strconv.FormatInt(zone.ID, 10)
		rr := mustGetRRSet(resolvedZoneID, rrName, rrType)
		if len(rr.Records) != 1 {
			panic(fmt.Errorf("rrset %s has %d records, dynamic updates require exactly one", rr.ID, len(rr.Records)))
		}

		return Record{
			Type:     rr.Type,
			ID:       buildRecordID(resolvedZoneID, rr.Name, rr.Type),
			ZoneID:   resolvedZoneID,
			Name:     rr.Name,
			Value:    rr.Records[0].Value,
			Fullname: fullDomain(rr.Name, zone.Name),
			Zone:     zone,
		}
	}

	// Backward compatibility: allow an rrset ID (name/type) without zone prefix.
	zones := mustListZones()
	var match *Record
	for _, z := range zones {
		zoneID := strconv.FormatInt(z.ID, 10)
		for _, t := range []string{"A", "AAAA"} {
			rrsets := mustListRRSets(zoneID, t)
			for _, rr := range rrsets {
				if rr.ID != recordID {
					continue
				}
				if len(rr.Records) != 1 {
					panic(fmt.Errorf("rrset %s has %d records, dynamic updates require exactly one", rr.ID, len(rr.Records)))
				}
				r := Record{
					Type:     rr.Type,
					ID:       buildRecordID(zoneID, rr.Name, rr.Type),
					ZoneID:   zoneID,
					Name:     rr.Name,
					Value:    rr.Records[0].Value,
					Fullname: fullDomain(rr.Name, z.Name),
					Zone:     z,
				}
				if match != nil {
					panic(fmt.Errorf("record id %s is ambiguous across zones, use zone-prefixed id", recordID))
				}
				match = &r
			}
		}
	}
	if match == nil {
		panic(fmt.Errorf("record id %s not found", recordID))
	}

	return *match
}

func buildRecordID(zoneID string, rrName string, rrType string) string {
	return fmt.Sprintf("%s/%s/%s", zoneID, rrName, rrType)
}

func parseRecordID(recordID string) (string, string, string, error) {
	parts := strings.SplitN(recordID, "/", 3)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid record id %q, expected format <zone-id>/<rr-name>/<rr-type>", recordID)
	}

	if parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("invalid record id %q, expected format <zone-id>/<rr-name>/<rr-type>", recordID)
	}

	return parts[0], parts[1], strings.ToUpper(parts[2]), nil
}

func fullDomain(name string, zoneName string) string {
	if name == "@" || name == "" {
		return zoneName
	}
	return name + "." + zoneName
}

func mustGetZone(zoneID string) Zone {
	var resp struct {
		Zone Zone `json:"zone"`
	}
	if err := requestJSON(http.MethodGet, "zones/"+url.PathEscape(zoneID), nil, &resp, apiToken, http.StatusOK); err != nil {
		panic(err)
	}
	if resp.Zone.ID == 0 {
		panic(fmt.Errorf("zone %s not found", zoneID))
	}
	return resp.Zone
}

func mustGetRRSet(zoneID string, rrName string, rrType string) rrSet {
	var resp rrSetResponse
	apiPath := fmt.Sprintf("zones/%s/rrsets/%s/%s", url.PathEscape(zoneID), url.PathEscape(rrName), url.PathEscape(rrType))
	if err := requestJSON(http.MethodGet, apiPath, nil, &resp, apiToken, http.StatusOK); err != nil {
		panic(err)
	}
	return resp.RRSet
}

func mustListZones() []Zone {
	result := make([]Zone, 0)
	page := 1

	for {
		var resp zonesResponse
		apiPath := fmt.Sprintf("zones?page=%d&per_page=50", page)
		if err := requestJSON(http.MethodGet, apiPath, nil, &resp, apiToken, http.StatusOK); err != nil {
			panic(err)
		}

		result = append(result, resp.Zones...)
		if resp.Meta.Pagination.NextPage == nil {
			break
		}
		page = *resp.Meta.Pagination.NextPage
	}

	return result
}

func mustListRRSets(zoneID string, rrType string) []rrSet {
	result := make([]rrSet, 0)
	page := 1

	for {
		var resp rrSetsResponse
		apiPath := fmt.Sprintf(
			"zones/%s/rrsets?page=%d&per_page=50&type=%s",
			url.PathEscape(zoneID),
			page,
			url.QueryEscape(rrType),
		)
		if err := requestJSON(http.MethodGet, apiPath, nil, &resp, apiToken, http.StatusOK); err != nil {
			panic(err)
		}

		result = append(result, resp.RRSets...)
		if resp.Meta.Pagination.NextPage == nil {
			break
		}
		page = *resp.Meta.Pagination.NextPage
	}

	return result
}

func requestJSON(method string, apiPath string, payload interface{}, out interface{}, token string, expectedStatus ...int) error {
	if token == "" {
		return errors.New("missing API token")
	}

	body := []byte{}
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = data
	}

	client := &http.Client{Timeout: 20 * time.Second}
	for i := 0; i < 3; i++ {
		var reqBody io.Reader
		if len(body) > 0 {
			reqBody = strings.NewReader(string(body))
		}

		req, err := http.NewRequest(method, hetznerURL+apiPath, reqBody)
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		if len(body) > 0 {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := client.Do(req)
		if err != nil {
			if i < 2 {
				time.Sleep(2 * time.Second)
				continue
			}
			return err
		}

		respBytes, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			return readErr
		}

		if !containsStatus(resp.StatusCode, expectedStatus) {
			if (resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500) && i < 2 {
				time.Sleep(2 * time.Second)
				continue
			}
			return fmt.Errorf("hetzner cloud api %s %s returned %d: %s", method, apiPath, resp.StatusCode, strings.TrimSpace(string(respBytes)))
		}

		if out != nil && len(respBytes) > 0 {
			dec := json.NewDecoder(strings.NewReader(string(respBytes)))
			dec.UseNumber()
			if err := dec.Decode(out); err != nil {
				return err
			}
		}

		return nil
	}

	return fmt.Errorf("request failed after retries: %s %s", method, apiPath)
}

func containsStatus(status int, allowed []int) bool {
	for _, s := range allowed {
		if status == s {
			return true
		}
	}
	return false
}
