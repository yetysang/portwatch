package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/wolveix/portwatch/internal/config"
	"github.com/wolveix/portwatch/internal/monitor"
	"github.com/wolveix/portwatch/internal/ports"
)

func mattermostChange(kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "127.0.0.1",
			Hostname: "localhost",
			Port:     8080,
			Proto:    "tcp",
			Process:  "nginx",
		},
	}
}

func TestMattermostHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	cfg := config.DefaultMattermostConfig()
	cfg.Enabled = true
	cfg.URL = ts.URL
	cfg.Channel = "#test"
	h := NewMattermostHandler(cfg)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty changes")
	}
}

func TestMattermostHandler_PostsOnChange(t *testing.T) {
	var gotBody []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := config.DefaultMattermostConfig()
	cfg.Enabled = true
	cfg.URL = ts.URL
	cfg.Channel = "#alerts"
	cfg.Timeout = 5 * time.Second
	h := NewMattermostHandler(cfg)

	if err := h.Handle([]monitor.Change{mattermostChange(monitor.Added)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(gotBody, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	if payload["channel"] != "#alerts" {
		t.Errorf("expected channel #alerts, got %v", payload["channel"])
	}
	if !strings.Contains(payload["text"].(string), "bound") {
		t.Errorf("expected 'bound' in text, got %q", payload["text"])
	}
}

func TestMattermostHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	cfg := config.DefaultMattermostConfig()
	cfg.Enabled = true
	cfg.URL = ts.URL
	cfg.Channel = "#alerts"
	h := NewMattermostHandler(cfg)

	if err := h.Handle([]monitor.Change{mattermostChange(monitor.Added)}); err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestFormatMattermostMsg_RemovedAction(t *testing.T) {
	c := mattermostChange(monitor.Removed)
	msg := formatMattermostMsg(c)
	if !strings.Contains(msg, "unbound") {
		t.Errorf("expected 'unbound' in message, got %q", msg)
	}
}

func TestFormatMattermostMsg_FallsBackToIP(t *testing.T) {
	c := mattermostChange(monitor.Added)
	c.Binding.Hostname = ""
	msg := formatMattermostMsg(c)
	if !strings.Contains(msg, "127.0.0.1") {
		t.Errorf("expected IP fallback in message, got %q", msg)
	}
}
