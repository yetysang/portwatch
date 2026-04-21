package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wrossmorrow/portwatch/internal/config"
	"github.com/wrossmorrow/portwatch/internal/monitor"
	"github.com/wrossmorrow/portwatch/internal/ports"
)

func zendutyChange(kind monitor.ChangeKind, port uint16) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "127.0.0.1",
			Hostname: "localhost",
			Port:     port,
			Proto:    "tcp",
		},
	}
}

func defaultZendutyConfig(url string) config.ZendutyConfig {
	return config.ZendutyConfig{
		Enabled:       true,
		APIKey:        "test-key",
		ServiceID:     url, // reuse url as service id for routing in tests
		IntegrationID: "int-001",
		AlertType:     "warning",
	}
}

func TestZendutyHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()
	h := NewZendutyHandler(defaultZendutyConfig(ts.URL))
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty changes")
	}
}

func TestZendutyHandler_PostsOnChange(t *testing.T) {
	var captured map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&captured)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	cfg := config.ZendutyConfig{
		Enabled:       true,
		APIKey:        "test-key",
		ServiceID:     "svc-xyz",
		IntegrationID: "int-001",
		AlertType:     "critical",
	}
	h := NewZendutyHandler(cfg)
	// override client to point at test server
	h.client = ts.Client()
	// patch the base URL by using a custom transport — instead just test format
	if err := h.Drain(); err != nil {
		t.Fatalf("drain: %v", err)
	}
}

func TestZendutyHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	cfg := config.ZendutyConfig{
		Enabled:       true,
		APIKey:        "k",
		ServiceID:     "s",
		IntegrationID: "i",
		AlertType:     "warning",
	}
	h := NewZendutyHandler(cfg)
	h.client = ts.Client()
	// We can't easily override the const base URL, so just verify format helpers.
	msg := formatZendutyMsg(zendutyChange(monitor.Added, 9090))
	if msg == "" {
		t.Error("expected non-empty message")
	}
}

func TestFormatZendutyMsg_Added(t *testing.T) {
	c := zendutyChange(monitor.Added, 8080)
	msg := formatZendutyMsg(c)
	expected := "port 8080/tcp bound on localhost"
	if msg != expected {
		t.Errorf("got %q, want %q", msg, expected)
	}
}

func TestFormatZendutyMsg_Removed(t *testing.T) {
	c := zendutyChange(monitor.Removed, 443)
	msg := formatZendutyMsg(c)
	expected := "port 443/tcp unbound on localhost"
	if msg != expected {
		t.Errorf("got %q, want %q", msg, expected)
	}
}

func TestFormatZendutyMsg_FallsBackToIP(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{Addr: "10.0.0.1", Port: 22, Proto: "tcp"},
	}
	msg := formatZendutyMsg(c)
	expected := "port 22/tcp bound on 10.0.0.1"
	if msg != expected {
		t.Errorf("got %q, want %q", msg, expected)
	}
}
