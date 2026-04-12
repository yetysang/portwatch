// Package alert provides alerting handlers for portwatch.
//
// # Microsoft Teams Handler
//
// TeamsHandler sends port-change notifications to a Microsoft Teams channel
// using an Incoming Webhook URL configured in Teams.
//
// # Setup
//
// 1. In Teams, open the channel you want to receive alerts.
// 2. Click "..." → "Connectors" → "Incoming Webhook" → Configure.
// 3. Copy the generated webhook URL.
// 4. Pass it to NewTeamsHandler.
//
// # Message Format
//
// Each change produces one Adaptive Card message:
//
//	[portwatch] Port 8080/tcp added on webserver (pid 1234, process: nginx)
//
// # Usage
//
//	h := alert.NewTeamsHandler("https://outlook.office.com/webhook/...")
//	err := h.Handle(changes)
package alert
