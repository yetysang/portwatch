package alert

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// ThrottleHandler wraps a Handler and suppresses repeated identical change sets
// within a configurable quiet window. Unlike RateLimitedHandler (which operates
// per-change-key), ThrottleHandler silences the entire batch if the full set of
// keys was already dispatched recently.
type ThrottleHandler struct {
	mu      sync.Mutex
	inner   Handler
	window  time.Duration
	lastKey string
	lastAt  time.Time
	now     func() time.Time
}

// NewThrottleHandler returns a ThrottleHandler that forwards to inner at most
// once per window for an identical batch of changes.
func NewThrottleHandler(inner Handler, window time.Duration) *ThrottleHandler {
	return &ThrottleHandler{
		inner:  inner,
		window: window,
		now:    time.Now,
	}
}

// Handle forwards changes to the inner handler unless the same batch was
// already forwarded within the throttle window.
func (t *ThrottleHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	key := batchKey(changes)
	now := t.now()

	t.mu.Lock()
	if key == t.lastKey && now.Sub(t.lastAt) < t.window {
		t.mu.Unlock()
		return nil
	}
	t.lastKey = key
	t.lastAt = now
	t.mu.Unlock()

	return t.inner.Handle(changes)
}

// batchKey produces a stable string key representing the full set of changes.
func batchKey(changes []monitor.Change) string {
	var b []byte
	for _, c := range changes {
		b = append(b, changeKey(c)...)
		b = append(b, '|')
	}
	return string(b)
}
