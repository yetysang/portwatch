// Package ports provides port scanning, filtering, and enrichment utilities.
package ports

import (
	"sync"
	"time"
)

// HistoryEntry records a binding observation at a point in time.
type HistoryEntry struct {
	Binding   Binding
	SeenAt    time.Time
	EventType string // "added" or "removed"
}

// History maintains a bounded in-memory log of port binding events.
type History struct {
	mu      sync.RWMutex
	entries []HistoryEntry
	maxSize int
}

// NewHistory creates a History that retains at most maxSize entries.
func NewHistory(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 256
	}
	return &History{maxSize: maxSize}
}

// Record appends a new entry, evicting the oldest if the buffer is full.
func (h *History) Record(eventType string, b Binding) {
	h.mu.Lock()
	defer h.mu.Unlock()
	entry := HistoryEntry{
		Binding:   b,
		SeenAt:    time.Now(),
		EventType: eventType,
	}
	if len(h.entries) >= h.maxSize {
		h.entries = h.entries[1:]
	}
	h.entries = append(h.entries, entry)
}

// Entries returns a shallow copy of all recorded entries.
func (h *History) Entries() []HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]HistoryEntry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Len returns the current number of entries.
func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.entries)
}

// Clear removes all entries from the history.
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = h.entries[:0]
}
