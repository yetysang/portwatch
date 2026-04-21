package alert

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

// InfluxDBHandler writes port change events to an InfluxDB v2 instance
// using the line protocol over HTTP.
type InfluxDBHandler struct {
	cfg    config.InfluxDBConfig
	client *http.Client
	now    func() time.Time
}

// NewInfluxDBHandler returns a new InfluxDBHandler using the given config.
func NewInfluxDBHandler(cfg config.InfluxDBConfig) *InfluxDBHandler {
	return &InfluxDBHandler{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
		now:    time.Now,
	}
}

// Handle sends each change as a line-protocol data point to InfluxDB.
func (h *InfluxDBHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	lines := make([]string, 0, len(changes))
	ts := h.now().UnixNano()
	for _, c := range changes {
		lines = append(lines, h.formatLine(c, ts))
	}
	body := strings.Join(lines, "\n")
	url := fmt.Sprintf("%s/api/v2/write?org=%s&bucket=%s&precision=ns",
		h.cfg.URL, h.cfg.Org, h.cfg.Bucket)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("influxdb: build request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+h.cfg.Token)
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("influxdb: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("influxdb: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// Drain is a no-op for InfluxDBHandler.
func (h *InfluxDBHandler) Drain() error { return nil }

func (h *InfluxDBHandler) formatLine(c monitor.Change, tsNano int64) string {
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	return fmt.Sprintf(
		"%s,action=%s,proto=%s,host=%s port=%d %d",
		h.cfg.Measurement,
		c.Kind,
		strings.ToLower(c.Binding.Proto),
		host,
		c.Binding.Port,
		tsNano,
	)
}
