package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func slackChange(kind monitor.ChangeKind, port uint16) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "127.0.0.1",
			Hostname: "localhost",
			Port:     port,
			Proto:    "tcp",
			Process:  "nginx",
			PID:      1234,
		},
	}
}

func TestSlackHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewSlackHandler(ts.URL, "")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty changes")
	}
}

func TestSlackHandler_PostsOnChange(t *testing.T) {
	var received slackPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewSlackHandler(ts.URL, "[portwatch]")
	c := slackChange(monitor.Added, 8080)
	if err := h.Handle([]monitor.Change{c}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(received.Text, "[portwatch]") {
		t.Errorf("expected prefix in message, got: %q", received.Text)
	}
	if !strings.Contains(received.Text, "8080") {
		t.Errorf("expected port in message, got: %q", received.Text)
	}
}

func TestSlackHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewSlackHandler(ts.URL, "")
	err := h.Handle([]monitor.Change{slackChange(monitor.Added, 9090)})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestFormatSlackMsg_RemovedAction(t *testing.T) {
	c := slackChange(monitor.Removed, 443)
	msg := formatSlackMsg(c, "")
	if !strings.Contains(msg, "removed") {
		t.Errorf("expected 'removed' in message, got: %q", msg)
	}
}

func TestFormatSlackMsg_FallsBackToIP(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{
			Addr:  "10.0.0.1",
			Port:  22,
			Proto: "tcp",
		},
	}
	msg := formatSlackMsg(c, "")
	if !strings.Contains(msg, "10.0.0.1") {
		t.Errorf("expected IP fallback in message, got: %q", msg)
	}
}
