package alert

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func prometheusChange(kind monitor.ChangeKind, port uint16) monitor.Change {
	return monitor.Change{
		Kind:    kind,
		Binding: ports.Binding{Port: port, Proto: "tcp", Addr: "0.0.0.0"},
	}
}

func TestPrometheusHandler_EmptyChangesNoCount(t *testing.T) {
	h := NewPrometheusHandler()
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.total != 0 {
		t.Errorf("expected total=0, got %d", h.total)
	}
}

func TestPrometheusHandler_CountsAdded(t *testing.T) {
	h := NewPrometheusHandler()
	_ = h.Handle([]monitor.Change{
		prometheusChange(monitor.Added, 8080),
		prometheusChange(monitor.Added, 9090),
	})
	if h.added != 2 {
		t.Errorf("expected added=2, got %d", h.added)
	}
	if h.removed != 0 {
		t.Errorf("expected removed=0, got %d", h.removed)
	}
	if h.total != 2 {
		t.Errorf("expected total=2, got %d", h.total)
	}
}

func TestPrometheusHandler_CountsRemoved(t *testing.T) {
	h := NewPrometheusHandler()
	_ = h.Handle([]monitor.Change{
		prometheusChange(monitor.Removed, 8080),
	})
	if h.removed != 1 {
		t.Errorf("expected removed=1, got %d", h.removed)
	}
	if h.added != 0 {
		t.Errorf("expected added=0, got %d", h.added)
	}
}

func TestPrometheusHandler_ServeHTTP_ContainsMetrics(t *testing.T) {
	h := NewPrometheusHandler()
	_ = h.Handle([]monitor.Change{
		prometheusChange(monitor.Added, 80),
		prometheusChange(monitor.Removed, 443),
	})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	body := rec.Body.String()
	for _, want := range []string{
		"portwatch_changes_total{kind=\"added\"}   1",
		"portwatch_changes_total{kind=\"removed\"} 1",
		"portwatch_changes_total{kind=\"any\"}     2",
		"# TYPE portwatch_changes_total counter",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing expected metric line %q in:\n%s", want, body)
		}
	}
}

func TestPrometheusHandler_DrainIsNoop(t *testing.T) {
	h := NewPrometheusHandler()
	if err := h.Drain(); err != nil {
		t.Errorf("Drain() returned unexpected error: %v", err)
	}
}
