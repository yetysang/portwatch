package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func TestDefaultSplunkConfig_Values(t *testing.T) {
	cfg := config.DefaultSplunkConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.URL != "" {
		t.Errorf("expected empty URL, got %q", cfg.URL)
	}
	if cfg.Token != "" {
		t.Errorf("expected empty Token, got %q", cfg.Token)
	}
	if cfg.Index != "main" {
		t.Errorf("expected default Index 'main', got %q", cfg.Index)
	}
	if cfg.Source != "portwatch" {
		t.Errorf("expected default Source 'portwatch', got %q", cfg.Source)
	}
}

func TestSplunkConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := config.DefaultSplunkConfig()
	cfg.Enabled = false
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestSplunkConfig_ValidateEnabledRequiresURL(t *testing.T) {
	cfg := config.DefaultSplunkConfig()
	cfg.Enabled = true
	cfg.Token = "splunk-hec-token"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when URL is missing")
	}
}

func TestSplunkConfig_ValidateEnabledRequiresToken(t *testing.T) {
	cfg := config.DefaultSplunkConfig()
	cfg.Enabled = true
	cfg.URL = "http://splunk.example.com:8088"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when Token is missing")
	}
}

func TestSplunkConfig_ValidateEnabledWithBothFields(t *testing.T) {
	cfg := config.DefaultSplunkConfig()
	cfg.Enabled = true
	cfg.URL = "http://splunk.example.com:8088"
	cfg.Token = "splunk-hec-token"
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
