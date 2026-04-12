package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func ntfyChange(kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "127.0.0.1",
			Port:     "9090",
			Proto:    "tcp",
			Hostname: "localhost",
			Process:  "myapp",
			PID:      42,
		},
	}
}

func TestNtfyHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewNtfyHandler(NtfyConfig{ServerURL: ts.URL, Topic: "alerts"})
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty changes")
	}
}

func TestNtfyHandler_PostsOnChange(t *testing.T) {
	var received ntfyPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewNtfyHandler(NtfyConfig{ServerURL: ts.URL, Topic: "portwatch"})
	if err := h.Handle([]monitor.Change{ntfyChange(monitor.Added)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Topic != "portwatch" {
		t.Errorf("topic = %q, want %q", received.Topic, "portwatch")
	}
	if received.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestNtfyHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	h := NewNtfyHandler(NtfyConfig{ServerURL: ts.URL, Topic: "alerts"})
	err := h.Handle([]monitor.Change{ntfyChange(monitor.Added)})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestFormatNtfyMsg_Added(t *testing.T) {
	c := ntfyChange(monitor.Added)
	p := formatNtfyMsg("test", c)
	if p.Priority != "default" {
		t.Errorf("priority = %q, want %q", p.Priority, "default")
	}
	if len(p.Tags) == 0 {
		t.Error("expected at least one tag")
	}
}

func TestFormatNtfyMsg_Removed(t *testing.T) {
	c := ntfyChange(monitor.Removed)
	p := formatNtfyMsg("test", c)
	if p.Priority != "high" {
		t.Errorf("priority = %q, want %q", p.Priority, "high")
	}
}
