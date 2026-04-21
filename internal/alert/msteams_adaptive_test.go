package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func adaptiveChange(kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "127.0.0.1",
			Hostname: "localhost",
			Port:     8080,
			Proto:    "tcp",
		},
	}
}

func defaultAdaptiveCfg(url string) config.AdaptiveCardConfig {
	return config.AdaptiveCardConfig{
		Enabled:    true,
		WebhookURL: url,
		Timeout:    5 * time.Second,
		ThemeColor: "0078D4",
	}
}

func TestAdaptiveCardHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewAdaptiveCardHandler(defaultAdaptiveCfg(ts.URL))
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty changes")
	}
}

func TestAdaptiveCardHandler_PostsOnChange(t *testing.T) {
	var received map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewAdaptiveCardHandler(defaultAdaptiveCfg(ts.URL))
	if err := h.Handle([]monitor.Change{adaptiveChange(monitor.Added)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received == nil {
		t.Fatal("expected payload to be sent")
	}
	if received["themeColor"] != "0078D4" {
		t.Errorf("expected themeColor 0078D4, got %v", received["themeColor"])
	}
}

func TestAdaptiveCardHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewAdaptiveCardHandler(defaultAdaptiveCfg(ts.URL))
	err := h.Handle([]monitor.Change{adaptiveChange(monitor.Added)})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestBuildAdaptiveCardPayload_RemovedAction(t *testing.T) {
	changes := []monitor.Change{adaptiveChange(monitor.Removed)}
	payload := buildAdaptiveCardPayload(changes, "FF0000")
	sections, ok := payload["sections"].([]map[string]any)
	if !ok || len(sections) == 0 {
		t.Fatal("expected sections in payload")
	}
	facts, ok := sections[0]["facts"].([]map[string]string)
	if !ok || len(facts) == 0 {
		t.Fatal("expected facts in section")
	}
	if facts[0]["value"] != "removed" {
		t.Errorf("expected action 'removed', got %q", facts[0]["value"])
	}
}
