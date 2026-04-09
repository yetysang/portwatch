package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Interval != 5*time.Second {
		t.Errorf("expected 5s interval, got %v", cfg.Interval)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected log_level info, got %q", cfg.LogLevel)
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 5*time.Second {
		t.Errorf("expected default interval, got %v", cfg.Interval)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	p := writeTemp(t, "interval: 10s\nlog_level: warn\nignore_ports: [22, 80]\n")
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 10*time.Second {
		t.Errorf("expected 10s, got %v", cfg.Interval)
	}
	if cfg.LogLevel != "warn" {
		t.Errorf("expected warn, got %q", cfg.LogLevel)
	}
	if len(cfg.IgnorePorts) != 2 {
		t.Errorf("expected 2 ignore_ports, got %d", len(cfg.IgnorePorts))
	}
}

func TestLoad_InvalidInterval(t *testing.T) {
	p := writeTemp(t, "interval: notaduration\n")
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for invalid interval")
	}
}

func TestLoad_IntervalTooShort(t *testing.T) {
	p := writeTemp(t, "interval: 500ms\n")
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for interval < 1s")
	}
}

func TestLoad_InvalidLogLevel(t *testing.T) {
	p := writeTemp(t, "interval: 5s\nlog_level: verbose\n")
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for unknown log_level")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
