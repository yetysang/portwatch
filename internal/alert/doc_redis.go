// Package alert provides alert handlers for portwatch.
//
// # Redis Stream Handler
//
// The RedisHandler publishes port change events to a Redis stream using the
// XADD command. Each change is serialised as a JSON string stored under the
// "event" field, alongside a RFC3339 timestamp in the "ts" field.
//
// Usage:
//
//	client := newRedisClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
//	h := alert.NewRedisHandler(client, cfg.Redis.Stream)
//
// The RedisPublisher interface allows injecting a test double without importing
// a concrete Redis client library, keeping the core package dependency-free.
//
// Configuration is held in config.RedisConfig and can be loaded from the
// standard portwatch TOML file under the [redis] section.
package alert
