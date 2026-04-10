package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all runtime configuration for portwatch.
type Config struct {
	// Interval between port scans.
	Interval time.Duration `yaml:"interval"`

	// LogLevel controls verbosity: "info" or "warn".
	LogLevel string `yaml:"log_level"`

	// IgnorePorts lists port numbers that should never trigger alerts.
	IgnorePorts []int `yaml:"ignore_ports"`

	// SnapshotPath is the file used to persist the last known state.
	SnapshotPath string `yaml:"snapshot_path"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:     15 * time.Second,
		LogLevel:     "info",
		IgnorePorts:  []int{},
		SnapshotPath: "/var/lib/portwatch/snapshot.json",
	}
}

// Load reads a YAML config file from path and merges it over the defaults.
// If path is empty, the default config is returned unchanged.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
