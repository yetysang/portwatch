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

// cloudwatchEvent is a single structured log event sent to CloudWatch Logs.
type cloudwatchEvent struct {
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

// cloudwatchPayload wraps one or more log events for the PutLogEvents API.
type cloudwatchPayload struct {
	LogGroupName  string            `json:"logGroupName"`
	LogStreamName string            `json:"logStreamName"`
	LogEvents     []cloudwatchEvent `json:"logEvents"`
}

// CloudWatchHandler ships port-change alerts to AWS CloudWatch Logs via the
// PutLogEvents REST endpoint (no SDK dependency).
type CloudWatchHandler struct {
	cfg    config.CloudWatchConfig
	client *http.Client
}

// NewCloudWatchHandler constructs a CloudWatchHandler. When cfg.Enabled is
// false the handler is a no-op.
func NewCloudWatchHandler(cfg config.CloudWatchConfig) *CloudWatchHandler {
	return &CloudWatchHandler{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Handle ships each change as a structured CloudWatch log event.
func (h *CloudWatchHandler) Handle(changes []monitor.Change) error {
	if !h.cfg.Enabled || len(changes) == 0 {
		return nil
	}

	events := make([]cloudwatchEvent, 0, len(changes))
	for _, c := range changes {
		msg := formatCloudWatchMsg(c)
		events = append(events, cloudwatchEvent{
			Timestamp: time.Now().UnixMilli(),
			Message:   msg,
		})
	}

	payload := cloudwatchPayload{
		LogGroupName:  h.cfg.LogGroup,
		LogStreamName: h.cfg.LogStream,
		LogEvents:     events,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cloudwatch: marshal payload: %w", err)
	}

	endpoint := fmt.Sprintf("https://logs.%s.amazonaws.com/", h.cfg.Region)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cloudwatch: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-amz-json-1.1")
	req.Header.Set("X-Amz-Target", "Logs_20140328.PutLogEvents")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("cloudwatch: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("cloudwatch: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// Drain is a no-op for CloudWatchHandler.
func (h *CloudWatchHandler) Drain() {}

func formatCloudWatchMsg(c monitor.Change) string {
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	return fmt.Sprintf("action=%s proto=%s host=%s port=%d pid=%d process=%s",
		c.Kind, c.Binding.Proto, host, c.Binding.Port, c.Binding.PID, c.Binding.Process)
}
