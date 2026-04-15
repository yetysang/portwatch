package config

import "testing"

func TestDefaultRedisConfig_Values(t *testing.T) {
	c := DefaultRedisConfig()
	if c.Enabled {
		t.Error("expected Enabled=false")
	}
	if c.Addr != "localhost:6379" {
		t.Errorf("unexpected Addr: %s", c.Addr)
	}
	if c.Stream != "portwatch:events" {
		t.Errorf("unexpected Stream: %s", c.Stream)
	}
	if c.DB != 0 {
		t.Errorf("unexpected DB: %d", c.DB)
	}
}

func TestRedisConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	c := RedisConfig{Enabled: false, Addr: "", Stream: ""}
	if err := c.Validate(); err != nil {
		t.Errorf("expected nil error for disabled config, got %v", err)
	}
}

func TestRedisConfig_ValidateEnabledRequiresAddr(t *testing.T) {
	c := DefaultRedisConfig()
	c.Enabled = true
	c.Addr = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty addr")
	}
}

func TestRedisConfig_ValidateEnabledRequiresStream(t *testing.T) {
	c := DefaultRedisConfig()
	c.Enabled = true
	c.Stream = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty stream")
	}
}

func TestRedisConfig_ValidateEnabledNegativeDB(t *testing.T) {
	c := DefaultRedisConfig()
	c.Enabled = true
	c.DB = -1
	if err := c.Validate(); err == nil {
		t.Error("expected error for negative db")
	}
}

func TestRedisConfig_ValidateEnabledWithValidFields(t *testing.T) {
	c := DefaultRedisConfig()
	c.Enabled = true
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
