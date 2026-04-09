package monitor

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// Change represents a detected port binding change.
type Change struct {
	Type    string // "added" or "removed"
	Binding ports.Binding
}

// Monitor watches for port binding changes at a given interval.
type Monitor struct {
	scanner  *ports.Scanner
	interval time.Duration
	previous map[string]ports.Binding
	Changes  chan Change
	Stop     chan struct{}
}

// New creates a new Monitor with the given polling interval.
func New(scanner *ports.Scanner, interval time.Duration) *Monitor {
	return &Monitor{
		scanner:  scanner,
		interval: interval,
		previous: make(map[string]ports.Binding),
		Changes:  make(chan Change, 64),
		Stop:     make(chan struct{}),
	}
}

// Run starts the monitoring loop. It blocks until Stop is closed.
func (m *Monitor) Run() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// Populate initial state without emitting changes.
	if bindings, err := m.scanner.Scan(); err == nil {
		for _, b := range bindings {
			m.previous[bindingKey(b)] = b
		}
	}

	for {
		select {
		case <-m.Stop:
			return
		case <-ticker.C:
			if err := m.poll(); err != nil {
				log.Printf("portwatch: scan error: %v", err)
			}
		}
	}
}

func (m *Monitor) poll() error {
	bindings, err := m.scanner.Scan()
	if err != nil {
		return err
	}

	current := make(map[string]ports.Binding, len(bindings))
	for _, b := range bindings {
		current[bindingKey(b)] = b
	}

	for key, b := range current {
		if _, existed := m.previous[key]; !existed {
			m.Changes <- Change{Type: "added", Binding: b}
		}
	}

	for key, b := range m.previous {
		if _, exists := current[key]; !exists {
			m.Changes <- Change{Type: "removed", Binding: b}
		}
	}

	m.previous = current
	return nil
}

func bindingKey(b ports.Binding) string {
	return b.Proto + "|" + b.LocalAddr + ":" + string(rune(b.LocalPort))
}
