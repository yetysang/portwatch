// Package alert provides handlers for dispatching port-change notifications
// to various backends.
//
// # Prometheus Handler
//
// The Prometheus handler exposes port-change metrics via an HTTP endpoint that
// can be scraped by a Prometheus server. It tracks two counters:
//
//   - portwatch_bindings_added_total   – incremented for every newly detected
//     port binding.
//   - portwatch_bindings_removed_total – incremented for every port binding
//     that has disappeared.
//
// The handler embeds a lightweight HTTP server; call ServeHTTP directly or
// register it with your own mux.
//
// # Configuration
//
// Enable the handler via the [config.PrometheusConfig] block:
//
//	[prometheus]
//	enabled     = true
//	listen_addr = ":9090"
//	path        = "/metrics"
//
// Fields:
//
//	enabled     – whether the handler is active (default: false).
//	listen_addr – TCP address the metrics endpoint listens on.
//	path        – URL path for the /metrics endpoint (must start with "/").
//
// # Example Prometheus scrape config
//
//	scrape_configs:
//	  - job_name: portwatch
//	    static_configs:
//	      - targets: ["localhost:9090"]
//	    metrics_path: /metrics
package alert
