// Package monitor detects changes in port bindings between scans.
package monitor

import (
	"github.com/example/portwatch/internal/ports"
)

// HistoryHook integrates a ports.History recorder into the monitor change
// pipeline so that every detected add/remove is persisted in memory.
type HistoryHook struct {
	history *ports.History
}

// NewHistoryHook wraps h so it can be used as a change observer.
func NewHistoryHook(h *ports.History) *HistoryHook {
	return &HistoryHook{history: h}
}

// RecordChanges iterates over a slice of Changes and appends each one to the
// underlying History using the canonical event-type strings "added" and
// "removed".
func (hh *HistoryHook) RecordChanges(changes []Change) {
	for _, c := range changes {
		event := eventLabel(c.Added)
		hh.history.Record(event, c.Binding)
	}
}

// History returns the underlying History store for read access.
func (hh *HistoryHook) History() *ports.History {
	return hh.history
}

func eventLabel(added bool) string {
	if added {
		return "added"
	}
	return "removed"
}
