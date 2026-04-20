package alert

import (
	"fmt"
	"sync"
	"time"
)

// CircuitState represents the state of the circuit breaker.
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// CircuitBreakerHandler wraps an alert Handler with a circuit breaker that
// stops forwarding events when the downstream handler fails repeatedly.
type CircuitBreakerHandler struct {
	inner      Handler
	mu         sync.Mutex
	state      CircuitState
	failures   int
	threshold  int
	resetAfter time.Duration
	openedAt   time.Time
	now        func() time.Time
}

// NewCircuitBreakerHandler returns a Handler that trips open after
// threshold consecutive errors from inner, and attempts recovery after
// resetAfter has elapsed.
func NewCircuitBreakerHandler(inner Handler, threshold int, resetAfter time.Duration) *CircuitBreakerHandler {
	if threshold <= 0 {
		threshold = 3
	}
	if resetAfter <= 0 {
		resetAfter = 30 * time.Second
	}
	return &CircuitBreakerHandler{
		inner:      inner,
		state:      CircuitClosed,
		threshold:  threshold,
		resetAfter: resetAfter,
		now:        time.Now,
	}
}

// Handle forwards changes to the inner handler unless the circuit is open.
func (c *CircuitBreakerHandler) Handle(changes []Change) error {
	if len(changes) == 0 {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	switch c.state {
	case CircuitOpen:
		if c.now().Sub(c.openedAt) < c.resetAfter {
			return fmt.Errorf("circuit open: downstream handler unavailable")
		}
		c.state = CircuitHalfOpen
	case CircuitHalfOpen:
		// allow a single probe through
	}

	err := c.inner.Handle(changes)
	if err != nil {
		c.failures++
		if c.state == CircuitHalfOpen || c.failures >= c.threshold {
			c.state = CircuitOpen
			c.openedAt = c.now()
		}
		return err
	}

	// success — reset
	c.failures = 0
	c.state = CircuitClosed
	return nil
}

// Drain is a no-op passthrough.
func (c *CircuitBreakerHandler) Drain() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.inner.Drain()
}

// State returns the current circuit state (for observability / testing).
func (c *CircuitBreakerHandler) State() CircuitState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}
