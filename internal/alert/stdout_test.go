package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func stdoutChange(removed bool) monitor.Change {
	return monitor.Change{
		Removed: removed,
		Binding: ports.Binding{
			Addr:        "127.0.0.1",
			Port:        8080,
			Proto:       "tcp",
			PID:         1234,
			ProcessName: "myapp",
		},
	}
}

func TestStdoutHandler_EmptyChangesNoOutput(t *testing.T) {
	var buf bytes.Buffer
	h := NewStdoutHandler(&buf, "")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestStdoutHandler_AddedChange(t *testing.T) {
	var buf bytes.Buffer
	h := NewStdoutHandler(&buf, "")
	if err := h.Handle([]monitor.Change{stdoutChange(false)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ADDED") {
		t.Errorf("expected ADDED in output, got: %s", out)
	}
	if !strings.Contains(out, "127.0.0.1:8080") {
		t.Errorf("expected address in output, got: %s", out)
	}
	if !strings.Contains(out, "myapp") {
		t.Errorf("expected process name in output, got: %s", out)
	}
}

func TestStdoutHandler_RemovedChange(t *testing.T) {
	var buf bytes.Buffer
	h := NewStdoutHandler(&buf, "")
	if err := h.Handle([]monitor.Change{stdoutChange(true)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "REMOVED") {
		t.Errorf("expected REMOVED in output, got: %s", out)
	}
}

func TestStdoutHandler_PrefixIncluded(t *testing.T) {
	var buf bytes.Buffer
	h := NewStdoutHandler(&buf, "portwatch")
	if err := h.Handle([]monitor.Change{stdoutChange(false)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[portwatch]") {
		t.Errorf("expected prefix in output, got: %s", buf.String())
	}
}

func TestStdoutHandler_UnknownProcess(t *testing.T) {
	var buf bytes.Buffer
	h := NewStdoutHandler(&buf, "")
	c := monitor.Change{
		Binding: ports.Binding{Addr: "0.0.0.0", Port: 443, Proto: "tcp"},
	}
	if err := h.Handle([]monitor.Change{c}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "proc=unknown") {
		t.Errorf("expected proc=unknown in output, got: %s", buf.String())
	}
}
