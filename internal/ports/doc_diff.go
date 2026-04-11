// Package ports provides scanning, enrichment, filtering, and diffing
// of network port bindings observed on the local host.
//
// # Diff and Baseline
//
// The Diff function computes the set of changes between two snapshots of
// port bindings represented as maps (keyed by "proto:addr:port"). Use
// BindingsToMap to convert a []Binding slice into the required map form.
//
//	prev := ports.BindingsToMap(previousBindings)
//	curr := ports.BindingsToMap(currentBindings)
//	changes := ports.Diff(prev, curr)
//
// Each returned monitor.Change carries the change type (Added/Removed)
// along with the full binding metadata.
//
// The Baseline type provides persistent acknowledgement of expected
// bindings. Entries are stored as JSON on disk and survive restarts.
// A binding present in the baseline can be suppressed from alerting
// by the alert handler or monitor layer.
//
//	bl, err := ports.NewBaseline("/var/lib/portwatch/baseline.json")
//	if err != nil { ... }
//	if !bl.Contains("tcp", "0.0.0.0", 8080) {
//		// alert: unexpected binding
//	}
package ports
