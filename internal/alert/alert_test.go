package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func makeChange(changeType, proto, addr string, port int) monitor.Change {
	return monitor.Change{
		Type: changeType,
		Binding: ports.Binding{
			Proto:     proto,
			LocalAddr: addr,
			LocalPort: port,
		},
	}
}

func TestHandler_InfoLevel(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, nil)

	a := h.Handle(makeChange("added", "tcp", "0.0.0.0", 8080))

	if a.Level != LevelInfo {
		t.Errorf("expected INFO, got %s", a.Level)
	}
	if !strings.Contains(buf.String(), "added") {
		t.Errorf("output missing 'added': %q", buf.String())
	}
}

func TestHandler_WarnLevel(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, []int{22, 3389})

	a := h.Handle(makeChange("added", "tcp", "0.0.0.0", 22))

	if a.Level != LevelWarn {
		t.Errorf("expected WARN, got %s", a.Level)
	}
	if !strings.Contains(buf.String(), "WARN") {
		t.Errorf("output missing WARN: %q", buf.String())
	}
}

// TestHandler_NonWatchedPort verifies that a port not in the watch list
// produces an INFO-level alert rather than a WARN-level alert.
func TestHandler_NonWatchedPort(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, []int{22, 3389})

	a := h.Handle(makeChange("added", "tcp", "0.0.0.0", 8080))

	if a.Level != LevelInfo {
		t.Errorf("expected INFO for non-watched port, got %s", a.Level)
	}
}

func TestHandler_OutputFormat(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, nil)
	h.Handle(makeChange("removed", "udp", "127.0.0.1", 53))

	line := buf.String()
	for _, want := range []string{"udp", "53", "removed", "127.0.0.1"} {
		if !strings.Contains(line, want) {
			t.Errorf("output missing %q: %s", want, line)
		}
	}
}

func TestHandler_Drain(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, nil)

	ch := make(chan monitor.Change, 3)
	ch <- makeChange("added", "tcp", "0.0.0.0", 80)
	ch <- makeChange("added", "tcp", "0.0.0.0", 443)
	ch <- makeChange("removed", "tcp", "0.0.0.0", 8080)
	close(ch)

	stop := make(chan struct{})
	h.Drain(ch, stop)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 alert lines, got %d", len(lines))
	}
}
