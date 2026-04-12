package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/monitor"
)

// DiscordHandler sends change alerts to a Discord webhook.
type DiscordHandler struct {
	webhookURL string
	client     *http.Client
}

type discordPayload struct {
	Content string         `json:"content,omitempty"`
	Embeds  []discordEmbed `json:"embeds,omitempty"`
}

type discordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
}

// NewDiscordHandler creates a DiscordHandler that posts to the given webhook URL.
func NewDiscordHandler(webhookURL string, client *http.Client) *DiscordHandler {
	if client == nil {
		client = http.DefaultClient
	}
	return &DiscordHandler{webhookURL: webhookURL, client: client}
}

// Handle posts each change as a Discord embed message.
func (h *DiscordHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		payload := discordPayload{
			Embeds: []discordEmbed{formatDiscordEmbed(c)},
		}
		data, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("discord: marshal: %w", err)
		}
		resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("discord: post: %w", err)
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("discord: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}

// Drain is a no-op for DiscordHandler.
func (h *DiscordHandler) Drain() error { return nil }

func formatDiscordEmbed(c monitor.Change) discordEmbed {
	action := "added"
	color := 0x2ECC71 // green
	if c.Kind == monitor.Removed {
		action = "removed"
		color = 0xE74C3C // red
	}
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	desc := fmt.Sprintf("**%s** port **%d/%s** %s",
		host, c.Binding.Port, c.Binding.Proto, action)
	if c.Binding.Process != "" {
		desc += fmt.Sprintf(" (process: %s, pid: %d)", c.Binding.Process, c.Binding.PID)
	}
	return discordEmbed{
		Title:       "portwatch: port " + action,
		Description: desc,
		Color:       color,
	}
}
