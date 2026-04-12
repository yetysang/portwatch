package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo Level = "INFO"
	LevelWarn Level = "WARN"
)

// Alert holds a formatted notification about a port change.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Change    monitor.Change
}

// Handler processes Changes and writes formatted alerts to a writer.
type Handler struct {
	out       io.Writer
	watchlist map[int]bool // ports that trigger WARN instead of INFO
}

// NewHandler creates an alert Handler writing to w.
// Ports in watchlist will be flagged at WARN level.
func NewHandler(w io.Writer, watchlist []int) *Handler {
	if w == nil {
		w = os.Stdout
	}
	wm := make(map[int]bool, len(watchlist))
	for _, p := range watchlist {
		wm[p] = true
	}
	return &Handler{out: w, watchlist: wm}
}

// Handle converts a Change into an Alert and writes it.
func (h *Handler) Handle(c monitor.Change) Alert {
	lvl := LevelInfo
	if h.watchlist[c.Binding.LocalPort] {
		lvl = LevelWarn
	}

	msg := fmt.Sprintf("[%s] port %s/%d (%s) %s",
		lvl,
		c.Binding.Proto,
		c.Binding.LocalPort,
		c.Binding.LocalAddr,
		c.Type,
	)

	a := Alert{
		Timestamp: time.Now().UTC(),
		Level:     lvl,
		Message:   msg,
		Change:    c,
	}

	fmt.Fprintf(h.out, "%s  %s\n", a.Timestamp.Format(time.RFC3339), msg)
	return a
}

// Drain reads from the Changes channel until it is closed or stop is signalled.
func (h *Handler) Drain(changes <-chan monitor.Change, stop <-chan struct{}) {
	for {
		select {
		case c, ok := <-changes:
			if !ok {
				return
			}
			h.Handle(c)
		case <-stop:
			return
		}
	}
}

// AddToWatchlist adds a port to the handler's watchlist so that it triggers
// WARN-level alerts. Calling this method is safe between Drain iterations but
// should not be called concurrently with Handle or Drain.
func (h *Handler) AddToWatchlist(port int) {
	h.watchlist[port] = true
}

// RemoveFromWatchlist removes a port from the watchlist, reverting its alerts
// to INFO level. No-op if the port was not previously watched.
func (h *Handler) RemoveFromWatchlist(port int) {
	delete(h.watchlist, port)
}
