// Package alert provides alerting handlers for portwatch.
//
// # Slack Handler
//
// SlackHandler sends port-change notifications to a Slack channel via an
// incoming webhook URL. Each Change in a batch produces a separate message.
//
// Usage:
//
//	h := alert.NewSlackHandler("https://hooks.slack.com/services/...", "[portwatch]")
//
// The prefix string is prepended to every message, making it easy to
// identify the source when multiple applications share a channel.
//
// Messages follow the format:
//
//	[prefix] added tcp/8080 (hostname) added [pid 1234, nginx]
//
// If no hostname is resolved the IP address is used as a fallback.
//
// Drain is a no-op; Slack messages are sent synchronously in Handle.
package alert
