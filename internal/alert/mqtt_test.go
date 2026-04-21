package alert

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func mqttChange(kind monitor.ChangeKind, port int) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Proto:    "tcp",
			Addr:     "0.0.0.0",
			Port:     port,
			Process:  "sshd",
			Hostname: "host1",
		},
	}
}

func TestFormatMQTTMsg_Added(t *testing.T) {
	ch := mqttChange(monitor.Added, 22)
	data, err := formatMQTTMsg(ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var msg mqttMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if msg.Action != "added" {
		t.Errorf("expected action=added, got %s", msg.Action)
	}
	if msg.Port != 22 {
		t.Errorf("expected port=22, got %d", msg.Port)
	}
	if msg.Proto != "tcp" {
		t.Errorf("expected proto=tcp, got %s", msg.Proto)
	}
	if msg.Hostname != "host1" {
		t.Errorf("expected hostname=host1, got %s", msg.Hostname)
	}
}

func TestFormatMQTTMsg_Removed(t *testing.T) {
	ch := mqttChange(monitor.Removed, 8080)
	data, err := formatMQTTMsg(ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var msg mqttMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if msg.Action != "removed" {
		t.Errorf("expected action=removed, got %s", msg.Action)
	}
	if msg.Port != 8080 {
		t.Errorf("expected port=8080, got %d", msg.Port)
	}
}

func TestFormatMQTTMsg_ContainsTimestamp(t *testing.T) {
	ch := mqttChange(monitor.Added, 443)
	data, err := formatMQTTMsg(ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var msg mqttMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if msg.At == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestMQTTHandler_HandleEmptyChangesIsNoop(t *testing.T) {
	// Verify Handle returns nil with no changes without requiring a real broker.
	h := &MQTTHandler{}
	if err := h.Handle(context.Background(), nil); err != nil {
		t.Errorf("expected nil error on empty changes, got: %v", err)
	}
	if err := h.Handle(context.Background(), []monitor.Change{}); err != nil {
		t.Errorf("expected nil error on empty slice, got: %v", err)
	}
}
