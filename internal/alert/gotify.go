package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

// GotifyConfig holds configuration for the Gotify push notification handler.
type GotifyConfig struct {
	Enabled  bool
	URL      string
	Token    string
	Priority int
}

type gotifyPayload struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

type gotifyHandler struct {
	cfg    GotifyConfig
	client *http.Client
}

// NewGotifyHandler returns a Handler that sends alerts to a Gotify server.
func NewGotifyHandler(cfg GotifyConfig, client *http.Client) Handler {
	if client == nil {
		client = http.DefaultClient
	}
	return &gotifyHandler{cfg: cfg, client: client}
}

func (g *gotifyHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	msg := formatGotifyMsg(changes)
	payload := gotifyPayload{
		Title:    "portwatch alert",
		Message:  msg,
		Priority: g.cfg.Priority,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("gotify: marshal payload: %w", err)
	}
	url := strings.TrimRight(g.cfg.URL, "/") + "/message?token=" + g.cfg.Token
	resp, err := g.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("gotify: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (g *gotifyHandler) Drain() error { return nil }

func formatGotifyMsg(changes []monitor.Change) string {
	var sb strings.Builder
	for _, c := range changes {
		host := c.Binding.Hostname
		if host == "" {
			host = c.Binding.Addr
		}
		proto := strings.ToUpper(c.Binding.Proto)
		fmt.Fprintf(&sb, "[%s] %s %s:%d\n", c.Kind, proto, host, c.Binding.Port)
	}
	return strings.TrimRight(sb.String(), "\n")
}
