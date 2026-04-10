// Package snapshot provides persistence for portwatch port-binding state.
//
// A Store saves and loads Snapshot values — timestamped lists of active
// port bindings — to a JSON file on disk.  The snapshot is written after
// every successful scan so that portwatch can resume without false-positive
// alerts after a restart.
//
// Typical usage:
//
//	store := snapshot.NewStore("/var/lib/portwatch/snapshot.json")
//
//	// On startup, load the previous state.
//	prev, err := store.Load()
//
//	// After each scan, persist the latest bindings.
//	err = store.Save(currentBindings)
package snapshot
