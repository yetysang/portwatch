package config

import "testing"

func TestDefaultMatrixConfig_Values(t *testing.T) {
	c := DefaultMatrixConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false")
	}
	if c.Homeserver != "" {
		t.Errorf("unexpected Homeserver: %q", c.Homeserver)
	}
	if c.Token != "" {
		t.Errorf("unexpected Token: %q", c.Token)
	}
	if c.RoomID != "" {
		t.Errorf("unexpected RoomID: %q", c.RoomID)
	}
}

func TestMatrixConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	c := MatrixConfig{Enabled: false}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got: %v", err)
	}
}

func TestMatrixConfig_ValidateEnabledRequiresHomeserver(t *testing.T) {
	c := MatrixConfig{Enabled: true, Token: "tok", RoomID: "!room:example.com"}
	if err := c.Validate(); err == nil {
		t.Error("expected error when homeserver is empty")
	}
}

func TestMatrixConfig_ValidateEnabledRequiresToken(t *testing.T) {
	c := MatrixConfig{Enabled: true, Homeserver: "https://matrix.example.com", RoomID: "!room:example.com"}
	if err := c.Validate(); err == nil {
		t.Error("expected error when token is empty")
	}
}

func TestMatrixConfig_ValidateEnabledRequiresRoomID(t *testing.T) {
	c := MatrixConfig{Enabled: true, Homeserver: "https://matrix.example.com", Token: "tok"}
	if err := c.Validate(); err == nil {
		t.Error("expected error when room_id is empty")
	}
}

func TestMatrixConfig_ValidateFullConfig(t *testing.T) {
	c := MatrixConfig{
		Enabled:    true,
		Homeserver: "https://matrix.example.com",
		Token:      "syt_abc123",
		RoomID:     "!abcdef:example.com",
	}
	if err := c.Validate(); err != nil {
		t.Errorf("expected valid config, got error: %v", err)
	}
}
