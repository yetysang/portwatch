package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func influxChange(kind monitor.ChangeKind, port uint16) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr: "127.0.0.1",
			Port: port,
			Proto: "tcp",
			Hostname: "localhost",
			Process: "nginx",
			PID: 1234,
		},
	}
}

func TestInfluxDBHandler_EmptyChangesNoRequest(t *testing.T) {
	requested := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	cfg := config.DefaultInfluxDBConfig()
	cfg.Enabled = true
	cfg.URL = ts.URL
	cfg.Token = "test-token"
	cfg.Org = "myorg"
	cfg.Bucket = "portwatch"

	h := alert.NewInfluxDBHandler(cfg)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if requested {
		t.Fatal("expected no HTTP request for empty changes")
	}
}

func TestInfluxDBHandler_PostsOnChange(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Token test-token" {
			t.Errorf("missing or wrong Authorization header: %s", r.Header.Get("Authorization"))
		}
		var err error
		received, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	cfg := config.DefaultInfluxDBConfig()
	cfg.Enabled = true
	cfg.URL = ts.URL
	cfg.Token = "test-token"
	cfg.Org = "myorg"
	cfg.Bucket = "portwatch"

	h := alert.NewInfluxDBHandler(cfg)
	changes := []monitor.Change{
		influxChange(monitor.ChangeAdded, 8080),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) == 0 {
		t.Fatal("expected request body, got none")
	}
	body := string(received)
	if !contains(body, "portwatch_binding") {
		t.Errorf("expected measurement name in body, got: %s", body)
	}
	if !contains(body, "port=8080") {
		t.Errorf("expected port tag in body, got: %s", body)
	}
	if !contains(body, "action=added") {
		t.Errorf("expected action tag in body, got: %s", body)
	}
}

func TestInfluxDBHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	cfg := config.DefaultInfluxDBConfig()
	cfg.Enabled = true
	cfg.URL = ts.URL
	cfg.Token = "test-token"
	cfg.Org = "myorg"
	cfg.Bucket = "portwatch"

	h := alert.NewInfluxDBHandler(cfg)
	changes := []monitor.Change{
		influxChange(monitor.ChangeAdded, 9090),
	}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestInfluxDBHandler_DisabledIsNoop(t *testing.T) {
	requested := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	cfg := config.DefaultInfluxDBConfig()
	cfg.Enabled = false
	cfg.URL = ts.URL
	cfg.Token = "test-token"
	cfg.Org = "myorg"
	cfg.Bucket = "portwatch"

	h := alert.NewInfluxDBHandler(cfg)
	changes := []monitor.Change{
		influxChange(monitor.ChangeAdded, 8080),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if requested {
		t.Fatal("expected no HTTP request when handler is disabled")
	}
}

func TestInfluxDBHandler_LineProtocolFormat(t *testing.T) {
	var receivedBody string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		receivedBody = string(b)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	cfg := config.DefaultInfluxDBConfig()
	cfg.Enabled = true
	cfg.URL = ts.URL
	cfg.Token = "tok"
	cfg.Org = "org"
	cfg.Bucket = "bkt"
	cfg.Measurement = "port_events"

	h := alert.NewInfluxDBHandler(cfg)
	changes := []monitor.Change{
		influxChange(monitor.ChangeRemoved, 443),
	}
	_ = h.Handle(changes)

	if !contains(receivedBody, "port_events") {
		t.Errorf("expected custom measurement name, got: %s", receivedBody)
	}
	if !contains(receivedBody, "action=removed") {
		t.Errorf("expected removed action, got: %s", receivedBody)
	}
	if !contains(receivedBody, "proto=tcp") {
		t.Errorf("expected proto tag, got: %s", receivedBody)
	}
}

// contains is a helper used across alert tests.
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

// Ensure unused imports are referenced.
var _ = json.Marshal
var _ = time.Now
