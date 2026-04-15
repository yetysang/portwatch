package config

import "fmt"

// RedisConfig holds configuration for the Redis stream alert handler.
type RedisConfig struct {
	Enabled bool   `toml:"enabled"`
	Addr    string `toml:"addr"`
	Password string `toml:"password"`
	DB      int    `toml:"db"`
	Stream  string `toml:"stream"`
}

// DefaultRedisConfig returns a RedisConfig with sensible defaults.
func DefaultRedisConfig() RedisConfig {
	return RedisConfig{
		Enabled: false,
		Addr:    "localhost:6379",
		Password: "",
		DB:      0,
		Stream:  "portwatch:events",
	}
}

// Validate returns an error if the RedisConfig is invalid.
func (c RedisConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Addr == "" {
		return fmt.Errorf("redis: addr must not be empty")
	}
	if c.Stream == "" {
		return fmt.Errorf("redis: stream must not be empty")
	}
	if c.DB < 0 {
		return fmt.Errorf("redis: db must be >= 0")
	}
	return nil
}
