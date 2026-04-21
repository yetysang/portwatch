package config

import (
	"testing"
	"time"
)

func TestDefaultInfluxDBConfig_Values(t *testing.T) {
	cfg := DefaultInfluxDBConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.Measurement != "portwatch_events" {
		t.Errorf("unexpected default measurement: %s", cfg.Measurement)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("unexpected default timeout: %v", cfg.Timeout)
	}
}

func TestInfluxDBConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := DefaultInfluxDBConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got: %v", err)
	}
}

func TestInfluxDBConfig_ValidateEnabledRequiresURL(t *testing.T) {
	cfg := DefaultInfluxDBConfig()
	cfg.Enabled = true
	cfg.Token = "tok"
	cfg.Org = "myorg"
	cfg.Bucket = "mybucket"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when url is missing")
	}
}

func TestInfluxDBConfig_ValidateEnabledRequiresToken(t *testing.T) {
	cfg := DefaultInfluxDBConfig()
	cfg.Enabled = true
	cfg.URL = "http://localhost:8086"
	cfg.Org = "myorg"
	cfg.Bucket = "mybucket"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when token is missing")
	}
}

func TestInfluxDBConfig_ValidateEnabledRequiresOrgAndBucket(t *testing.T) {
	cfg := DefaultInfluxDBConfig()
	cfg.Enabled = true
	cfg.URL = "http://localhost:8086"
	cfg.Token = "tok"
	cfg.Org = "myorg"
	// bucket missing
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when bucket is missing")
	}
}

func TestInfluxDBConfig_ValidateEnabledWithAllFields(t *testing.T) {
	cfg := DefaultInfluxDBConfig()
	cfg.Enabled = true
	cfg.URL = "http://localhost:8086"
	cfg.Token = "tok"
	cfg.Org = "myorg"
	cfg.Bucket = "mybucket"
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error with all fields set, got: %v", err)
	}
}
