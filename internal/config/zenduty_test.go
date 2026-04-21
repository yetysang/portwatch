package config

import "testing"

func TestDefaultZendutyConfig_Values(t *testing.T) {
	cfg := DefaultZendutyConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false")
	}
	if cfg.AlertType != "warning" {
		t.Errorf("expected AlertType=warning, got %q", cfg.AlertType)
	}
}

func TestZendutyConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := DefaultZendutyConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestZendutyConfig_ValidateEnabledRequiresAPIKey(t *testing.T) {
	cfg := DefaultZendutyConfig()
	cfg.Enabled = true
	cfg.ServiceID = "svc-123"
	cfg.IntegrationID = "int-456"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing api_key")
	}
}

func TestZendutyConfig_ValidateEnabledRequiresServiceID(t *testing.T) {
	cfg := DefaultZendutyConfig()
	cfg.Enabled = true
	cfg.APIKey = "key-abc"
	cfg.IntegrationID = "int-456"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing service_id")
	}
}

func TestZendutyConfig_ValidateEnabledRequiresIntegrationID(t *testing.T) {
	cfg := DefaultZendutyConfig()
	cfg.Enabled = true
	cfg.APIKey = "key-abc"
	cfg.ServiceID = "svc-123"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing integration_id")
	}
}

func TestZendutyConfig_ValidateInvalidAlertType(t *testing.T) {
	cfg := DefaultZendutyConfig()
	cfg.Enabled = true
	cfg.APIKey = "key-abc"
	cfg.ServiceID = "svc-123"
	cfg.IntegrationID = "int-456"
	cfg.AlertType = "unknown"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid alert_type")
	}
}

func TestZendutyConfig_ValidateEnabledWithValidFields(t *testing.T) {
	cfg := DefaultZendutyConfig()
	cfg.Enabled = true
	cfg.APIKey = "key-abc"
	cfg.ServiceID = "svc-123"
	cfg.IntegrationID = "int-456"
	cfg.AlertType = "critical"
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
