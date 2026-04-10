package ports

import (
	"errors"
	"testing"
)

func stubLookupAddr(hosts []string, err error) func(string) ([]string, error) {
	return func(string) ([]string, error) {
		return hosts, err
	}
}

func stubLookupPort(name string) func(string, int) string {
	return func(string, int) string {
		return name
	}
}

func TestResolver_ResolveUsesHostname(t *testing.T) {
	r := &Resolver{
		lookupAddr: stubLookupAddr([]string{"example.com."}, nil),
		lookupPort: stubLookupPort("http"),
	}
	b := Binding{IP: "93.184.216.34", Port: 80, Proto: "tcp"}
	rb := r.Resolve(b)
	if rb.Hostname != "example.com." {
		t.Errorf("expected hostname %q, got %q", "example.com.", rb.Hostname)
	}
	if rb.ServiceName != "http" {
		t.Errorf("expected service %q, got %q", "http", rb.ServiceName)
	}
}

func TestResolver_ResolveFallsBackToIP(t *testing.T) {
	r := &Resolver{
		lookupAddr: stubLookupAddr(nil, errors.New("no reverse")),
		lookupPort: stubLookupPort("ssh"),
	}
	b := Binding{IP: "127.0.0.1", Port: 22, Proto: "tcp"}
	rb := r.Resolve(b)
	if rb.Hostname != "127.0.0.1" {
		t.Errorf("expected IP fallback %q, got %q", "127.0.0.1", rb.Hostname)
	}
}

func TestResolver_ResolveAll(t *testing.T) {
	r := &Resolver{
		lookupAddr: stubLookupAddr([]string{"localhost."}, nil),
		lookupPort: stubLookupPort("http"),
	}
	bindings := []Binding{
		{IP: "127.0.0.1", Port: 80, Proto: "tcp"},
		{IP: "127.0.0.1", Port: 443, Proto: "tcp"},
	}
	resolved := r.ResolveAll(bindings)
	if len(resolved) != 2 {
		t.Fatalf("expected 2 resolved bindings, got %d", len(resolved))
	}
	for _, rb := range resolved {
		if rb.Hostname == "" {
			t.Error("expected non-empty hostname")
		}
	}
}

func TestLookupServiceName_WellKnown(t *testing.T) {
	cases := []struct {
		port    int
		expected string
	}{
		{22, "ssh"},
		{80, "http"},
		{443, "https"},
		{3306, "mysql"},
		{9999, "9999"},
	}
	for _, tc := range cases {
		got := lookupServiceName("tcp", tc.port)
		if got != tc.expected {
			t.Errorf("port %d: expected %q, got %q", tc.port, tc.expected, got)
		}
	}
}

func TestNewResolver_NotNil(t *testing.T) {
	r := NewResolver()
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}
