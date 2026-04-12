package config

import "testing"

func TestDefaultLokiConfig_Values(t *testing.T) {
	cfg := DefaultLokiConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.URL == "" {
		t.Error("expected non-empty default URL")
	}
	if cfg.JobLabel == "" {
		t.Error("expected non-empty default JobLabel")
	}
}

func TestLokiConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := LokiConfig{Enabled: false, URL: "", JobLabel: ""}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got: %v", err)
	}
}

func TestLokiConfig_ValidateEnabledRequiresURL(t *testing.T) {
	cfg := LokiConfig{Enabled: true, URL: "", JobLabel: "portwatch"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when URL is empty")
	}
}

func TestLokiConfig_ValidateEnabledRequiresJobLabel(t *testing.T) {
	cfg := LokiConfig{Enabled: true, URL: "http://localhost:3100", JobLabel: ""}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when JobLabel is empty")
	}
}

func TestLokiConfig_ValidateEnabledWithBothFields(t *testing.T) {
	cfg := LokiConfig{Enabled: true, URL: "http://localhost:3100", JobLabel: "portwatch"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error with valid config, got: %v", err)
	}
}
