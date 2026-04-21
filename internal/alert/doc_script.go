// Package alert provides alert handlers for portwatch.
//
// # Script Handler
//
// The ScriptHandler executes an external script for every port-binding change
// detected by portwatch. The change is serialised as a JSON object and written
// to the script's stdin, making it easy to integrate with any language or
// toolchain.
//
// Example JSON payload delivered to stdin:
//
//	{
//	  "kind":    "added",
//	  "proto":   "tcp",
//	  "addr":    "0.0.0.0",
//	  "port":    9090,
//	  "process": "prometheus",
//	  "pid":     4321
//	}
//
// Configuration fields (ScriptConfig):
//
//	enabled   – enable the handler (default: false)
//	path      – absolute path to the executable (required when enabled)
//	args      – optional extra arguments passed to the script
//	timeout   – per-invocation timeout, default 10 s, max 5 m
//	env_vars  – additional KEY=VALUE environment variables injected at runtime
//
// The handler returns an error if the script exits with a non-zero status or
// if the timeout is exceeded, allowing upstream retry / circuit-breaker
// wrappers to react accordingly.
package alert
