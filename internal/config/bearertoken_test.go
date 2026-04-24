package config

import "testing"

func TestDefaultBearerTokenConfig_Values(t *testing.T) {
	cfg := DefaultBearerTokenConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.Header != "Authorization" {
		t.Errorf("expected Header=Authorization, got %q", cfg.Header)
	}
	if cfg.Token != "" {
		t.Errorf("expected empty Token by default, got %q", cfg.Token)
	}
}

func TestBearerTokenConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := BearerTokenConfig{Enabled: false, Token: "", Header: ""}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestBearerTokenConfig_ValidateEnabledRequiresToken(t *testing.T) {
	cfg := BearerTokenConfig{Enabled: true, Token: "", Header: "Authorization"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when token is empty")
	}
}

func TestBearerTokenConfig_ValidateEnabledRequiresHeader(t *testing.T) {
	cfg := BearerTokenConfig{Enabled: true, Token: "secret", Header: ""}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when header is empty")
	}
}

func TestBearerTokenConfig_ValidateEnabledWithBothFields(t *testing.T) {
	cfg := BearerTokenConfig{Enabled: true, Token: "mysecret", Header: "Authorization"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error with valid config, got %v", err)
	}
}

func TestBearerTokenConfig_ValidateCustomHeader(t *testing.T) {
	cfg := BearerTokenConfig{Enabled: true, Token: "tok", Header: "X-API-Key"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error with custom header, got %v", err)
	}
}
