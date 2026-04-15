package config

import (
	"testing"
)

func TestDefaultKafkaConfig_Values(t *testing.T) {
	cfg := DefaultKafkaConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if len(cfg.Brokers) != 1 || cfg.Brokers[0] != "localhost:9092" {
		t.Errorf("unexpected default brokers: %v", cfg.Brokers)
	}
	if cfg.Topic != "portwatch" {
		t.Errorf("unexpected default topic: %s", cfg.Topic)
	}
	if cfg.ClientID != "portwatch" {
		t.Errorf("unexpected default client_id: %s", cfg.ClientID)
	}
}

func TestKafkaConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := KafkaConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got: %v", err)
	}
}

func TestKafkaConfig_ValidateEnabledRequiresBrokers(t *testing.T) {
	cfg := KafkaConfig{Enabled: true, Brokers: nil, Topic: "portwatch", ClientID: "portwatch"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when brokers is empty")
	}
}

func TestKafkaConfig_ValidateEnabledRequiresNonEmptyBroker(t *testing.T) {
	cfg := KafkaConfig{Enabled: true, Brokers: []string{""}, Topic: "portwatch", ClientID: "portwatch"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when a broker entry is empty")
	}
}

func TestKafkaConfig_ValidateEnabledRequiresTopic(t *testing.T) {
	cfg := KafkaConfig{Enabled: true, Brokers: []string{"localhost:9092"}, Topic: "", ClientID: "portwatch"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when topic is empty")
	}
}

func TestKafkaConfig_ValidateEnabledRequiresClientID(t *testing.T) {
	cfg := KafkaConfig{Enabled: true, Brokers: []string{"localhost:9092"}, Topic: "portwatch", ClientID: ""}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when client_id is empty")
	}
}

func TestKafkaConfig_ValidateEnabledWithAllFields(t *testing.T) {
	cfg := KafkaConfig{
		Enabled:  true,
		Brokers:  []string{"broker1:9092", "broker2:9092"},
		Topic:    "alerts",
		ClientID: "portwatch-prod",
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for valid config, got: %v", err)
	}
}
