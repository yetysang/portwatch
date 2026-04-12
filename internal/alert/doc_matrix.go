// Package alert provides handlers for dispatching port change notifications
// to various alerting backends.
//
// # Matrix Handler
//
// The Matrix handler sends port change events as messages to a Matrix room
// using the Matrix Client-Server API (/_matrix/client/v3/rooms/{roomID}/send).
//
// # Configuration
//
//	[matrix]
//	enabled    = true
//	homeserver = "https://matrix.example.org"
//	token      = "syt_..."
//	room_id    = "!roomid:example.org"
//
// # Message Format
//
// Each change is sent as a plain-text message containing the action (added/removed),
// protocol, address, port, and process name when available.
//
// If multiple changes occur in a single tick they are batched into one message.
package alert
