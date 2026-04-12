package config

import "fmt"

// LokiConfig holds Grafana Loki alerting configuration.
type LokiConfig struct {
	Enabled  bool   `toml:"enabled"`
	URL      string `toml:"url"`
	JobLabel string `toml:"job_label"`
}

// DefaultLokiConfig returns a LokiConfig with sensible defaults.
func DefaultLokiConfig() LokiConfig {
	return LokiConfig{
		Enabled:  false,
		URL:      "http://localhost:3100",
		JobLabel: "portwatch",
	}
}

// Validate checks the LokiConfig for correctness.
func (c LokiConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("loki: url must not be empty when enabled")
	}
	if c.JobLabel == "" {
		return fmt.Errorf("loki: job_label must not be empty when enabled")
	}
	return nil
}
