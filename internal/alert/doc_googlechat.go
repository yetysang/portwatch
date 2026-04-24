// Package alert provides alert handlers for portwatch.
//
// # Google Chat Handler
//
// The GoogleChatHandler sends port-change notifications to a Google Chat space
// using an incoming webhook URL. Each [monitor.Change] is formatted as a plain
// text message and posted individually.
//
// # Configuration
//
//	googlechat:
//	  enabled: true
//	  webhook_url: "https://chat.googleapis.com/v1/spaces/XXXX/messages?key=YYYY&token=ZZZZ"
//	  timeout: 5s
//
// # Message Format
//
// Messages follow the pattern:
//
//	[portwatch] Port <port>/<proto> bound|unbound on <hostname> (process: <name>, pid: <pid>)
//
// The process and PID fields are omitted when not available.
package alert
