package main

import (
	"bytes"
	"testing"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/monitor"
	"portwatch/internal/ports"
)

// mockScanner satisfies the ports.Scanner interface used by monitor.New.
type mockScanner struct {
	calls    int
	results [][]ports.Binding
}

func (m *mockScanner) Scan() ([]ports.Binding, error) {
	if m.calls >= len(m.results) {
		return nil, nil
	}
	b := m.results[m.calls]
	m.calls++
	return b, nil
}

func TestTick_NoChanges(t *testing.T) {
	cfg := config.DefaultConfig()
	scanner := &mockScanner{
		results: [][]ports.Binding{
			{{Addr: "127.0.0.1", Port: 8080, Proto: "tcp"}},
			{{Addr: "127.0.0.1", Port: 8080, Proto: "tcp"}},
		},
	}

	var buf bytes.Buffer
	h := alert.NewHandler(&buf, cfg)
	mon := monitor.New(scanner, cfg)

	// First tick seeds the baseline — no changes expected.
	if err := tick(mon, h); err != nil {
		t.Fatalf("unexpected error on first tick: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output on seed tick, got: %s", buf.String())
	}

	// Second tick with identical state — still no changes.
	if err := tick(mon, h); err != nil {
		t.Fatalf("unexpected error on second tick: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output when state unchanged, got: %s", buf.String())
	}
}

func TestTick_DetectsNewBinding(t *testing.T) {
	cfg := config.DefaultConfig()
	scanner := &mockScanner{
		results: [][]ports.Binding{
			{{Addr: "0.0.0.0", Port: 22, Proto: "tcp"}},
			{
				{Addr: "0.0.0.0", Port: 22, Proto: "tcp"},
				{Addr: "0.0.0.0", Port: 9999, Proto: "tcp"},
			},
		},
	}

	var buf bytes.Buffer
	h := alert.NewHandler(&buf, cfg)
	mon := monitor.New(scanner, cfg)

	// Seed.
	if err := tick(mon, h); err != nil {
		t.Fatalf("seed tick error: %v", err)
	}

	// Second tick should detect port 9999 as added.
	if err := tick(mon, h); err != nil {
		t.Fatalf("second tick error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("expected alert output for new binding, got none")
	}
}
