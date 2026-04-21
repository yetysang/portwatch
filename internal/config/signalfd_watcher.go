package config

import "time"

// WatcherConfig controls the file-based config reload watcher.
type WatcherConfig struct {
	// Enabled turns on inotify/poll-based config file watching.
	Enabled bool `toml:"enabled" yaml:"enabled"`

	// PollInterval is the fallback polling interval when inotify is unavailable.
	PollInterval time.Duration `toml:"poll_interval" yaml:"poll_interval"`

	// DebounceDelay is the quiet period after the last change event before
	// the reload is triggered. Prevents rapid successive reloads.
	DebounceDelay time.Duration `toml:"debounce_delay" yaml:"debounce_delay"`
}

// DefaultWatcherConfig returns a WatcherConfig with sensible defaults.
func DefaultWatcherConfig() WatcherConfig {
	return WatcherConfig{
		Enabled:       false,
		PollInterval:  5 * time.Second,
		DebounceDelay: 500 * time.Millisecond,
	}
}

// Validate checks WatcherConfig for logical errors.
func (w WatcherConfig) Validate() error {
	if !w.Enabled {
		return nil
	}
	if w.PollInterval < 500*time.Millisecond {
		return newValidationError("watcher.poll_interval", "must be at least 500ms")
	}
	if w.PollInterval > 5*time.Minute {
		return newValidationError("watcher.poll_interval", "must not exceed 5 minutes")
	}
	if w.DebounceDelay < 0 {
		return newValidationError("watcher.debounce_delay", "must not be negative")
	}
	if w.DebounceDelay > 30*time.Second {
		return newValidationError("watcher.debounce_delay", "must not exceed 30 seconds")
	}
	return nil
}
