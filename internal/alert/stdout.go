// Package alert provides handlers for dispatching port change notifications.
package alert

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// StdoutHandler writes change events to an io.Writer in a human-readable format.
type StdoutHandler struct {
	out    io.Writer
	prefix string
}

// NewStdoutHandler returns a StdoutHandler that writes to w.
// If w is nil, os.Stdout is used.
func NewStdoutHandler(w io.Writer, prefix string) *StdoutHandler {
	if w == nil {
		w = os.Stdout
	}
	return &StdoutHandler{out: w, prefix: prefix}
}

// Handle writes each change to the configured writer.
func (h *StdoutHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	ts := time.Now().UTC().Format(time.RFC3339)
	var errs []string
	for _, c := range changes {
		line := h.formatChange(ts, c)
		if _, err := fmt.Fprintln(h.out, line); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("stdout handler write errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

func (h *StdoutHandler) formatChange(ts string, c monitor.Change) string {
	verb := "ADDED"
	if c.Removed {
		verb = "REMOVED"
	}
	addr := fmt.Sprintf("%s:%d", c.Binding.Addr, c.Binding.Port)
	proc := c.Binding.ProcessName
	if proc == "" {
		proc = "unknown"
	}
	prefix := h.prefix
	if prefix != "" {
		prefix = "[" + prefix + "] "
	}
	return fmt.Sprintf("%s%s [%s] %s (proto=%s pid=%d proc=%s)",
		prefix, ts, verb, addr, c.Binding.Proto, c.Binding.PID, proc)
}
