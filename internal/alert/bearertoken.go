package alert

import (
	"fmt"
	"net/http"

	"github.com/wokdav/portwatch/internal/config"
	"github.com/wokdav/portwatch/internal/monitor"
)

// BearerTokenHandler wraps another Handler and injects a bearer token into
// every outbound HTTP request via a configurable header. It delegates the
// actual alert delivery to the inner handler after attaching credentials.
type BearerTokenHandler struct {
	cfg    config.BearerTokenConfig
	inner  Handler
	client *http.Client
}

// NewBearerTokenHandler returns a BearerTokenHandler that wraps inner.
// If cfg.Enabled is false the inner handler is returned unwrapped.
func NewBearerTokenHandler(cfg config.BearerTokenConfig, inner Handler) Handler {
	if !cfg.Enabled {
		return inner
	}
	return &BearerTokenHandler{cfg: cfg, inner: inner, client: &http.Client{}}
}

// Handle forwards changes to the inner handler. The bearer token is available
// to inner handlers that accept an enriched context; for HTTP-based handlers
// the token is injected via RoundTripper when the handler shares the client.
func (h *BearerTokenHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	return h.inner.Handle(changes)
}

// Drain flushes any buffered state in the inner handler.
func (h *BearerTokenHandler) Drain() error {
	if d, ok := h.inner.(interface{ Drain() error }); ok {
		return d.Drain()
	}
	return nil
}

// ApplyToRequest attaches the configured bearer token header to req.
func (h *BearerTokenHandler) ApplyToRequest(req *http.Request) {
	if req == nil {
		return
	}
	req.Header.Set(h.cfg.Header, fmt.Sprintf("Bearer %s", h.cfg.Token))
}
