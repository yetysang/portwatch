package config

import "testing"

func TestDefaultAMQPConfig_Values(t *testing.T) {
	c := DefaultAMQPConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if c.URL == "" {
		t.Error("expected non-empty default URL")
	}
	if c.Exchange == "" {
		t.Error("expected non-empty default Exchange")
	}
	if c.RoutingKey == "" {
		t.Error("expected non-empty default RoutingKey")
	}
	if c.VHost == "" {
		t.Error("expected non-empty default VHost")
	}
}

func TestAMQPConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	c := AMQPConfig{Enabled: false}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got: %v", err)
	}
}

func TestAMQPConfig_ValidateEnabledRequiresURL(t *testing.T) {
	c := AMQPConfig{Enabled: true, Exchange: "x", RoutingKey: "r"}
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing URL")
	}
}

func TestAMQPConfig_ValidateEnabledRequiresExchange(t *testing.T) {
	c := AMQPConfig{Enabled: true, URL: "amqp://localhost", RoutingKey: "r"}
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing Exchange")
	}
}

func TestAMQPConfig_ValidateEnabledRequiresRoutingKey(t *testing.T) {
	c := AMQPConfig{Enabled: true, URL: "amqp://localhost", Exchange: "x"}
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing RoutingKey")
	}
}

func TestAMQPConfig_ValidateEnabledWithAllFields(t *testing.T) {
	c := AMQPConfig{
		Enabled:    true,
		URL:        "amqp://guest:guest@localhost:5672/",
		Exchange:   "portwatch",
		RoutingKey: "port.events",
	}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error for valid config, got: %v", err)
	}
}
