package alert

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// SummaryHandler periodically emits a human-readable summary of all changes
// seen since the last flush, grouped by added/removed.
type SummaryHandler struct {
	out     io.Writer
	prefix  string
	buf     []monitor.Change
	lastAt  time.Time
	interval time.Duration
	now     func() time.Time
}

// NewSummaryHandler returns a SummaryHandler that writes to w.
// interval controls how often the accumulated summary is flushed.
func NewSummaryHandler(w io.Writer, prefix string, interval time.Duration) *SummaryHandler {
	if w == nil {
		w = os.Stdout
	}
	return &SummaryHandler{
		out:      w,
		prefix:   prefix,
		interval: interval,
		now:      time.Now,
	}
}

// Handle accumulates changes and flushes a summary when the interval elapses.
func (h *SummaryHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	h.buf = append(h.buf, changes...)
	if h.now().Sub(h.lastAt) >= h.interval {
		return h.Flush()
	}
	return nil
}

// Flush writes the accumulated summary immediately and resets the buffer.
func (h *SummaryHandler) Flush() error {
	if len(h.buf) == 0 {
		return nil
	}
	var added, removed []string
	for _, c := range h.buf {
		desc := fmt.Sprintf("%s:%d", c.Binding.IP, c.Binding.Port)
		if c.Added {
			added = append(added, desc)
		} else {
			removed = append(removed, desc)
		}
	}
	lines := []string{fmt.Sprintf("%s[summary] %s", h.prefix, h.now().Format(time.RFC3339))}
	if len(added) > 0 {
		lines = append(lines, fmt.Sprintf("  added   (%d): %s", len(added), strings.Join(added, ", ")))
	}
	if len(removed) > 0 {
		lines = append(lines, fmt.Sprintf("  removed (%d): %s", len(removed), strings.Join(removed, ", ")))
	}
	_, err := fmt.Fprintln(h.out, strings.Join(lines, "\n"))
	h.buf = h.buf[:0]
	h.lastAt = h.now()
	return err
}
