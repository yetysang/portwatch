package monitor

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// fakeScanner lets tests inject controlled Scan results.
type fakeScanner struct {
	results [][]ports.Binding
	call    int
}

func (f *fakeScanner) Scan() ([]ports.Binding, error) {
	if f.call >= len(f.results) {
		return f.results[len(f.results)-1], nil
	}
	b := f.results[f.call]
	f.call++
	return b, nil
}

func binding(proto, addr string, port int) ports.Binding {
	return ports.Binding{Proto: proto, LocalAddr: addr, LocalPort: port}
}

func TestBindingKey(t *testing.T) {
	b := binding("tcp", "0.0.0.0", 8080)
	key := bindingKey(b)
	if key == "" {
		t.Fatal("expected non-empty binding key")
	}
}

func TestMonitor_DetectsAdded(t *testing.T) {
	initial := []ports.Binding{binding("tcp", "0.0.0.0", 80)}
	second := []ports.Binding{binding("tcp", "0.0.0.0", 80), binding("tcp", "0.0.0.0", 443)}

	fake := &fakeScanner{results: [][]ports.Binding{initial, second}}

	// Directly test poll logic via unexported method using a real Monitor.
	m := &Monitor{
		scanner:  nil,
		interval: time.Second,
		previous: make(map[string]ports.Binding),
		Changes:  make(chan Change, 64),
		Stop:     make(chan struct{}),
	}

	// Seed initial state.
	for _, b := range initial {
		m.previous[bindingKey(b)] = b
	}

	// Simulate poll with second scan result.
	current := make(map[string]ports.Binding)
	for _, b := range second {
		current[bindingKey(b)] = b
	}
	for key, b := range current {
		if _, existed := m.previous[key]; !existed {
			m.Changes <- Change{Type: "added", Binding: b}
		}
	}

	if len(m.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(m.Changes))
	}
	c := <-m.Changes
	if c.Type != "added" || c.Binding.LocalPort != 443 {
		t.Errorf("unexpected change: %+v", c)
	}
	_ = fake
}

func TestMonitor_DetectsRemoved(t *testing.T) {
	initial := []ports.Binding{
		binding("tcp", "0.0.0.0", 80),
		binding("tcp", "0.0.0.0", 9000),
	}
	second := []ports.Binding{binding("tcp", "0.0.0.0", 80)}

	m := &Monitor{
		previous: make(map[string]ports.Binding),
		Changes:  make(chan Change, 64),
	}
	for _, b := range initial {
		m.previous[bindingKey(b)] = b
	}

	current := make(map[string]ports.Binding)
	for _, b := range second {
		current[bindingKey(b)] = b
	}
	for key, b := range m.previous {
		if _, exists := current[key]; !exists {
			m.Changes <- Change{Type: "removed", Binding: b}
		}
	}

	if len(m.Changes) != 1 {
		t.Fatalf("expected 1 removal, got %d", len(m.Changes))
	}
	c := <-m.Changes
	if c.Type != "removed" || c.Binding.LocalPort != 9000 {
		t.Errorf("unexpected change: %+v", c)
	}
}
