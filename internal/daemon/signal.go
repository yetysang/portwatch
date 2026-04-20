// Package daemon provides OS-signal lifecycle management for portwatch.
package daemon

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/patrickward/portwatch/internal/config"
)

// SignalHandler listens for OS signals and coordinates graceful shutdown or
// config reload.
type SignalHandler struct {
	cfg    config.SignalConfig
	logger *slog.Logger
	onReload func()
}

// NewSignalHandler creates a SignalHandler that calls onReload when SIGHUP is
// received (if ReloadOnHUP is enabled).
func NewSignalHandler(cfg config.SignalConfig, logger *slog.Logger, onReload func()) *SignalHandler {
	if onReload == nil {
		onReload = func() {}
	}
	return &SignalHandler{cfg: cfg, logger: logger, onReload: onReload}
}

// Run blocks until an interrupt signal is received, then cancels the supplied
// context and waits up to GracePeriodSeconds for the caller to finish.
// It returns once the grace period has elapsed or done is closed.
func (h *SignalHandler) Run(ctx context.Context, cancel context.CancelFunc, done <-chan struct{}) {
	sigs := make(chan os.Signal, 2)
	notify := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	if h.cfg.ReloadOnHUP {
		notify = append(notify, syscall.SIGHUP)
	}
	signal.Notify(sigs, notify...)
	defer signal.Stop(sigs)

	for {
		select {
		case sig := <-sigs:
			switch sig {
			case syscall.SIGHUP:
				h.logger.Info("received SIGHUP, reloading config")
				h.onReload()
			default:
				h.logger.Info("received signal, shutting down", "signal", sig)
				cancel()
				grace := time.Duration(h.cfg.GracePeriodSeconds) * time.Second
				select {
				case <-done:
					h.logger.Info("clean shutdown complete")
				case <-time.After(grace):
					h.logger.Warn("grace period expired, forcing exit")
				}
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
