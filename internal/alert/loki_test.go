package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func lokiChange(kind monitor.ChangeKind, port int) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Proto:    "tcp",
			Port:     port,
			IP:       "127.0.0.1",
			Hostname: "localhost",
			Process:  "nginx",
			PID:      1234,
		},
	}
}

func TestLokiHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewLokiHandler(LokiConfig{Enabled: true, URL: ts.URL, JobLabel: "portwatch"})
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty changes")
	}
}

func TestLokiHandler_PostsOnChange(t *testing.T) {
	var captured []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	h := NewLokiHandler(LokiConfig{Enabled: true, URL: ts.URL, JobLabel: "portwatch"})
	if err := h.Handle([]monitor.Change{lokiChange(monitor.ChangeAdded, 8080)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(captured, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	streams, ok := payload["streams"].([]any)
	if !ok || len(streams) == 0 {
		t.Fatal("expected at least one stream")
	}
}

func TestLokiHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewLokiHandler(LokiConfig{Enabled: true, URL: ts.URL, JobLabel: "portwatch"})
	err := h.Handle([]monitor.Change{lokiChange(monitor.ChangeAdded, 9090)})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestFormatLokiMsg_FallsBackToIP(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.ChangeAdded,
		Binding: ports.Binding{Proto: "udp", Port: 53, IP: "10.0.0.1"},
	}
	msg := formatLokiMsg(c)
	if msg == "" {
		t.Fatal("expected non-empty message")
	}
	if !contains(msg, "10.0.0.1") {
		t.Errorf("expected IP in message, got: %s", msg)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
