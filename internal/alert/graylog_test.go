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

func graylogChange(kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "127.0.0.1",
			Hostname: "localhost",
			Port:     9000,
			Proto:    "tcp",
			Process:  "graylog",
		},
	}
}

func defaultGraylogCfg(url string) config.GraylogConfig {
	return config.GraylogConfig{
		Enabled:  true,
		URL:      url,
		Source:   "portwatch",
		Facility: "portwatch",
		Timeout:  5 * time.Second,
	}
}

func TestGraylogHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewGraylogHandler(defaultGraylogCfg(ts.URL))
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP request for empty changes")
	}
}

func TestGraylogHandler_PostsOnChange(t *testing.T) {
	var received graylogPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	h := NewGraylogHandler(defaultGraylogCfg(ts.URL))
	if err := h.Handle([]monitor.Change{graylogChange(monitor.ChangeAdded)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Version != "1.1" {
		t.Errorf("expected version 1.1, got %q", received.Version)
	}
	if received.Port != 9000 {
		t.Errorf("expected port 9000, got %d", received.Port)
	}
	if received.Action != "added" {
		t.Errorf("expected action added, got %q", received.Action)
	}
}

func TestGraylogHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewGraylogHandler(defaultGraylogCfg(ts.URL))
	err := h.Handle([]monitor.Change{graylogChange(monitor.ChangeAdded)})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestFormatGraylogMsg_FallsBackToIP(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.ChangeRemoved,
		Binding: ports.Binding{
			Addr:  "10.0.0.1",
			Port:  8080,
			Proto: "udp",
		},
	}
	cfg := config.DefaultGraylogConfig()
	msg := formatGraylogMsg(c, cfg)
	if msg.Host != cfg.Source {
		t.Errorf("expected host %q, got %q", cfg.Source, msg.Host)
	}
	if msg.Proto != "udp" {
		t.Errorf("expected proto udp, got %q", msg.Proto)
	}
}
