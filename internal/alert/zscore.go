package alert

import (
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

// ZScoreHandler wraps another Handler and suppresses forwarding unless the
// number of changes in the current tick is statistically anomalous relative
// to a rolling window of historical tick sizes.
type ZScoreHandler struct {
	mu       sync.Mutex
	cfg      config.ZScoreConfig
	next     Handler
	window   []float64
	lastFire time.Time
	now      func() time.Time
}

// NewZScoreHandler constructs a ZScoreHandler that forwards to next only when
// the change count exceeds cfg.Threshold standard deviations above the mean.
func NewZScoreHandler(cfg config.ZScoreConfig, next Handler) *ZScoreHandler {
	return &ZScoreHandler{
		cfg:  cfg,
		next: next,
		now:  time.Now,
	}
}

// Handle records the tick's change count and forwards to the wrapped handler
// only when anomalous activity is detected.
func (h *ZScoreHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	count := float64(len(changes))
	h.record(count)

	if !h.isAnomalous(count) {
		return nil
	}

	now := h.now()
	if now.Sub(h.lastFire) < h.cfg.Cooldown {
		return nil
	}
	h.lastFire = now

	slog.Info("zscore anomaly detected",
		"change_count", count,
		"mean", fmt.Sprintf("%.2f", h.mean()),
		"stddev", fmt.Sprintf("%.2f", h.stddev()),
	)
	return h.next.Handle(changes)
}

// Drain flushes the wrapped handler.
func (h *ZScoreHandler) Drain() error { return h.next.Drain() }

func (h *ZScoreHandler) record(v float64) {
	h.window = append(h.window, v)
	if len(h.window) > h.cfg.WindowSize {
		h.window = h.window[len(h.window)-h.cfg.WindowSize:]
	}
}

func (h *ZScoreHandler) isAnomalous(v float64) bool {
	if len(h.window) < h.cfg.MinSamples {
		return false
	}
	std := h.stddev()
	if std == 0 {
		return false
	}
	z := (v - h.mean()) / std
	return z >= h.cfg.Threshold
}

func (h *ZScoreHandler) mean() float64 {
	if len(h.window) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range h.window {
		sum += v
	}
	return sum / float64(len(h.window))
}

func (h *ZScoreHandler) stddev() float64 {
	if len(h.window) < 2 {
		return 0
	}
	m := h.mean()
	ss := 0.0
	for _, v := range h.window {
		d := v - m
		ss += d * d
	}
	return math.Sqrt(ss / float64(len(h.window)))
}
