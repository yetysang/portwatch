package alert

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func throttleChange(port uint16) monitor.Change {
	return monitor.Change{
		Type: monitor.Added,
		Binding: ports.Binding{Port: port, Proto: "tcp", Addr: "0.0.0.0"},
	}
}

type countingHandler struct {
	calls atomic.Int32
	err   error
}

func (c *countingHandler) Handle(_ []monitor.Change) error {
	c.calls.Add(1)
	return c.err
}

func TestThrottleHandler_FirstCallForwarded(t *testing.T) {
	inner := &countingHandler{}
	th := NewThrottleHandler(inner, 5*time.Second)

	if err := th.Handle([]monitor.Change{throttleChange(8080)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", inner.calls.Load())
	}
}

func TestThrottleHandler_DuplicateWithinWindowSuppressed(t *testing.T) {
	inner := &countingHandler{}
	fixed := time.Unix(1_000_000, 0)
	th := NewThrottleHandler(inner, 10*time.Second)
	th.now = func() time.Time { return fixed }

	changes := []monitor.Change{throttleChange(9090)}
	_ = th.Handle(changes)
	_ = th.Handle(changes)
	_ = th.Handle(changes)

	if inner.calls.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", inner.calls.Load())
	}
}

func TestThrottleHandler_AllowedAfterWindow(t *testing.T) {
	inner := &countingHandler{}
	clock := time.Unix(1_000_000, 0)
	th := NewThrottleHandler(inner, 5*time.Second)
	th.now = func() time.Time { return clock }

	changes := []monitor.Change{throttleChange(3000)}
	_ = th.Handle(changes)

	clock = clock.Add(6 * time.Second)
	_ = th.Handle(changes)

	if inner.calls.Load() != 2 {
		t.Fatalf("expected 2 calls, got %d", inner.calls.Load())
	}
}

func TestThrottleHandler_DifferentBatchNotSuppressed(t *testing.T) {
	inner := &countingHandler{}
	fixed := time.Unix(1_000_000, 0)
	th := NewThrottleHandler(inner, 10*time.Second)
	th.now = func() time.Time { return fixed }

	_ = th.Handle([]monitor.Change{throttleChange(1111)})
	_ = th.Handle([]monitor.Change{throttleChange(2222)})

	if inner.calls.Load() != 2 {
		t.Fatalf("expected 2 calls, got %d", inner.calls.Load())
	}
}

func TestThrottleHandler_EmptyChangesSkipped(t *testing.T) {
	inner := &countingHandler{}
	th := NewThrottleHandler(inner, time.Second)

	_ = th.Handle(nil)
	_ = th.Handle([]monitor.Change{})

	if inner.calls.Load() != 0 {
		t.Fatalf("expected 0 calls, got %d", inner.calls.Load())
	}
}

func TestThrottleHandler_PropagatesError(t *testing.T) {
	want := errors.New("downstream failure")
	inner := &countingHandler{err: want}
	th := NewThrottleHandler(inner, time.Second)

	got := th.Handle([]monitor.Change{throttleChange(5000)})
	if !errors.Is(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
