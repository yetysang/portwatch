// Package daemon contains runtime lifecycle helpers for portwatch.
package daemon

import (
	"context"
	"crypto/sha256"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/rorycl/portwatch/internal/config"
)

// ReloadFunc is called when a config file change is detected.
type ReloadFunc func(path string) error

// ConfigWatcher polls a config file for changes and triggers a reload
// callback when the file content changes.
type ConfigWatcher struct {
	cfg    config.WatcherConfig
	path   string
	log    *slog.Logger
	onLoad ReloadFunc
}

// NewConfigWatcher creates a ConfigWatcher. It does not start polling.
func NewConfigWatcher(cfg config.WatcherConfig, path string, log *slog.Logger, fn ReloadFunc) *ConfigWatcher {
	return &ConfigWatcher{
		cfg:    cfg,
		path:   path,
		log:    log,
		onLoad: fn,
	}
}

// Run starts the polling loop. It blocks until ctx is cancelled.
func (w *ConfigWatcher) Run(ctx context.Context) {
	if !w.cfg.Enabled {
		w.log.Debug("config watcher disabled")
		return
	}
	ticker := time.NewTicker(w.cfg.PollInterval)
	defer ticker.Stop()

	last := w.hashFile()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			current := w.hashFile()
			if current == "" || current == last {
				continue
			}
			w.log.Info("config file changed, debouncing", "path", w.path, "delay", w.cfg.DebounceDelay)
			select {
			case <-ctx.Done():
				return
			case <-time.After(w.cfg.DebounceDelay):
			}
			if err := w.onLoad(w.path); err != nil {
				w.log.Error("config reload failed", "err", err)
			} else {
				w.log.Info("config reloaded", "path", w.path)
				last = current
			}
		}
	}
}

func (w *ConfigWatcher) hashFile() string {
	f, err := os.Open(w.path)
	if err != nil {
		return ""
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return string(h.Sum(nil))
}
