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

func gotifyChange(kind, addr string, port int) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:  addr,
			Port:  port,
			Proto: "tcp",
		},
	}
}

func TestGotifyHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	cfg := GotifyConfig{Enabled: true, URL: ts.URL, Token: "tok", Priority: 5}
	h := NewGotifyHandler(cfg, ts.Client())
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP request for empty changes")
	}
}

func TestGotifyHandler_PostsOnChange(t *testing.T) {
	var captured struct {
		Title    string `json:"title"`
		Message  string `json:"message"`
		Priority int    `json:"priority"`
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := GotifyConfig{Enabled: true, URL: ts.URL, Token: "mytoken", Priority: 7}
	h := NewGotifyHandler(cfg, ts.Client())
	changes := []monitor.Change{gotifyChange("added", "127.0.0.1", 8080)}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.Title != "portwatch alert" {
		t.Errorf("expected title 'portwatch alert', got %q", captured.Title)
	}
	if captured.Priority != 7 {
		t.Errorf("expected priority 7, got %d", captured.Priority)
	}
	if !strings.Contains(captured.Message, "8080") {
		t.Errorf("expected message to contain port, got %q", captured.Message)
	}
}

func TestGotifyHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	cfg := GotifyConfig{Enabled: true, URL: ts.URL, Token: "bad", Priority: 5}
	h := NewGotifyHandler(cfg, ts.Client())
	changes := []monitor.Change{gotifyChange("added", "0.0.0.0", 9000)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestFormatGotifyMsg_UsesHostname(t *testing.T) {
	changes := []monitor.Change{{
		Kind: "added",
		Binding: ports.Binding{
			Addr:     "1.2.3.4",
			Hostname: "myhost",
			Port:     443,
			Proto:    "tcp",
		},
	}}
	msg := formatGotifyMsg(changes)
	if !strings.Contains(msg, "myhost") {
		t.Errorf("expected hostname in message, got %q", msg)
	}
	if strings.Contains(msg, "1.2.3.4") {
		t.Errorf("expected IP to be replaced by hostname, got %q", msg)
	}
}
