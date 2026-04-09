package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	Interval    time.Duration `yaml:"-"`
	IntervalRaw string        `yaml:"interval"`
	LogLevel    string        `yaml:"log_level"`
	AlertFile   string        `yaml:"alert_file"`
	IgnorePorts []int         `yaml:"ignore_ports"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		IntervalRaw: "5s",
		Interval:    5 * time.Second,
		LogLevel:    "info",
		AlertFile:   "",
		IgnorePorts: []int{},
	}
}

// Load reads a YAML config file from path and returns a validated Config.
// If path is empty, the default config is returned.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.IntervalRaw != "" {
		d, err := time.ParseDuration(c.IntervalRaw)
		if err != nil {
			return fmt.Errorf("config: invalid interval %q: %w", c.IntervalRaw, err)
		}
		if d < time.Second {
			return fmt.Errorf("config: interval %q is too short (minimum 1s)", c.IntervalRaw)
		}
		c.Interval = d
	}

	switch c.LogLevel {
	case "info", "warn", "error", "debug":
	default:
		return fmt.Errorf("config: unknown log_level %q", c.LogLevel)
	}

	return nil
}
