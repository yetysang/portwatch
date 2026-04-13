package alert

import (
	"runtime"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func execChange(kind monitor.ChangeKind, port uint16) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "127.0.0.1",
			Hostname: "localhost",
			Port:     port,
			Proto:    "TCP",
		},
	}
}

func TestExecHandler_EmptyChangesNoRun(t *testing.T) {
	h := NewExecHandler("false") // would fail if run
	if err := h.Handle(nil); err != nil {
		t.Fatalf("expected no error on empty changes, got %v", err)
	}
}

func TestExecHandler_RunsCommandOnChange(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}
	h := NewExecHandler("true")
	changes := []monitor.Change{execChange(monitor.Added, 8080)}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecHandler_ErrorOnFailingCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}
	h := NewExecHandler("false")
	changes := []monitor.Change{execChange(monitor.Added, 9090)}
	err := h.Handle(changes)
	if err == nil {
		t.Fatal("expected error from failing command")
	}
}

func TestFormatExecMsg_Added(t *testing.T) {
	changes := []monitor.Change{execChange(monitor.Added, 443)}
	msg := formatExecMsg(changes)
	if !strings.Contains(msg, "localhost:443/tcp added") {
		t.Errorf("unexpected message: %q", msg)
	}
}

func TestFormatExecMsg_Removed(t *testing.T) {
	changes := []monitor.Change{execChange(monitor.Removed, 80)}
	msg := formatExecMsg(changes)
	if !strings.Contains(msg, "removed") {
		t.Errorf("expected 'removed' in message: %q", msg)
	}
}

func TestFormatExecMsg_FallsBackToIP(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{Addr: "10.0.0.1", Port: 22, Proto: "TCP"},
	}
	msg := formatExecMsg([]monitor.Change{c})
	if !strings.Contains(msg, "10.0.0.1:22") {
		t.Errorf("expected IP fallback in message: %q", msg)
	}
}

func TestFormatExecMsg_MultipleChanges(t *testing.T) {
	changes := []monitor.Change{
		execChange(monitor.Added, 80),
		execChange(monitor.Removed, 443),
	}
	msg := formatExecMsg(changes)
	if !strings.Contains(msg, ";") {
		t.Errorf("expected semicolon separator in message: %q", msg)
	}
}

func TestExecHandler_DrainIsNoop(t *testing.T) {
	h := NewExecHandler("true")
	if err := h.Drain(); err != nil {
		t.Fatalf("Drain should be noop, got: %v", err)
	}
}
