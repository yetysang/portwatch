package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

// OpsGenieConfig holds configuration for the OpsGenie alert handler.
type OpsGenieConfig struct {
	APIKey  string
	Team    string
	BaseURL string // defaults to https://api.opsgenie.com
}

type opsGenieHandler struct {
	cfg    OpsGenieConfig
	client *http.Client
}

// NewOpsGenieHandler returns a Handler that sends alerts to OpsGenie.
func NewOpsGenieHandler(cfg OpsGenieConfig) Handler {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.opsgenie.com"
	}
	return &opsGenieHandler{cfg: cfg, client: &http.Client{}}
}

func (h *opsGenieHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		if err := h.send(c); err != nil {
			return err
		}
	}
	return nil
}

func (h *opsGenieHandler) Drain() error { return nil }

func (h *opsGenieHandler) send(c monitor.Change) error {
	msg := formatOpsGenieMsg(c)
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("opsgenie: marshal: %w", err)
	}
	url := strings.TrimRight(h.cfg.BaseURL, "/") + "/v2/alerts"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("opsgenie: request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+h.cfg.APIKey)
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: http: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status %d", resp.StatusCode)
	}
	return nil
}

type opsGeniePayload struct {
	Message  string   `json:"message"`
	Alias    string   `json:"alias"`
	Priority string   `json:"priority"`
	Tags     []string `json:"tags"`
	Details  map[string]string `json:"details,omitempty"`
}

func formatOpsGenieMsg(c monitor.Change) opsGeniePayload {
	action := "added"
	priority := "P3"
	if c.Kind == monitor.Removed {
		action = "removed"
		priority = "P5"
	}
	proto := strings.ToLower(c.Binding.Proto)
	alias := fmt.Sprintf("portwatch-%s-%d-%s", proto, c.Binding.Port, action)
	msg := fmt.Sprintf("Port %s/%d %s", strings.ToUpper(proto), c.Binding.Port, action)
	if c.Binding.Process != "" {
		msg += " by " + c.Binding.Process
	}
	details := map[string]string{
		"addr":  c.Binding.Addr,
		"proto": strings.ToUpper(proto),
		"port":  fmt.Sprintf("%d", c.Binding.Port),
	}
	if c.Binding.Process != "" {
		details["process"] = c.Binding.Process
	}
	return opsGeniePayload{
		Message:  msg,
		Alias:    alias,
		Priority: priority,
		Tags:     []string{"portwatch", proto, action},
		Details:  details,
	}
}
