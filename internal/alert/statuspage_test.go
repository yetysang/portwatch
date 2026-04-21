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

func statuspageChange(kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{Port: 8080, Proto: "tcp"},
	}
}

func defaultStatusPageCfg(url string) config.StatusPageConfig {
	return config.StatusPageConfig{
		Enabled:     true,
		APIKey:      "test-key",
		PageID:      "page123",
		ComponentID: "comp456",
		Timeout:     5 * time.Second,
	}
}

func TestStatusPageHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()
	h := NewStatusPageHandler(defaultStatusPageCfg(ts.URL))
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP request for empty changes")
	}
}

func TestStatusPageHandler_PostsOnChange(t *testing.T) {
	var gotBody map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	cfg := defaultStatusPageCfg(ts.URL)
	h := NewStatusPageHandler(cfg)
	if err := h.Handle([]monitor.Change{statuspageChange(monitor.Added)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStatusPageHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()
	cfg := defaultStatusPageCfg(ts.URL)
	h := NewStatusPageHandler(cfg)
	err := h.Handle([]monitor.Change{statuspageChange(monitor.Added)})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestStatusPageHandler_RemovedSetsOperational(t *testing.T) {
	var gotStatus string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		gotStatus = body["component"]["status"]
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	cfg := defaultStatusPageCfg(ts.URL)
	h := NewStatusPageHandler(cfg)
	_ = h.Handle([]monitor.Change{statuspageChange(monitor.Removed)})
	if gotStatus != "operational" {
		t.Errorf("expected operational, got %q", gotStatus)
	}
}
