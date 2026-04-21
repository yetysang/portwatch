package config

import (
	"fmt"
	"time"
)

// CloudWatchConfig holds settings for the AWS CloudWatch Logs alert handler.
type CloudWatchConfig struct {
	Enabled       bool          `toml:"enabled"`
	Region        string        `toml:"region"`
	LogGroup      string        `toml:"log_group"`
	LogStream     string        `toml:"log_stream"`
	AccessKeyID   string        `toml:"access_key_id"`
	SecretKey     string        `toml:"secret_access_key"`
	FlushInterval time.Duration `toml:"flush_interval"`
}

// DefaultCloudWatchConfig returns a CloudWatchConfig with sensible defaults.
func DefaultCloudWatchConfig() CloudWatchConfig {
	return CloudWatchConfig{
		Enabled:       false,
		Region:        "us-east-1",
		LogGroup:      "/portwatch/alerts",
		LogStream:     "portwatch",
		FlushInterval: 5 * time.Second,
	}
}

// Validate returns an error if the CloudWatchConfig is invalid.
func (c CloudWatchConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Region == "" {
		return fmt.Errorf("cloudwatch: region is required")
	}
	if c.LogGroup == "" {
		return fmt.Errorf("cloudwatch: log_group is required")
	}
	if c.LogStream == "" {
		return fmt.Errorf("cloudwatch: log_stream is required")
	}
	if c.FlushInterval < time.Second {
		return fmt.Errorf("cloudwatch: flush_interval must be at least 1s")
	}
	return nil
}
