package alert

import (
	"sync"
	"time"

	"github.com/patrickdappollonio/portwatch/internal/config"
	"github.com/patrickdappollonio/portwatch/internal/monitor"
)

// RollupHandler accumulates changes within a time window and flushes them as a
// single batch to the downstream handler, capped at MaxBatch entries.
type RollupHandler struct {
	cfg      config.RollupConfig
	next     Handler
	mu       sync.Mutex
	buf      []monitor.Change
	deadline time.Time
	now      func() time.Time
}

// NewRollupHandler wraps next with rollup behaviour described by cfg.
func NewRollupHandler(cfg config.RollupConfig, next Handler) *RollupHandler {
	return &RollupHandler{
		cfg:  cfg,
		next: next,
		now:  time.Now,
	}
}

// Handle buffers changes and flushes when the window expires or the batch is full.
func (h *RollupHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	if !h.cfg.Enabled {
		return h.next.Handle(changes)
	}

	h.mu.Lock()
	now := h.now()
	if h.deadline.IsZero() {
		h.deadline = now.Add(h.cfg.Window)
	}
	h.buf = append(h.buf, changes...)
	shoulFlush := now.After(h.deadline) || len(h.buf) >= h.cfg.MaxBatch
	var batch []monitor.Change
	if shouldFlush := shouldFlush; shouldFlush {
		batch = h.buf
		h.buf = nil
		h.deadline = time.Time{}
	}
	h.mu.Unlock()

	if batch != nil {
		return h.next.Handle(batch)
	}
	return nil
}

// Drain flushes any buffered changes immediately.
func (h *RollupHandler) Drain() error {
	h.mu.Lock()
	batch := h.buf
	h.buf = nil
	h.deadline = time.Time{}
	h.mu.Unlock()

	if len(batch) == 0 {
		return nil
	}
	return h.next.Handle(batch)
}
