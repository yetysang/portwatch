package config

import "fmt"

// SNSConfig holds AWS SNS alerting configuration.
type SNSConfig struct {
	Enabled  bool   `toml:"enabled"`
	TopicARN string `toml:"topic_arn"`
	Region   string `toml:"region"`
}

// DefaultSNSConfig returns an SNSConfig with sensible defaults.
func DefaultSNSConfig() SNSConfig {
	return SNSConfig{
		Enabled:  false,
		TopicARN: "",
		Region:   "us-east-1",
	}
}

// Validate checks that the SNSConfig is coherent.
func (c SNSConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.TopicARN == "" {
		return fmt.Errorf("sns: topic_arn is required when enabled")
	}
	if c.Region == "" {
		return fmt.Errorf("sns: region is required when enabled")
	}
	return nil
}
