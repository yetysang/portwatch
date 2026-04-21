// Package alert provides alert handlers for portwatch.
//
// # Mattermost Handler
//
// The MattermostHandler sends port-change notifications to a Mattermost
// channel via an incoming webhook integration.
//
// Configuration (via MattermostConfig):
//
//	[mattermost]
//	enabled  = true
//	url      = "https://mattermost.example.com/hooks/<token>"
//	channel  = "#security-alerts"
//	username = "portwatch"          # optional display name
//	icon_url = "https://..."        # optional icon
//	timeout  = "10s"
//
// Each change produces one POST to the webhook URL with a JSON payload
// compatible with the Mattermost incoming-webhook schema.
//
// The handler implements alert.Handler and can be composed with
// RateLimitedHandler, ThrottleHandler, or any other middleware handler.
package alert
