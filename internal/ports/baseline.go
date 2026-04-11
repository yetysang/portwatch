// Package ports provides utilities for scanning and comparing port bindings.
package ports

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BaselineEntry records a binding that has been acknowledged as expected.
type BaselineEntry struct {
	Proto     string    `json:"proto"`
	Addr      string    `json:"addr"`
	Port      int       `json:"port"`
	Process   string    `json:"process,omitempty"`
	AddedAt   time.Time `json:"added_at"`
}

// Baseline holds the set of acknowledged port bindings.
type Baseline struct {
	path    string
	entries map[string]BaselineEntry
}

// NewBaseline creates a Baseline backed by the given file path.
// If the file does not exist, an empty baseline is returned.
func NewBaseline(path string) (*Baseline, error) {
	b := &Baseline{
		path:    path,
		entries: make(map[string]BaselineEntry),
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return b, nil
		}
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &b.entries); err != nil {
		return nil, fmt.Errorf("baseline: parse %s: %w", path, err)
	}
	return b, nil
}

// Contains reports whether the given binding is in the baseline.
func (b *Baseline) Contains(proto, addr string, port int) bool {
	key := proto + ":" + addr + ":" + itoa(port)
	_, ok := b.entries[key]
	return ok
}

// Add adds a binding to the baseline and persists it to disk.
func (b *Baseline) Add(entry BaselineEntry) error {
	key := entry.Proto + ":" + entry.Addr + ":" + itoa(entry.Port)
	entry.AddedAt = time.Now().UTC()
	b.entries[key] = entry
	return b.save()
}

// Len returns the number of entries in the baseline.
func (b *Baseline) Len() int { return len(b.entries) }

func (b *Baseline) save() error {
	if err := os.MkdirAll(filepath.Dir(b.path), 0o755); err != nil {
		return fmt.Errorf("baseline: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(b.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	return os.WriteFile(b.path, data, 0o644)
}
