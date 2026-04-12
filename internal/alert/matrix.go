package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// MatrixHandler sends alert messages to a Matrix room via the Client-Server API.
type MatrixHandler struct {
	homeserver string
	token      string
	roomID     string
	client     *http.Client
}

type matrixTextBody struct {
	MsgType       string `json:"msgtype"`
	Body          string `json:"body"`
	FormattedBody string `json:"formatted_body,omitempty"`
	Format        string `json:"format,omitempty"`
}

// NewMatrixHandler returns a handler that posts messages to a Matrix room.
func NewMatrixHandler(homeserver, token, roomID string) *MatrixHandler {
	return &MatrixHandler{
		homeserver: homeserver,
		token:      token,
		roomID:     roomID,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Handle sends each change as a Matrix message.
func (h *MatrixHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		if err := h.post(formatMatrixMsg(c)); err != nil {
			return err
		}
	}
	return nil
}

// Drain is a no-op for the Matrix handler.
func (h *MatrixHandler) Drain() error { return nil }

func (h *MatrixHandler) post(text string) error {
	body := matrixTextBody{
		MsgType: "m.text",
		Body:    text,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("matrix: marshal: %w", err)
	}
	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/send/m.room.message",
		h.homeserver, h.roomID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("matrix: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.token)
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("matrix: request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("matrix: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatMatrixMsg(c monitor.Change) string {
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	action := "bound"
	if c.Kind == monitor.Removed {
		action = "unbound"
	}
	proc := c.Binding.Process
	if proc == "" {
		proc = "unknown"
	}
	return fmt.Sprintf("[portwatch] Port %s/%s %s on %s (pid %d, %s)",
		c.Binding.Port, c.Binding.Proto, action, host, c.Binding.PID, proc)
}
