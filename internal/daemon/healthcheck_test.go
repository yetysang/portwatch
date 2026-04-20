package daemon

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func defaultHealthCfg(addr string) config.HealthCheckConfig {
	return config.HealthCheckConfig{
		Enabled:     true,
		ListenAddr:  addr,
		Path:        "/healthz",
		ReadTimeout: 5 * time.Second,
	}
}

func TestHealthServer_NotReadyReturns503(t *testing.T) {
	cfg := defaultHealthCfg(":0")
	h := NewHealthServer(cfg)
	// not ready by default
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	h.handle(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rec.Code)
	}
}

func TestHealthServer_ReadyReturns200(t *testing.T) {
	cfg := defaultHealthCfg(":0")
	h := NewHealthServer(cfg)
	h.SetReady(true)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	h.handle(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestHealthServer_ResponseBodyIsJSON(t *testing.T) {
	cfg := defaultHealthCfg(":0")
	h := NewHealthServer(cfg)
	h.SetReady(true)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	h.handle(rec, req)
	body, _ := io.ReadAll(rec.Body)
	var m map[string]string
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if m["status"] != "ok" {
		t.Errorf("unexpected status field: %s", m["status"])
	}
}

func TestHealthServer_DisabledStartIsNoop(t *testing.T) {
	cfg := config.DefaultHealthCheckConfig() // Enabled: false
	h := NewHealthServer(cfg)
	if err := h.Start(); err != nil {
		t.Fatalf("Start on disabled server should not error: %v", err)
	}
	if err := h.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown on disabled server should not error: %v", err)
	}
}

func TestHealthServer_ContentTypeHeader(t *testing.T) {
	cfg := defaultHealthCfg(":0")
	h := NewHealthServer(cfg)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	h.handle(rec, req)
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("unexpected Content-Type: %s", ct)
	}
}
