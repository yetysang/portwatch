package config

import (
	"testing"
	"time"
)

func TestDefaultCloudWatchConfig_Values(t *testing.T) {
	cfg := DefaultCloudWatchConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false")
	}
	if cfg.Region != "us-east-1" {
		t.Errorf("unexpected Region: %s", cfg.Region)
	}
	if cfg.LogGroup != "/portwatch/alerts" {
		t.Errorf("unexpected LogGroup: %s", cfg.LogGroup)
	}
	if cfg.LogStream != "portwatch" {
		t.Errorf("unexpected LogStream: %s", cfg.LogStream)
	}
	if cfg.FlushInterval != 5*time.Second {
		t.Errorf("unexpected FlushInterval: %v", cfg.FlushInterval)
	}
}

func TestCloudWatchConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := CloudWatchConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestCloudWatchConfig_ValidateEnabledRequiresRegion(t *testing.T) {
	cfg := DefaultCloudWatchConfig()
	cfg.Enabled = true
	cfg.Region = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing region")
	}
}

func TestCloudWatchConfig_ValidateEnabledRequiresLogGroup(t *testing.T) {
	cfg := DefaultCloudWatchConfig()
	cfg.Enabled = true
	cfg.LogGroup = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing log_group")
	}
}

func TestCloudWatchConfig_ValidateEnabledRequiresLogStream(t *testing.T) {
	cfg := DefaultCloudWatchConfig()
	cfg.Enabled = true
	cfg.LogStream = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing log_stream")
	}
}

func TestCloudWatchConfig_ValidateFlushIntervalTooShort(t *testing.T) {
	cfg := DefaultCloudWatchConfig()
	cfg.Enabled = true
	cfg.FlushInterval = 100 * time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for flush_interval below 1s")
	}
}

func TestCloudWatchConfig_ValidateEnabledWithValidFields(t *testing.T) {
	cfg := DefaultCloudWatchConfig()
	cfg.Enabled = true
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
