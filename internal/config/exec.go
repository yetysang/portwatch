package config

import "fmt"

// ExecConfig holds configuration for the exec alert handler.
type ExecConfig struct {
	// Enabled controls whether the exec handler is active.
	Enabled bool `toml:"enabled" yaml:"enabled"`

	// Command is the executable to run when changes are detected.
	Command string `toml:"command" yaml:"command"`

	// Args are optional static arguments passed before the change summary.
	Args []string `toml:"args" yaml:"args"`
}

// DefaultExecConfig returns an ExecConfig with safe defaults.
func DefaultExecConfig() ExecConfig {
	return ExecConfig{
		Enabled: false,
		Command: "",
		Args:    []string{},
	}
}

// Validate checks that the ExecConfig is consistent.
func (c ExecConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Command == "" {
		return fmt.Errorf("exec: command must not be empty when enabled")
	}
	return nil
}
