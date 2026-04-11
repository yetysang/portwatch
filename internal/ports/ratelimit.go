// Package ports provides utilities for scanning and analysing port bindings.
package ports

import (
	"sync"
	"time"
)

// RateLimiter suppresses repeated alerts for the same port within a cooldown
// window. It is safe for concurrent use.
type RateLimiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
	now      func() time.Time // injectable for testing
}

// NewRateLimiter returns a RateLimiter that silences duplicate events for the
// same key until cooldown has elapsed since the last allowed event.
func NewRateLimiter(cooldown time.Duration) *RateLimiter {
	return &RateLimiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the event identified by key should be forwarded.
// Subsequent calls with the same key within the cooldown window return false.
func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.now()
	if t, ok := r.last[key]; ok && now.Sub(t) < r.cooldown {
		return false
	}
	r.last[key] = now
	return true
}

// Reset removes all recorded events, allowing every key through on the next
// call to Allow.
func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.last = make(map[string]time.Time)
}

// Len returns the number of keys currently tracked.
func (r *RateLimiter) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.last)
}
