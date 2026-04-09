package ports

import (
	"os"
	"testing"
)

func TestParseHexAddr_Valid(t *testing.T) {
	addr, port, err := parseHexAddr("0100007F:0050")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr != "0100007F" {
		t.Errorf("expected addr '0100007F', got '%s'", addr)
	}
	if port != 80 {
		t.Errorf("expected port 80, got %d", port)
	}
}

func TestParseHexAddr_InvalidFormat(t *testing.T) {
	_, _, err := parseHexAddr("badformat")
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestParseHexAddr_InvalidPort(t *testing.T) {
	_, _, err := parseHexAddr("0100007F:ZZZZ")
	if err == nil {
		t.Fatal("expected error for invalid port hex, got nil")
	}
}

func TestScanner_ScanReturnsBindings(t *testing.T) {
	if _, err := os.Stat("/proc/net/tcp"); os.IsNotExist(err) {
		t.Skip("skipping: /proc/net/tcp not available on this system")
	}

	s := NewScanner()
	bindings, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}

	// We can't assert exact values, but we expect a non-nil slice
	if bindings == nil {
		t.Error("expected non-nil bindings slice")
	}
}

func TestBinding_Fields(t *testing.T) {
	b := Binding{
		Protocol:  "tcp",
		LocalAddr: "0100007F",
		Port:      8080,
		PID:       1234,
		State:     "0A",
	}

	if b.Protocol != "tcp" {
		t.Errorf("unexpected protocol: %s", b.Protocol)
	}
	if b.Port != 8080 {
		t.Errorf("unexpected port: %d", b.Port)
	}
	if b.PID != 1234 {
		t.Errorf("unexpected PID: %d", b.PID)
	}
}
