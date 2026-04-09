package config

import (
	"testing"
	"time"
)

func TestValidate_DefaultConfigIsValid(t *testing.T) {
	cfg := DefaultConfig()
	if err := Validate(cfg); err != nil {
		t.Fatalf("expected default config to be valid, got: %v", err)
	}
}

func TestValidate_IntervalTooShort(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Interval = 100 * time.Millisecond
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for short interval")
	}
	if !IsValidationError(err) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
}

func TestValidate_IntervalTooLong(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Interval = 48 * time.Hour
	if err := Validate(cfg); err == nil {
		t.Fatal("expected validation error for excessive interval")
	}
}

func TestValidate_BadLogLevel(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LogLevel = "verbose"
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for unknown log level")
	}
	if !IsValidationError(err) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
}

func TestValidate_InvalidIgnorePort(t *testing.T) {
	cfg := DefaultConfig()
	cfg.IgnorePorts = []int{0, 70000}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for out-of-range ports")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Errs) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(ve.Errs), ve.Errs)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Interval = 1 * time.Millisecond
	cfg.LogLevel = "bad"
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected multiple validation errors")
	}
	ve := err.(*ValidationError)
	if len(ve.Errs) < 2 {
		t.Fatalf("expected at least 2 errors, got %d", len(ve.Errs))
	}
}

func TestIsValidationError(t *testing.T) {
	if IsValidationError(nil) {
		t.Fatal("nil should not be a ValidationError")
	}
	cfg := DefaultConfig()
	cfg.LogLevel = "nope"
	err := Validate(cfg)
	if !IsValidationError(err) {
		t.Fatal("expected IsValidationError to return true")
	}
}
