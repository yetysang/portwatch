package alert

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

// DeadManHandler sends an HTTP GET heartbeat to a configured URL every time
// Handle is called with any changes. It acts as a dead-man's switch: if the
// upstream service (e.g. healthchecks.io) does not receive a ping within the
// expected interval, it fires an alert on its own.
type DeadManHandler struct {
	cfg    config.DeadManConfig
	client *http.Client
	logger *slog.Logger
}

// NewDeadManHandler constructs a DeadManHandler. If cfg.Enabled is false the
// handler is a no-op.
func NewDeadManHandler(cfg config.DeadManConfig, logger *slog.Logger) *DeadManHandler {
	return &DeadManHandler{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
		logger: logger,
	}
}

// Handle pings the dead-man URL whenever there are changes to report.
// If the list of changes is empty no request is sent — the heartbeat is
// intentionally tied to scan activity so a silent system still triggers the
// upstream switch.
func (h *DeadManHandler) Handle(changes []monitor.Change) error {
	if !h.cfg.Enabled || len(changes) == 0 {
		return nil
	}
	return h.ping()
}

// Heartbeat sends a ping unconditionally. Call this from the main tick loop
// after every successful scan so the switch resets even on quiet intervals.
func (h *DeadManHandler) Heartbeat() error {
	if !h.cfg.Enabled {
		return nil
	}
	return h.ping()
}

// Drain is a no-op for this handler.
func (h *DeadManHandler) Drain() {}

func (h *DeadManHandler) ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.cfg.URL, nil)
	if err != nil {
		return err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		h.logger.Warn("deadman: heartbeat ping failed", "url", h.cfg.URL, "err", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		h.logger.Warn("deadman: unexpected response", "url", h.cfg.URL, "status", resp.StatusCode)
		return fmt.Errorf("deadman: ping returned HTTP %d", resp.StatusCode)
	}

	h.logger.Debug("deadman: heartbeat sent", "url", h.cfg.URL, "at", time.Now().UTC().Format(time.RFC3339))
	return nil
}
