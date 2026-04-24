package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/wneessen/portwatch/internal/config"
	"github.com/wneessen/portwatch/internal/monitor"
	"github.com/wneessen/portwatch/internal/ports"
)

func googlechatChange(kind monitor.ChangeKind, port uint16) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "127.0.0.1",
			Hostname: "localhost",
			Port:     port,
			Proto:    "tcp",
			Process:  "sshd",
			PID:      1234,
		},
		At: time.Now(),
	}
}

func TestGoogleChatHandler_EmptyChangesNoRequest(t *testing.T) {
	requested := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested = true
	}))
	defer ts.Close()

	h := NewGoogleChatHandler(config.GoogleChatConfig{Enabled: true, WebhookURL: ts.URL, Timeout: 5 * time.Second})
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if requested {
		t.Error("expected no HTTP request for empty changes")
	}
}

func TestGoogleChatHandler_PostsOnChange(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewGoogleChatHandler(config.GoogleChatConfig{Enabled: true, WebhookURL: ts.URL, Timeout: 5 * time.Second})
	changes := []monitor.Change{googlechatChange(monitor.Added, 22)}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["text"] == "" {
		t.Error("expected non-empty text in payload")
	}
}

func TestGoogleChatHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewGoogleChatHandler(config.GoogleChatConfig{Enabled: true, WebhookURL: ts.URL, Timeout: 5 * time.Second})
	changes := []monitor.Change{googlechatChange(monitor.Added, 8080)}
	if err := h.Handle(changes); err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestFormatGoogleChatMsg_AddedAction(t *testing.T) {
	c := googlechatChange(monitor.Added, 443)
	msg := formatGoogleChatMsg(c)
	if msg == "" {
		t.Fatal("expected non-empty message")
	}
	for _, want := range []string{"443", "bound", "localhost"} {
		if !containsStr(msg, want) {
			t.Errorf("expected message to contain %q, got: %s", want, msg)
		}
	}
}

func TestFormatGoogleChatMsg_RemovedAction(t *testing.T) {
	c := googlechatChange(monitor.Removed, 80)
	msg := formatGoogleChatMsg(c)
	if !containsStr(msg, "unbound") {
		t.Errorf("expected 'unbound' in message, got: %s", msg)
	}
}
