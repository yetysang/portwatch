package ports

import (
	"testing"
	"time"
)

func sampleBinding(port uint16) Binding {
	return Binding{IP: "127.0.0.1", Port: port, Protocol: "tcp"}
}

func TestHistory_NewDefaultMaxSize(t *testing.T) {
	h := NewHistory(0)
	if h.maxSize != 256 {
		t.Fatalf("expected default maxSize 256, got %d", h.maxSize)
	}
}

func TestHistory_RecordAndLen(t *testing.T) {
	h := NewHistory(10)
	h.Record("added", sampleBinding(8080))
	h.Record("removed", sampleBinding(9090))
	if h.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", h.Len())
	}
}

func TestHistory_EntriesContent(t *testing.T) {
	h := NewHistory(10)
	before := time.Now()
	h.Record("added", sampleBinding(3000))
	after := time.Now()

	entries := h.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.EventType != "added" {
		t.Errorf("expected eventType 'added', got %q", e.EventType)
	}
	if e.Binding.Port != 3000 {
		t.Errorf("expected port 3000, got %d", e.Binding.Port)
	}
	if e.SeenAt.Before(before) || e.SeenAt.After(after) {
		t.Errorf("SeenAt %v not in expected range", e.SeenAt)
	}
}

func TestHistory_Eviction(t *testing.T) {
	h := NewHistory(3)
	for i := uint16(1); i <= 5; i++ {
		h.Record("added", sampleBinding(i))
	}
	if h.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", h.Len())
	}
	entries := h.Entries()
	if entries[0].Binding.Port != 3 {
		t.Errorf("expected oldest surviving port 3, got %d", entries[0].Binding.Port)
	}
}

func TestHistory_Clear(t *testing.T) {
	h := NewHistory(10)
	h.Record("added", sampleBinding(80))
	h.Record("added", sampleBinding(443))
	h.Clear()
	if h.Len() != 0 {
		t.Fatalf("expected 0 entries after clear, got %d", h.Len())
	}
}

func TestHistory_EntriesReturnsCopy(t *testing.T) {
	h := NewHistory(10)
	h.Record("added", sampleBinding(8080))
	entries := h.Entries()
	entries[0].EventType = "mutated"
	original := h.Entries()
	if original[0].EventType != "added" {
		t.Error("Entries() should return a copy, not a reference")
	}
}
