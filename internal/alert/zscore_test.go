package alert

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func zscoreCfg() config.ZScoreConfig {
	return config.ZScoreConfig{
		Enabled:    true,
		WindowSize: 20,
		Threshold:  2.0,
		MinSamples: 5,
		Cooldown:   10 * time.Second,
	}
}

func zscoreChange(n int) []monitor.Change {
	out := make([]monitor.Change, n)
	for i := range out {
		out[i] = monitor.Change{
			Kind:    monitor.Added,
			Binding: ports.Binding{Port: uint16(8000 + i), Proto: "tcp"},
		}
	}
	return out
}

func TestZScoreHandler_EmptyChangesSkipped(t *testing.T) {
	called := false
	next := HandlerFunc(func(ch []monitor.Change) error { called = true; return nil })
	h := NewZScoreHandler(zscoreCfg(), next)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("next should not be called for empty changes")
	}
}

func TestZScoreHandler_BelowMinSamplesNotForwarded(t *testing.T) {
	called := false
	next := HandlerFunc(func(ch []monitor.Change) error { called = true; return nil })
	h := NewZScoreHandler(zscoreCfg(), next)
	// feed 4 ticks (< min_samples=5)
	for i := 0; i < 4; i++ {
		_ = h.Handle(zscoreChange(2))
	}
	if called {
		t.Fatal("next should not be called before min_samples reached")
	}
}

func TestZScoreHandler_AnomalyForwarded(t *testing.T) {
	called := 0
	next := HandlerFunc(func(ch []monitor.Change) error { called++; return nil })
	cfg := zscoreCfg()
	h := NewZScoreHandler(cfg, next)
	// Establish a stable baseline: 10 ticks of 1 change each
	for i := 0; i < 10; i++ {
		_ = h.Handle(zscoreChange(1))
	}
	// Now send a spike of 50 changes — should exceed threshold
	_ = h.Handle(zscoreChange(50))
	if called != 1 {
		t.Fatalf("expected next called once, got %d", called)
	}
}

func TestZScoreHandler_CooldownSuppressesDuplicate(t *testing.T) {
	called := 0
	next := HandlerFunc(func(ch []monitor.Change) error { called++; return nil })
	cfg := zscoreCfg()
	fixed := time.Now()
	h := NewZScoreHandler(cfg, next)
	h.now = func() time.Time { return fixed }
	// Establish baseline
	for i := 0; i < 10; i++ {
		_ = h.Handle(zscoreChange(1))
	}
	// First spike fires
	_ = h.Handle(zscoreChange(50))
	// Second spike within cooldown should be suppressed
	_ = h.Handle(zscoreChange(50))
	if called != 1 {
		t.Fatalf("expected exactly 1 forward, got %d", called)
	}
}

func TestZScoreHandler_DrainDelegates(t *testing.T) {
	sentinel := errors.New("drain called")
	next := HandlerFunc(func(_ []monitor.Change) error { return nil })
	h := NewZScoreHandler(zscoreCfg(), &drainableStub{next: next, drainErr: sentinel})
	if err := h.Drain(); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

// drainableStub is a helper that lets us inject a drain error.
type drainableStub struct {
	next     Handler
	drainErr error
}

func (d *drainableStub) Handle(ch []monitor.Change) error { return d.next.Handle(ch) }
func (d *drainableStub) Drain() error                     { return d.drainErr }
