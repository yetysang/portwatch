package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/warden-protocol/portwatch/internal/config"
	"github.com/warden-protocol/portwatch/internal/monitor"
	"github.com/warden-protocol/portwatch/internal/ports"
)

func hipchatChange() monitor.Change {
	return monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{
			Proto:    "tcp",
			Addr:     "0.0.0.0",
			Port:     9200,
			Hostname: "elastic-host",
			Process:  "elasticsearch",
			PID:      4321,
		},
	}
}

func TestHipChatHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	cfg := config.DefaultHipChatConfig()
	cfg.Enabled = true
	cfg.BaseURL = ts.URL
	cfg.AuthToken = "tok"
	cfg.RoomID = "42"
	h := NewHipChatHandler(cfg)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP request for empty changes")
	}
}

func TestHipChatHandler_PostsOnChange(t *testing.T) {
	var gotAuth, gotBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		var buf strings.Builder
		j := json.NewDecoder(r.Body)
		var p hipchatPayload
		_ = j.Decode(&p)
		gotBody = p.Message
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	cfg := config.DefaultHipChatConfig()
	cfg.Enabled = true
	cfg.BaseURL = ts.URL
	cfg.AuthToken = "mytoken"
	cfg.RoomID = "99"
	cfg.Timeout = 5 * time.Second
	h := NewHipChatHandler(cfg)
	if err := h.Handle([]monitor.Change{hipchatChange()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth != "Bearer mytoken" {
		t.Errorf("expected Bearer mytoken, got %q", gotAuth)
	}
	if !strings.Contains(gotBody, "9200") {
		t.Errorf("expected port 9200 in message, got %q", gotBody)
	}
	_ = gotBody
}

func TestHipChatHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	cfg := config.DefaultHipChatConfig()
	cfg.Enabled = true
	cfg.BaseURL = ts.URL
	cfg.AuthToken = "tok"
	cfg.RoomID = "1"
	h := NewHipChatHandler(cfg)
	if err := h.Handle([]monitor.Change{hipchatChange()}); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestFormatHipChatMsg_FallsBackToIP(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.Removed,
		Binding: ports.Binding{Proto: "udp", Addr: "127.0.0.1", Port: 53},
	}
	msg := formatHipChatMsg(c)
	if !strings.Contains(msg, "127.0.0.1") {
		t.Errorf("expected IP fallback in message, got %q", msg)
	}
	if !strings.Contains(msg, "unbound") {
		t.Errorf("expected 'unbound' in message, got %q", msg)
	}
}
