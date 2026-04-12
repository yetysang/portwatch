package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func teamsChange(kind monitor.ChangeKind, port, proto string) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			IP:       "127.0.0.1",
			Hostname: "localhost",
			Port:     port,
			Proto:    proto,
			PID:      1234,
			Process:  "nginx",
		},
	}
}

func TestTeamsHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewTeamsHandler(ts.URL)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty changes")
	}
}

func TestTeamsHandler_PostsOnChange(t *testing.T) {
	var captured []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, r.ContentLength)
		r.Body.Read(buf)
		captured = buf
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewTeamsHandler(ts.URL)
	if err := h.Handle([]monitor.Change{teamsChange(monitor.Added, "8080", "tcp")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload teamsPayload
	if err := json.Unmarshal(captured, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	if payload.Type != "message" {
		t.Errorf("expected type 'message', got %q", payload.Type)
	}
	if len(payload.Attachments) == 0 {
		t.Fatal("expected at least one attachment")
	}
	body := payload.Attachments[0].Content.Body
	if len(body) == 0 || body[0].Text == "" {
		t.Error("expected non-empty text block")
	}
}

func TestTeamsHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewTeamsHandler(ts.URL)
	err := h.Handle([]monitor.Change{teamsChange(monitor.Added, "443", "tcp")})
	if err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestFormatTeamsMsg_RemovedAction(t *testing.T) {
	c := teamsChange(monitor.Removed, "22", "tcp")
	msg := formatTeamsMsg(c)
	if msg == "" {
		t.Fatal("expected non-empty message")
	}
	for _, want := range []string{"removed", "22", "tcp", "localhost", "nginx"} {
		if !contains(msg, want) {
			t.Errorf("expected %q in message %q", want, msg)
		}
	}
}

func TestFormatTeamsMsg_FallsBackToIP(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{IP: "10.0.0.1", Port: "9090", Proto: "udp"},
	}
	msg := formatTeamsMsg(c)
	if !contains(msg, "10.0.0.1") {
		t.Errorf("expected IP fallback in message %q", msg)
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
