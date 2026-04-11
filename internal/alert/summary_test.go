package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func summaryChange(ip string, port int, added bool) monitor.Change {
	return monitor.Change{
		Binding: ports.Binding{IP: ip, Port: port},
		Added:   added,
	}
}

func fixedNow(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestSummaryHandler_EmptyChangesNoOutput(t *testing.T) {
	var buf bytes.Buffer
	h := NewSummaryHandler(&buf, "", time.Minute)
	_ = h.Handle(nil)
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestSummaryHandler_AccumulatesWithinInterval(t *testing.T) {
	var buf bytes.Buffer
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	h := NewSummaryHandler(&buf, "", time.Hour)
	h.now = fixedNow(now)
	h.lastAt = now // pretend we just flushed

	_ = h.Handle([]monitor.Change{summaryChange("127.0.0.1", 8080, true)})
	if buf.Len() != 0 {
		t.Fatalf("expected no flush yet, got %q", buf.String())
	}
	if len(h.buf) != 1 {
		t.Fatalf("expected 1 buffered change, got %d", len(h.buf))
	}
}

func TestSummaryHandler_FlushesAfterInterval(t *testing.T) {
	var buf bytes.Buffer
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	h := NewSummaryHandler(&buf, "[pw] ", time.Minute)
	h.now = fixedNow(base)
	h.lastAt = base.Add(-2 * time.Minute)

	_ = h.Handle([]monitor.Change{
		summaryChange("0.0.0.0", 443, true),
		summaryChange("0.0.0.0", 80, false),
	})
	out := buf.String()
	if !strings.Contains(out, "[pw] [summary]") {
		t.Errorf("expected prefix in output, got %q", out)
	}
	if !strings.Contains(out, "added") {
		t.Errorf("expected 'added' section, got %q", out)
	}
	if !strings.Contains(out, "removed") {
		t.Errorf("expected 'removed' section, got %q", out)
	}
}

func TestSummaryHandler_FlushClearsBuffer(t *testing.T) {
	var buf bytes.Buffer
	h := NewSummaryHandler(&buf, "", time.Minute)
	h.now = fixedNow(time.Now())
	h.buf = []monitor.Change{summaryChange("127.0.0.1", 9090, true)}

	_ = h.Flush()
	if len(h.buf) != 0 {
		t.Fatalf("expected buffer cleared after flush, got %d entries", len(h.buf))
	}
}

func TestSummaryHandler_FlushEmptyIsNoop(t *testing.T) {
	var buf bytes.Buffer
	h := NewSummaryHandler(&buf, "", time.Minute)
	_ = h.Flush()
	if buf.Len() != 0 {
		t.Fatalf("expected no output for empty flush, got %q", buf.String())
	}
}
