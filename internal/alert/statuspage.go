package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

const statuspageBaseURL = "https://api.statuspage.io/v1"

// StatusPageHandler posts component status updates to Statuspage.io.
type StatusPageHandler struct {
	cfg    config.StatusPageConfig
	client *http.Client
}

// NewStatusPageHandler creates a new StatusPageHandler.
func NewStatusPageHandler(cfg config.StatusPageConfig) *StatusPageHandler {
	return &StatusPageHandler{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

type statuspageComponentBody struct {
	Component statuspageStatus `json:"component"`
}

type statuspageStatus struct {
	Status string `json:"status"`
}

// Handle posts a component status update for each change batch.
func (h *StatusPageHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	status := "operational"
	for _, c := range changes {
		if c.Kind == monitor.Added {
			status = "under_maintenance"
			break
		}
	}
	body := statuspageComponentBody{
		Component: statuspageStatus{Status: status},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("statuspage: marshal payload: %w", err)
	}
	url := fmt.Sprintf("%s/pages/%s/components/%s",
		statuspageBaseURL, h.cfg.PageID, h.cfg.ComponentID)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("statuspage: build request: %w", err)
	}
	req.Header.Set("Authorization", "OAuth "+h.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("statuspage: do request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("statuspage: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// Drain is a no-op for this handler.
func (h *StatusPageHandler) Drain() error { return nil }
