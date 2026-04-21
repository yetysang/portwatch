package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/wvictim/portwatch/internal/config"
	"github.com/wvictim/portwatch/internal/monitor"
	"github.com/wvictim/portwatch/internal/ports"
)

func newrelicChange(kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "127.0.0.1",
			Hostname: "localhost",
			Port:     9200,
			Proto:    "tcp",
			Process:  "elasticsearch",
			PID:      4321,
		},
	}
}

func defaultNewRelicCfg(url, accountID string) config.NewRelicConfig {
	return config.NewRelicConfig{
		Enabled:   true,
		APIKey:    "test-api-key",
		AccountID: accountID,
		Region:    "US",
		Timeout:   2 * time.Second,
	}
}

func TestNewRelicHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	cfg := defaultNewRelicCfg(ts.URL, "123456")
	h := NewNewRelicHandler(cfg)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no request for empty changes")
	}
}

func TestNewRelicHandler_PostsOnChange(t *testing.T) {
	var received []map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := defaultNewRelicCfg(ts.URL, "123456")
	// Override endpoint via a custom handler wrapping
	h := &newRelicHandler{
		cfg:      cfg,
		client:   &http.Client{Timeout: 2 * time.Second},
		endpoint: ts.URL,
	}
	if err := h.Handle([]monitor.Change{newrelicChange(monitor.Added)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) != 1 {
		t.Fatalf("expected 1 event, got %d", len(received))
	}
	if received[0]["action"] != "added" {
		t.Errorf("expected action=added, got %v", received[0]["action"])
	}
}

func TestNewRelicHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	h := &newRelicHandler{
		cfg:      defaultNewRelicCfg(ts.URL, "999"),
		client:   &http.Client{Timeout: 2 * time.Second},
		endpoint: ts.URL,
	}
	err := h.Handle([]monitor.Change{newrelicChange(monitor.Added)})
	if err == nil {
		t.Fatal("expected error on non-2xx status")
	}
}

func TestFormatNewRelicEvent_RemovedAction(t *testing.T) {
	c := newrelicChange(monitor.Removed)
	ev := formatNewRelicEvent(c)
	if ev.Action != "removed" {
		t.Errorf("expected removed, got %s", ev.Action)
	}
	if ev.EventType != "PortWatchEvent" {
		t.Errorf("unexpected event type: %s", ev.EventType)
	}
}

func TestFormatNewRelicEvent_FallsBackToIP(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{
			Addr:  "10.0.0.1",
			Port:  8080,
			Proto: "tcp",
		},
	}
	ev := formatNewRelicEvent(c)
	if ev.Addr != "10.0.0.1" {
		t.Errorf("expected IP fallback, got %s", ev.Addr)
	}
}
