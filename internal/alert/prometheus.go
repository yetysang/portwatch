package alert

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/user/portwatch/internal/monitor"
)

// PrometheusHandler exposes port change metrics via an HTTP endpoint
// compatible with Prometheus scraping.
type PrometheusHandler struct {
	mu      sync.Mutex
	added   int64
	removed int64
	total   int64
}

// NewPrometheusHandler creates a PrometheusHandler. Register its ServeHTTP
// method on your metrics endpoint (e.g. /metrics).
func NewPrometheusHandler() *PrometheusHandler {
	return &PrometheusHandler{}
}

// Handle records counts for each change in the batch.
func (h *PrometheusHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, c := range changes {
		switch c.Kind {
		case monitor.Added:
			h.added++
		case monitor.Removed:
			h.removed++
		}
		h.total++
	}
	return nil
}

// Drain is a no-op for the Prometheus handler.
func (h *PrometheusHandler) Drain() error { return nil }

// ServeHTTP writes Prometheus text-format metrics.
func (h *PrometheusHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	h.mu.Lock()
	added, removed, total := h.added, h.removed, h.total
	h.mu.Unlock()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "# HELP portwatch_changes_total Total port binding changes observed.\n")
	fmt.Fprintf(w, "# TYPE portwatch_changes_total counter\n")
	fmt.Fprintf(w, "portwatch_changes_total{kind=\"added\"}   %d\n", added)
	fmt.Fprintf(w, "portwatch_changes_total{kind=\"removed\"} %d\n", removed)
	fmt.Fprintf(w, "portwatch_changes_total{kind=\"any\"}     %d\n", total)
}
