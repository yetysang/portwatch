// Package daemon provides runtime lifecycle helpers for portwatch.
package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/user/portwatch/internal/config"
)

// HealthServer exposes a lightweight HTTP endpoint for liveness probes.
type HealthServer struct {
	cfg    config.HealthCheckConfig
	ready  atomic.Bool
	server *http.Server
}

// NewHealthServer creates a HealthServer from cfg.
func NewHealthServer(cfg config.HealthCheckConfig) *HealthServer {
	h := &HealthServer{cfg: cfg}
	mux := http.NewServeMux()
	mux.HandleFunc(cfg.Path, h.handle)
	h.server = &http.Server{
		Addr:        cfg.ListenAddr,
		Handler:     mux,
		ReadTimeout: cfg.ReadTimeout,
	}
	return h
}

// SetReady marks the daemon as ready to serve traffic.
func (h *HealthServer) SetReady(v bool) { h.ready.Store(v) }

// Start begins listening in a background goroutine.
// It returns an error if the listener cannot be bound.
func (h *HealthServer) Start() error {
	if !h.cfg.Enabled {
		return nil
	}
	errCh := make(chan error, 1)
	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	select {
	case err := <-errCh:
		return fmt.Errorf("healthcheck: %w", err)
	case <-time.After(50 * time.Millisecond):
		return nil
	}
}

// Shutdown gracefully stops the HTTP server.
func (h *HealthServer) Shutdown(ctx context.Context) error {
	if !h.cfg.Enabled {
		return nil
	}
	return h.server.Shutdown(ctx)
}

func (h *HealthServer) handle(w http.ResponseWriter, _ *http.Request) {
	type response struct {
		Status string `json:"status"`
	}
	w.Header().Set("Content-Type", "application/json")
	if h.ready.Load() {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response{Status: "ok"})
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(response{Status: "starting"})
	}
}
