package monitor

import (
	"testing"

	"github.com/example/portwatch/internal/ports"
)

func makeChange(added bool, port uint16) Change {
	return Change{
		Binding: ports.Binding{IP: "0.0.0.0", Port: port, Protocol: "tcp"},
		Added:   added,
	}
}

func TestHistoryHook_RecordAdded(t *testing.T) {
	h := ports.NewHistory(10)
	hook := NewHistoryHook(h)
	hook.RecordChanges([]Change{makeChange(true, 8080)})

	entries := h.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].EventType != "added" {
		t.Errorf("expected 'added', got %q", entries[0].EventType)
	}
	if entries[0].Binding.Port != 8080 {
		t.Errorf("expected port 8080, got %d", entries[0].Binding.Port)
	}
}

func TestHistoryHook_RecordRemoved(t *testing.T) {
	h := ports.NewHistory(10)
	hook := NewHistoryHook(h)
	hook.RecordChanges([]Change{makeChange(false, 443)})

	entries := h.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].EventType != "removed" {
		t.Errorf("expected 'removed', got %q", entries[0].EventType)
	}
}

func TestHistoryHook_MultipleChanges(t *testing.T) {
	h := ports.NewHistory(10)
	hook := NewHistoryHook(h)
	changes := []Change{
		makeChange(true, 80),
		makeChange(true, 443),
		makeChange(false, 8080),
	}
	hook.RecordChanges(changes)

	if h.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", h.Len())
	}
}

func TestHistoryHook_EmptyChanges(t *testing.T) {
	h := ports.NewHistory(10)
	hook := NewHistoryHook(h)
	hook.RecordChanges(nil)

	if h.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", h.Len())
	}
}

func TestHistoryHook_HistoryAccessor(t *testing.T) {
	h := ports.NewHistory(10)
	hook := NewHistoryHook(h)
	if hook.History() != h {
		t.Error("History() should return the same *ports.History passed to NewHistoryHook")
	}
}

func TestEventLabel(t *testing.T) {
	if eventLabel(true) != "added" {
		t.Error("eventLabel(true) should return 'added'")
	}
	if eventLabel(false) != "removed" {
		t.Error("eventLabel(false) should return 'removed'")
	}
}
