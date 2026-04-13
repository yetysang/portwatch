// Package alert provides alerting handlers for portwatch.
//
// # Exec Handler
//
// The ExecHandler invokes an external command whenever port binding
// changes are detected. A formatted summary of the changes is appended
// as the final argument to the command.
//
// Example configuration (TOML):
//
//	[exec]
//	enabled = true
//	command = "/usr/local/bin/notify.sh"
//	args    = ["--level", "warn"]
//
// The handler runs the command synchronously and returns an error if
// the process exits with a non-zero status. Stderr output is included
// in the error message for easier debugging.
//
// Use this handler to integrate portwatch with custom scripts,
// notification systems, or orchestration tools that accept CLI input.
package alert
