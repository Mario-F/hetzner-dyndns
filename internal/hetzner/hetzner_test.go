package hetzner

import "testing"

func TestParseRecordID(t *testing.T) {
	zoneID, rrName, rrType, err := parseRecordID("123/@/a")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if zoneID != "123" {
		t.Fatalf("unexpected zoneID: %s", zoneID)
	}
	if rrName != "@" {
		t.Fatalf("unexpected rrName: %s", rrName)
	}
	if rrType != "A" {
		t.Fatalf("unexpected rrType: %s", rrType)
	}
}

func TestParseRecordIDInvalid(t *testing.T) {
	_, _, _, err := parseRecordID("broken-id")
	if err == nil {
		t.Fatal("expected parseRecordID to fail for invalid format")
	}
}

func TestBuildRecordID(t *testing.T) {
	got := buildRecordID("123", "@", "A")
	if got != "123/@/A" {
		t.Fatalf("unexpected id: %s", got)
	}
}

func TestFullDomain(t *testing.T) {
	if fullDomain("@", "example.com") != "example.com" {
		t.Fatal("expected root rr name to resolve to zone name")
	}
	if fullDomain("www", "example.com") != "www.example.com" {
		t.Fatal("expected subdomain to be combined with zone name")
	}
}
