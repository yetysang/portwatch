package daemon

import (
	"context"
	"log/slog"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rorycl/portwatch/internal/config"
)

func defaultWatcherCfg() config.WatcherConfig {
	cfg := config.DefaultWatcherConfig()
	cfg.Enabled = true
	cfg.PollInterval = 50 * time.Millisecond
	cfg.DebounceDelay = 10 * time.Millisecond
	return cfg
}

func TestConfigWatcher_DisabledDoesNothing(t *testing.T) {
	cfg := config.DefaultWatcherConfig() // Enabled: false
	var called int32
	w := NewConfigWatcher(cfg, "/nonexistent", slog.Default(), func(string) error {
		atomic.AddInt32(&called, 1)
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	w.Run(ctx)
	if atomic.LoadInt32(&called) != 0 {
		t.Fatal("expected no reload calls when disabled")
	}
}

func TestConfigWatcher_DetectsFileChange(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "portwatch-cfg-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	_ = tmp.Close()

	var calls int32
	w := NewConfigWatcher(defaultWatcherCfg(), tmp.Name(), slog.Default(), func(string) error {
		atomic.AddInt32(&calls, 1)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		w.Run(ctx)
		close(done)
	}()

	time.Sleep(80 * time.Millisecond)
	if err := os.WriteFile(tmp.Name(), []byte("interval = \"10s\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)
	cancel()
	<-done

	if atomic.LoadInt32(&calls) == 0 {
		t.Fatal("expected at least one reload call after file change")
	}
}

func TestConfigWatcher_NoReloadWhenFileUnchanged(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "portwatch-cfg-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	_ = tmp.WriteString("interval = \"10s\"\n")
	_ = tmp.Close()

	var calls int32
	w := NewConfigWatcher(defaultWatcherCfg(), tmp.Name(), slog.Default(), func(string) error {
		atomic.AddInt32(&calls, 1)
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	if atomic.LoadInt32(&calls) != 0 {
		t.Fatalf("expected no reload calls for unchanged file, got %d", calls)
	}
}
