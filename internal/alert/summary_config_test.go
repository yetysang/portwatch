package alert

import (
	"testing"
	"time"
)

func TestDefaultSummaryConfig_Values(t *testing.T) {
	cfg := DefaultSummaryConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Interval != 5*time.Minute {
		t.Errorf("expected 5m interval, got %v", cfg.Interval)
	}
	if cfg.Prefix == "" {
		t.Error("expected non-empty default prefix")
	}
}

func TestSummaryConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := SummaryConfig{Enabled: false, Interval: 0}
	if err := cfg.Validate(); err != nil {
		t.Errorf("disabled config should always be valid, got %v", err)
	}
}

func TestSummaryConfig_ValidateIntervalTooShort(t *testing.T) {
	cfg := SummaryConfig{Enabled: true, Interval: 5 * time.Second}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for interval < 10s")
	}
}

func TestSummaryConfig_ValidateIntervalTooLong(t *testing.T) {
	cfg := SummaryConfig{Enabled: true, Interval: 25 * time.Hour}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for interval > 24h")
	}
}

func TestSummaryConfig_ValidateValidConfig(t *testing.T) {
	cfg := SummaryConfig{Enabled: true, Interval: 10 * time.Minute, Prefix: "[pw] "}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected valid config, got %v", err)
	}
}
