// Package alert provides handlers for dispatching port change notifications
// to various alerting backends.
//
// # Ntfy Handler
//
// NewNtfyHandler sends push notifications via a self-hosted or public ntfy.sh
// server whenever port bindings are added or removed.
//
// Configuration is provided through config.NtfyConfig:
//
//	- URL      – base URL of the ntfy server (e.g. https://ntfy.sh)
//	- Topic    – the topic to publish messages to
//	- Token    – optional Bearer token for authenticated servers
//	- Priority – message priority: min, low, default, high, or urgent
//	- Enabled  – set to false to disable without removing config
//
// Example config.yaml snippet:
//
//	ntfy:
//	  enabled: true
//	  url: https://ntfy.sh
//	  topic: portwatch-alerts
//	  priority: high
//	  token: tk_mytoken
package alert
