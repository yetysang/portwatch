// Package alert provides change notification handlers for portwatch.
package alert

import (
	"fmt"

	"github.com/example/portwatch/internal/monitor"
	"github.com/example/portwatch/internal/ports"
)

// RateLimitedHandler wraps a Handler and suppresses repeated alerts for the
// same port/protocol pair within a configurable cooldown window.
type RateLimitedHandler struct {
	inner *Handler
	rl    *ports.RateLimiter
}

// NewRateLimitedHandler wraps h with rate-limiting provided by rl.
func NewRateLimitedHandler(h *Handler, rl *ports.RateLimiter) *RateLimitedHandler {
	return &RateLimitedHandler{inner: h, rl: rl}
}

// Handle forwards only those changes whose key has not been seen within the
// cooldown window. Suppressed changes are silently dropped.
func (r *RateLimitedHandler) Handle(changes []monitor.Change) {
	allowed := changes[:0:0]
	for _, c := range changes {
		key := changeKey(c)
		if r.rl.Allow(key) {
			allowed = append(allowed, c)
		}
	}
	if len(allowed) > 0 {
		r.inner.Handle(allowed)
	}
}

// Drain flushes all buffered output from the inner handler.
func (r *RateLimitedHandler) Drain() []string {
	return r.inner.Drain()
}

// changeKey produces a stable string key for a monitor.Change.
func changeKey(c monitor.Change) string {
	return fmt.Sprintf("%s:%s:%d", c.Kind, c.Binding.Proto, c.Binding.Port)
}
