package alert

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func scriptChange() monitor.Change {
	return monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{
			Proto:   "tcp",
			Addr:    "0.0.0.0",
			Port:    9090,
			Process: "myapp",
			PID:     1234,
		},
	}
}

func TestScriptHandler_EmptyChangesNoRun(t *testing.T) {
	cfg := config.DefaultScriptConfig()
	cfg.Enabled = true
	cfg.Path = "/nonexistent/script.sh"
	h := NewScriptHandler(cfg)
	if err := h.Handle(nil); err != nil {
		t.Errorf("unexpected error on empty changes: %v", err)
	}
}

func TestScriptHandler_RunsScriptOnChange(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell script test on Windows")
	}

	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "out.json")
	scriptPath := filepath.Join(tmpDir, "notify.sh")

	scriptContent := "#!/bin/sh\ncat > " + outFile + "\n"
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	cfg := config.DefaultScriptConfig()
	cfg.Enabled = true
	cfg.Path = scriptPath
	cfg.Timeout = 5 * time.Second
	h := NewScriptHandler(cfg)

	if err := h.Handle([]monitor.Change{scriptChange()}); err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if m["port"].(float64) != 9090 {
		t.Errorf("expected port 9090, got %v", m["port"])
	}
	if m["kind"] != "added" {
		t.Errorf("expected kind 'added', got %v", m["kind"])
	}
}

func TestScriptHandler_ErrorOnFailingScript(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell script test on Windows")
	}

	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "fail.sh")
	if err := os.WriteFile(scriptPath, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	cfg := config.DefaultScriptConfig()
	cfg.Enabled = true
	cfg.Path = scriptPath
	h := NewScriptHandler(cfg)

	if err := h.Handle([]monitor.Change{scriptChange()}); err == nil {
		t.Error("expected error from failing script")
	}
}

func TestBuildScriptPayload_Fields(t *testing.T) {
	c := scriptChange()
	data, err := buildScriptPayload(c)
	if err != nil {
		t.Fatalf("buildScriptPayload: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"kind", "proto", "addr", "port", "process", "pid"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing key %q in payload", key)
		}
	}
}
