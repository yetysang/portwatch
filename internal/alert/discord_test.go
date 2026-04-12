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

func discordChange(kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "127.0.0.1",
			Hostname: "localhost",
			Port:     8080,
			Proto:    "tcp",
			Process:  "myapp",
			PID:      1234,
		},
	}
}

func TestDiscordHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewDiscordHandler(ts.URL, ts.Client())
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP request for empty changes")
	}
}

func TestDiscordHandler_PostsOnChange(t *testing.T) {
	var received discordPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	h := NewDiscordHandler(ts.URL, ts.Client())
	if err := h.Handle([]monitor.Change{discordChange(monitor.Added)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received.Embeds) != 1 {
		t.Fatalf("expected 1 embed, got %d", len(received.Embeds))
	}
	if received.Embeds[0].Title != "portwatch: port added" {
		t.Errorf("unexpected title: %q", received.Embeds[0].Title)
	}
}

func TestDiscordHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	h := NewDiscordHandler(ts.URL, ts.Client())
	err := h.Handle([]monitor.Change{discordChange(monitor.Added)})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestFormatDiscordEmbed_RemovedColor(t *testing.T) {
	c := discordChange(monitor.Removed)
	embed := formatDiscordEmbed(c)
	if embed.Color != 0xE74C3C {
		t.Errorf("expected red color for removed, got %#x", embed.Color)
	}
	if embed.Title != "portwatch: port removed" {
		t.Errorf("unexpected title: %q", embed.Title)
	}
}

func TestFormatDiscordEmbed_FallsBackToIP(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{Addr: "10.0.0.1", Port: 443, Proto: "tcp"},
	}
	embed := formatDiscordEmbed(c)
	if embed.Description == "" {
		t.Error("expected non-empty description")
	}
}
