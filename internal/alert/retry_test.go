package alert

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func retryChange() monitor.Change {
	return monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{Addr: "127.0.0.1", Port: 9090, Proto: "tcp"},
	}
}

func noSleep(_ time.Duration) {}

func makeRetryHandler(next Handler, maxAttempts int) *retryHandler {
	cfg := RetryConfig{MaxAttempts: maxAttempts, Delay: 0, Multiplier: 1.0}
	h := NewRetryHandler(next, cfg).(*retryHandler)
	h.sleep = noSleep
	return h
}

func TestRetryHandler_EmptyChangesNoCall(t *testing.T) {
	var called int
	next := &mockHandler{handleFn: func(_ []monitor.Change) error { called++; return nil }}
	h := makeRetryHandler(next, 3)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 0 {
		t.Errorf("expected 0 calls, got %d", called)
	}
}

func TestRetryHandler_SuccessOnFirstAttempt(t *testing.T) {
	var calls int
	next := &mockHandler{handleFn: func(_ []monitor.Change) error { calls++; return nil }}
	h := makeRetryHandler(next, 3)
	if err := h.Handle([]monitor.Change{retryChange()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetryHandler_RetriesOnError(t *testing.T) {
	var calls int
	wantErr := errors.New("temporary failure")
	next := &mockHandler{handleFn: func(_ []monitor.Change) error { calls++; return wantErr }}
	h := makeRetryHandler(next, 3)
	err := h.Handle([]monitor.Change{retryChange()})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected wantErr, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetryHandler_SucceedsOnSecondAttempt(t *testing.T) {
	var calls int
	next := &mockHandler{handleFn: func(_ []monitor.Change) error {
		calls++
		if calls < 2 {
			return errors.New("fail")
		}
		return nil
	}}
	h := makeRetryHandler(next, 3)
	if err := h.Handle([]monitor.Change{retryChange()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestRetryHandler_DefaultConfig(t *testing.T) {
	cfg := DefaultRetryConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", cfg.Multiplier)
	}
}
