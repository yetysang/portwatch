package alert

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// RetryConfig holds configuration for the retry handler.
type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

// DefaultRetryConfig returns a RetryConfig with sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// retryHandler wraps a Handler and retries on error with exponential backoff.
type retryHandler struct {
	next   Handler
	cfg    RetryConfig
	sleep  func(time.Duration)
}

// NewRetryHandler returns a Handler that retries the wrapped handler on failure.
func NewRetryHandler(next Handler, cfg RetryConfig) Handler {
	return &retryHandler{
		next:  next,
		cfg:   cfg,
		sleep: time.Sleep,
	}
}

func (r *retryHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	delay := r.cfg.Delay
	var lastErr error
	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		err := r.next.Handle(changes)
		if err == nil {
			return nil
		}
		lastErr = err
		log.Printf("[retry] attempt %d/%d failed: %v", attempt, r.cfg.MaxAttempts, err)
		if attempt < r.cfg.MaxAttempts {
			r.sleep(delay)
			delay = time.Duration(float64(delay) * r.cfg.Multiplier)
		}
	}
	return lastErr
}

func (r *retryHandler) Drain() error {
	return r.next.Drain()
}
