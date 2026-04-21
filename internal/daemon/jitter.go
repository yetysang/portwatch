// Package daemon provides runtime helpers for the portwatch daemon process.
package daemon

import (
	"math/rand"
	"time"

	"github.com/user/portwatch/internal/config"
)

// Jitter returns a random duration in [0, cfg.MaxJitter) when jitter is
// enabled, or zero otherwise. It is safe to call from multiple goroutines.
func Jitter(cfg config.JitterConfig) time.Duration {
	if !cfg.Enabled || cfg.MaxJitter <= 0 {
		return 0
	}
	//nolint:gosec // non-cryptographic jitter is intentional
	return time.Duration(rand.Int63n(int64(cfg.MaxJitter)))
}

// SleepWithJitter blocks for the base interval plus a random jitter drawn
// from cfg. It returns immediately when ctx is done.
func SleepWithJitter(base time.Duration, cfg config.JitterConfig) time.Duration {
	j := Jitter(cfg)
	total := base + j
	time.Sleep(total)
	return total
}

// TickerWithJitter returns a channel that receives after each (interval +
// random jitter) period. The caller is responsible for draining the channel.
// The returned stop function must be called to release resources.
func TickerWithJitter(interval time.Duration, cfg config.JitterConfig) (<-chan time.Time, func()) {
	ch := make(chan time.Time, 1)
	stop := make(chan struct{})

	go func() {
		for {
			wait := interval + Jitter(cfg)
			select {
			case <-time.After(wait):
				select {
				case ch <- time.Now():
				default:
				}
			case <-stop:
				close(ch)
				return
			}
		}
	}()

	return ch, func() { close(stop) }
}
