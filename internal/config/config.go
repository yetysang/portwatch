package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds all runtime configuration for portwatch.
type Config struct {
	// Interval between port scans.
	Interval time.Duration `json:"interval"`

	// LogLevel controls verbosity: "info", "warn", "error".
	LogLevel string `json:"log_level"`

	// IgnorePorts lists port numbers to suppress from alerts.
	IgnorePorts []int `json:"ignore_ports"`

	// WebhookURL is an optional HTTP endpoint to receive change events.
	// Leave empty to disable webhook delivery.
	WebhookURL string `json:"webhook_url"`

	// WebhookTimeout is the HTTP client timeout for webhook posts.
	WebhookTimeout time.Duration `json:"webhook_timeout"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:       15 * time.Second,
		LogLevel:       "info",
		IgnorePorts:    []int{},
		WebhookURL:     "",
		WebhookTimeout: 5 * time.Second,
	}
}

// Load reads a JSON config file from path and merges it over defaults.
// If path is empty, the default config is returned unchanged.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
