package alert_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/monitor"
)

func signedChange() monitor.Change {
	return monitor.Change{Kind: monitor.Added, Port: 9090, Proto: "tcp", Addr: "0.0.0.0"}
}

func TestWebhookHandler_SignsPayloadWithSecret(t *testing.T) {
	secret := "test-secret"
	var gotSig, gotBody string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		gotSig = r.Header.Get("X-PortWatch-Signature")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := alert.NewWebhookHandler(alert.WebhookConfig{
		URL:     ts.URL,
		Secret:  secret,
		Timeout: 2 * time.Second,
	})

	changes := []monitor.Change{signedChange()}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotSig == "" {
		t.Fatal("expected X-PortWatch-Signature header to be set")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(gotBody))
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if gotSig != expected {
		t.Errorf("signature mismatch:\n got  %s\n want %s", gotSig, expected)
	}
}

func TestWebhookHandler_NoSignatureWithoutSecret(t *testing.T) {
	var gotSig string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSig = r.Header.Get("X-PortWatch-Signature")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := alert.NewWebhookHandler(alert.WebhookConfig{
		URL:     ts.URL,
		Timeout: 2 * time.Second,
	})

	if err := h.Handle([]monitor.Change{signedChange()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotSig != "" {
		t.Errorf("expected no signature header, got %q", gotSig)
	}
}

func TestWebhookHandler_PayloadIsValidJSON(t *testing.T) {
	var body []byte

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := alert.NewWebhookHandler(alert.WebhookConfig{URL: ts.URL})
	changes := []monitor.Change{signedChange()}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !json.Valid(body) {
		t.Errorf("payload is not valid JSON: %s", body)
	}

	if !strings.Contains(string(body), "9090") {
		t.Errorf("payload missing expected port: %s", body)
	}
}
