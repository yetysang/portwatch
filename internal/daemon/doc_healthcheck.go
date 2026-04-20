// Package daemon provides runtime lifecycle helpers for the portwatch daemon,
// including graceful signal handling and the optional HTTP health-check server.
//
// # Health Check
//
// The HealthServer exposes a single endpoint (default /healthz) that returns
// HTTP 200 with {"status":"ok"} once the daemon has completed its first scan
// and is considered ready, or HTTP 503 with {"status":"starting"} while it is
// still initialising.
//
// Enable and configure the endpoint in portwatch.toml:
//
//	[healthcheck]
//	enabled     = true
//	listen_addr = ":9090"
//	path        = "/healthz"
//	read_timeout = "5s"
//
// The server is started by calling Start() after daemon initialisation and
// shut down via Shutdown(ctx) during graceful termination.
package daemon
