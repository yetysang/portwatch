// Package alert provides alerting handlers for portwatch.
//
// # New Relic Handler
//
// The New Relic handler sends port change events to New Relic Insights
// using the custom events API. Each binding change (added or removed)
// is submitted as a "PortWatchEvent" event type.
//
// Configuration:
//
//	newrelic:
//	  enabled: true
//	  api_key: "YOUR_INSERT_KEY"
//	  account_id: "1234567"
//	  region: "US"   # or "EU"
//	  timeout: 5s
//
// Events include the following attributes:
//   - eventType: always "PortWatchEvent"
//   - action: "added" or "removed"
//   - port: the port number
//   - proto: "tcp" or "udp"
//   - addr: resolved hostname or IP address
//   - process: process name (if available)
//   - pid: process ID (if available)
//   - timestamp: Unix millisecond timestamp
package alert
