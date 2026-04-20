package alert

import (
	"errors"
	"testing"
	"time"
)

// failingHandler is a test Handler that always returns an error.
type failingHandler struct{ calls int }

func (f *failingHandler) Handle(_ []Change) error {
	f.calls++
	return errors.New("downstream error")
}
func (f *failingHandler) Drain() error { return nil }

// succeedingHandler is a test Handler that always succeeds.
type succeedingHandler struct{ calls int }

func (s *succeedingHandler) Handle(_ []Change) error { s.calls++; return nil }
func (s *succeedingHandler) Drain() error             { return nil }

func circuitChange() []Change {
	return []Change{{Kind: ChangeAdded, Port: 8080, Proto: "tcp", Addr: "0.0.0.0"}}
}

func TestCircuitBreaker_InitiallyClosed(t *testing.T) {
	h := &succeedingHandler{}
	cb := NewCircuitBreakerHandler(h, 3, 10*time.Second)
	if cb.State() != CircuitClosed {
		t.Fatalf("expected Closed, got %v", cb.State())
	}
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	h := &failingHandler{}
	cb := NewCircuitBreakerHandler(h, 3, 10*time.Second)
	for i := 0; i < 3; i++ {
		_ = cb.Handle(circuitChange())
	}
	if cb.State() != CircuitOpen {
		t.Fatalf("expected Open after %d failures, got %v", 3, cb.State())
	}
}

func TestCircuitBreaker_BlocksWhenOpen(t *testing.T) {
	h := &failingHandler{}
	cb := NewCircuitBreakerHandler(h, 1, 10*time.Second)
	_ = cb.Handle(circuitChange()) // trips open

	callsBefore := h.calls
	err := cb.Handle(circuitChange())
	if err == nil {
		t.Fatal("expected error when circuit is open")
	}
	if h.calls != callsBefore {
		t.Fatal("inner handler should not be called when circuit is open")
	}
}

func TestCircuitBreaker_HalfOpenAfterReset(t *testing.T) {
	h := &succeedingHandler{}
	now := time.Now()
	cb := NewCircuitBreakerHandler(h, 1, 5*time.Second)
	cb.now = func() time.Time { return now }

	// force open state
	cb.mu.Lock()
	cb.state = CircuitOpen
	cb.openedAt = now.Add(-10 * time.Second)
	cb.mu.Unlock()

	err := cb.Handle(circuitChange())
	if err != nil {
		t.Fatalf("expected success after reset window, got %v", err)
	}
	if cb.State() != CircuitClosed {
		t.Fatalf("expected Closed after successful probe, got %v", cb.State())
	}
}

func TestCircuitBreaker_ResetsOnSuccess(t *testing.T) {
	h := &succeedingHandler{}
	cb := NewCircuitBreakerHandler(h, 3, 10*time.Second)
	cb.mu.Lock()
	cb.failures = 2
	cb.mu.Unlock()

	_ = cb.Handle(circuitChange())
	cb.mu.Lock()
	f := cb.failures
	cb.mu.Unlock()
	if f != 0 {
		t.Fatalf("expected failures reset to 0, got %d", f)
	}
}

func TestCircuitBreaker_EmptyChangesSkipped(t *testing.T) {
	h := &failingHandler{}
	cb := NewCircuitBreakerHandler(h, 1, 10*time.Second)
	err := cb.Handle(nil)
	if err != nil {
		t.Fatalf("unexpected error for empty changes: %v", err)
	}
	if h.calls != 0 {
		t.Fatal("inner handler should not be called for empty changes")
	}
}
