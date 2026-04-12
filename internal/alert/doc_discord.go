// Package alert provides alerting handlers for portwatch.
//
// # Discord Handler
//
// The DiscordHandler sends port change notifications to a Discord channel
// via an incoming webhook URL. Each change is delivered as a rich embed
// with a colour-coded indicator: green for added bindings, red for removed.
//
// Usage:
//
//	h := alert.NewDiscordHandler(cfg.Discord.WebhookURL, nil)
//
// The handler implements the Handler interface and can be composed with
// other handlers via MultiHandler, ThrottleHandler, or DedupHandler.
//
// Configuration is provided through config.DiscordConfig. Set Enabled to
// true and supply a valid WebhookURL to activate Discord alerts.
package alert
