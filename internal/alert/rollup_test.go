package alert

import (
	"testing"
	"time"

	"github.com/patrickdappollonio/portwatch/internal/config"
	"github.com/patrickdappollonio/portwatch/internal/monitor"
	"github.com/patrickdappollonio/portwatch/internal/ports"
)

func rollupChange(port uint16) monitor.Change {
	return monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{Port: port, Proto: "tcp"},
	}
}

func rollupCfg(enabled bool, window time.Duration, max int) config.RollupConfig {
	return config.RollupConfig{Enabled: enabled, Window: window, MaxBatch: max}
}

func TestRollupHandler_DisabledForwardsImmediately(t *testing.T) {
	var got []monitor.Change
	next := HandlerFunc(func(ch []monitor.Change) error { got = append(got, ch...); return nil })
	h := NewRollupHandler(rollupCfg(false, time.Second, 10), next)

	_ = h.Handle([]monitor.Change{rollupChange(80)})

	if len(got) != 1 {
		t.Fatalf("expected 1 forwarded change, got %d", len(got))
	}
}

func TestRollupHandler_EmptyChangesSkipped(t *testing.T) {
	called := false
	next := HandlerFunc(func(ch []monitor.Change) error { called = true; return nil })
	h := NewRollupHandler(rollupCfg(true, time.Second, 10), next)

	_ = h.Handle(nil)

	if called {
		t.Fatal("expected handler not to be called for empty changes")
	}
}

func TestRollupHandler_BuffersWithinWindow(t *testing.T) {
	var got []monitor.Change
	next := HandlerFunc(func(ch []monitor.Change) error { got = append(got, ch...); return nil })

	fixedNow := time.Now()
	h := NewRollupHandler(rollupCfg(true, 5*time.Second, 100), next)
	h.now = func() time.Time { return fixedNow }

	_ = h.Handle([]monitor.Change{rollupChange(80)})
	_ = h.Handle([]monitor.Change{rollupChange(443)})

	if len(got) != 0 {
		t.Fatalf("expected 0 forwarded (still in window), got %d", len(got))
	}
}

func TestRollupHandler_FlushesAfterWindow(t *testing.T) {
	var got []monitor.Change
	next := HandlerFunc(func(ch []monitor.Change) error { got = append(got, ch...); return nil })

	start := time.Now()
	h := NewRollupHandler(rollupCfg(true, time.Second, 100), next)
	h.now = func() time.Time { return start }

	_ = h.Handle([]monitor.Change{rollupChange(80)})

	// Advance past window.
	h.now = func() time.Time { return start.Add(2 * time.Second) }
	_ = h.Handle([]monitor.Change{rollupChange(443)})

	if len(got) != 2 {
		t.Fatalf("expected 2 flushed changes, got %d", len(got))
	}
}

func TestRollupHandler_FlushesOnMaxBatch(t *testing.T) {
	var got []monitor.Change
	next := HandlerFunc(func(ch []monitor.Change) error { got = append(got, ch...); return nil })

	fixedNow := time.Now()
	h := NewRollupHandler(rollupCfg(true, time.Minute, 2), next)
	h.now = func() time.Time { return fixedNow }

	_ = h.Handle([]monitor.Change{rollupChange(80), rollupChange(443)})

	if len(got) != 2 {
		t.Fatalf("expected 2 flushed on max_batch, got %d", len(got))
	}
}

func TestRollupHandler_DrainFlushesBuffer(t *testing.T) {
	var got []monitor.Change
	next := HandlerFunc(func(ch []monitor.Change) error { got = append(got, ch...); return nil })

	fixedNow := time.Now()
	h := NewRollupHandler(rollupCfg(true, time.Minute, 100), next)
	h.now = func() time.Time { return fixedNow }

	_ = h.Handle([]monitor.Change{rollupChange(8080)})
	if len(got) != 0 {
		t.Fatal("expected nothing flushed before Drain")
	}

	_ = h.Drain()
	if len(got) != 1 {
		t.Fatalf("expected 1 change after Drain, got %d", len(got))
	}
}
