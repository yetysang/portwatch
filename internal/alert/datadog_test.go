package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func datadogChange(kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "0.0.0.0",
			Port:     "9090",
			Proto:    "tcp",
			Hostname: "myhost",
			Process:  "prometheus",
			PID:      42,
		},
	}
}

func TestDatadogHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	cfg := DatadogConfig{Enabled: true, APIKey: "test-key", Site: ts.Listener.Addr().String()}
	h := NewDatadogHandler(cfg)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty changes")
	}
}

func TestDatadogHandler_PostsOnChange(t *testing.T) {
	var received datadogEvent
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if r.Header.Get("DD-API-KEY") != "secret" {
			t.Errorf("missing or wrong DD-API-KEY header")
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	cfg := DatadogConfig{Enabled: true, APIKey: "secret", Site: ts.Listener.Addr().String()}
	h := NewDatadogHandler(cfg)
	if err := h.Handle([]monitor.Change{datadogChange(monitor.Added)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Title == "" {
		t.Error("expected non-empty event title")
	}
	if received.AlertType != "warning" {
		t.Errorf("expected alert_type=warning for Added, got %q", received.AlertType)
	}
}

func TestDatadogHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	cfg := DatadogConfig{Enabled: true, APIKey: "bad", Site: ts.Listener.Addr().String()}
	h := NewDatadogHandler(cfg)
	err := h.Handle([]monitor.Change{datadogChange(monitor.Added)})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestFormatDatadogEvent_RemovedAlertType(t *testing.T) {
	evt := formatDatadogEvent(datadogChange(monitor.Removed))
	if evt.AlertType != "info" {
		t.Errorf("expected alert_type=info for Removed, got %q", evt.AlertType)
	}
}

func TestFormatDatadogEvent_TagsIncludePort(t *testing.T) {
	evt := formatDatadogEvent(datadogChange(monitor.Added))
	found := false
	for _, tag := range evt.Tags {
		if tag == "port:9090" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected tag port:9090 in %v", evt.Tags)
	}
}
