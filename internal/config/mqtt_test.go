package config

import (
	"testing"
	"time"
)

func TestDefaultMQTTConfig_Values(t *testing.T) {
	cfg := DefaultMQTTConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Broker != "tcp://localhost:1883" {
		t.Errorf("unexpected default broker: %s", cfg.Broker)
	}
	if cfg.Topic != "portwatch/alerts" {
		t.Errorf("unexpected default topic: %s", cfg.Topic)
	}
	if cfg.ClientID != "portwatch" {
		t.Errorf("unexpected default client_id: %s", cfg.ClientID)
	}
	if cfg.QoS != 0 {
		t.Errorf("unexpected default qos: %d", cfg.QoS)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("unexpected default timeout: %v", cfg.Timeout)
	}
}

func TestMQTTConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := MQTTConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got: %v", err)
	}
}

func TestMQTTConfig_ValidateEnabledRequiresBroker(t *testing.T) {
	cfg := DefaultMQTTConfig()
	cfg.Enabled = true
	cfg.Broker = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when broker is empty")
	}
}

func TestMQTTConfig_ValidateEnabledRequiresTopic(t *testing.T) {
	cfg := DefaultMQTTConfig()
	cfg.Enabled = true
	cfg.Topic = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when topic is empty")
	}
}

func TestMQTTConfig_ValidateInvalidQoS(t *testing.T) {
	cfg := DefaultMQTTConfig()
	cfg.Enabled = true
	cfg.QoS = 3
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when qos > 2")
	}
}

func TestMQTTConfig_ValidateEnabledWithValidFields(t *testing.T) {
	cfg := DefaultMQTTConfig()
	cfg.Enabled = true
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error with valid config, got: %v", err)
	}
}
