package alert

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

// stubHandler records calls for assertion in tests.
type stubHandler struct {
	handleCalls int
	drainCalls  int
	handleErr   error
	drainErr    error
	lastChanges []monitor.Change
}

func (s *stubHandler) Handle(changes []monitor.Change) error {
	s.handleCalls++
	s.lastChanges = changes
	return s.handleErr
}

func (s *stubHandler) Drain() error {
	s.drainCalls++
	return s.drainErr
}

func TestMultiHandler_Len(t *testing.T) {
	m := NewMultiHandler(&stubHandler{}, &stubHandler{})
	if m.Len() != 2 {
		t.Fatalf("expected 2 handlers, got %d", m.Len())
	}
}

func TestMultiHandler_EmptyChangesSkipsHandle(t *testing.T) {
	h := &stubHandler{}
	m := NewMultiHandler(h)
	if err := m.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.handleCalls != 0 {
		t.Fatalf("expected 0 handle calls, got %d", h.handleCalls)
	}
}

func TestMultiHandler_HandleForwardsToAll(t *testing.T) {
	h1, h2 := &stubHandler{}, &stubHandler{}
	m := NewMultiHandler(h1, h2)
	changes := []monitor.Change{{Kind: monitor.Added}}
	if err := m.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h1.handleCalls != 1 || h2.handleCalls != 1 {
		t.Fatalf("expected both handlers called once")
	}
}

func TestMultiHandler_DrainForwardsToAll(t *testing.T) {
	h1, h2 := &stubHandler{}, &stubHandler{}
	m := NewMultiHandler(h1, h2)
	if err := m.Drain(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h1.drainCalls != 1 || h2.drainCalls != 1 {
		t.Fatalf("expected both handlers drained once")
	}
}

func TestMultiHandler_CollectsErrors(t *testing.T) {
	h1 := &stubHandler{handleErr: errors.New("backend1 down")}
	h2 := &stubHandler{handleErr: errors.New("backend2 down")}
	m := NewMultiHandler(h1, h2)
	changes := []monitor.Change{{Kind: monitor.Added}}
	err := m.Handle(changes)
	if err == nil {
		t.Fatal("expected combined error, got nil")
	}
	if h1.handleCalls != 1 || h2.handleCalls != 1 {
		t.Fatal("all handlers must be called even when one errors")
	}
}

func TestMultiHandler_DrainCollectsErrors(t *testing.T) {
	h := &stubHandler{drainErr: errors.New("flush failed")}
	m := NewMultiHandler(h)
	if err := m.Drain(); err == nil {
		t.Fatal("expected error from Drain")
	}
}
