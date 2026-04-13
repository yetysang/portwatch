package config

import "testing"

func TestDefaultExecConfig_Values(t *testing.T) {
	cfg := DefaultExecConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.Command != "" {
		t.Errorf("expected empty Command, got %q", cfg.Command)
	}
	if len(cfg.Args) != 0 {
		t.Errorf("expected empty Args slice, got %v", cfg.Args)
	}
}

func TestExecConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := ExecConfig{Enabled: false, Command: ""}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestExecConfig_ValidateEnabledRequiresCommand(t *testing.T) {
	cfg := ExecConfig{Enabled: true, Command: ""}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error when command is empty")
	}
}

func TestExecConfig_ValidateEnabledWithCommand(t *testing.T) {
	cfg := ExecConfig{Enabled: true, Command: "/usr/local/bin/notify.sh"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error with valid command, got %v", err)
	}
}

func TestExecConfig_ValidateEnabledWithArgs(t *testing.T) {
	cfg := ExecConfig{
		Enabled: true,
		Command: "/usr/bin/notify-send",
		Args:    []string{"portwatch"},
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error with args, got %v", err)
	}
}

func TestExecConfig_ValidateErrorMessage(t *testing.T) {
	cfg := ExecConfig{Enabled: true, Command: ""}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}
	if err.Error() == "" {
		t.Error("error message should not be empty")
	}
}
