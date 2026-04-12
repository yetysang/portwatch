package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/monitor"
)

// TelegramHandler sends port change alerts to a Telegram chat via the Bot API.
type TelegramHandler struct {
	botToken string
	chatID   string
	client   *http.Client
}

type telegramPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// NewTelegramHandler creates a TelegramHandler that posts messages to the given chat.
func NewTelegramHandler(botToken, chatID string, client *http.Client) *TelegramHandler {
	if client == nil {
		client = http.DefaultClient
	}
	return &TelegramHandler{
		botToken: botToken,
		chatID:   chatID,
		client:   client,
	}
}

// Handle sends each change as a Telegram message.
func (h *TelegramHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		msg := formatTelegramMsg(c)
		payload := telegramPayload{
			ChatID:    h.chatID,
			Text:      msg,
			ParseMode: "Markdown",
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("telegram: marshal payload: %w", err)
		}
		url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", h.botToken)
		resp, err := h.client.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("telegram: post: %w", err)
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("telegram: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}

// Drain is a no-op for TelegramHandler.
func (h *TelegramHandler) Drain() error { return nil }

func formatTelegramMsg(c monitor.Change) string {
	action := "added"
	if c.Kind == monitor.Removed {
		action = "removed"
	}
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	proto := c.Binding.Proto
	if proto == "" {
		proto = "tcp"
	}
	msg := fmt.Sprintf("*portwatch* — port `%s/%d` *%s* on `%s`",
		proto, c.Binding.Port, action, host)
	if c.Binding.Process != "" {
		msg += fmt.Sprintf(" (process: `%s`", c.Binding.Process)
		if c.Binding.PID > 0 {
			msg += fmt.Sprintf(", pid: `%d`", c.Binding.PID)
		}
		msg += ")"
	}
	return msg
}
