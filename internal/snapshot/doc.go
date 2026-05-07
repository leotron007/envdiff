// Package snapshot provides save/load functionality for persisting diff
// results as JSON snapshots. Snapshots can be used as baselines in CI
// pipelines to detect regressions or unexpected environment changes between
// runs.
//
// Usage:
//
//	// Save a snapshot after diffing
//	err := snapshot.Save(".envdiff-snapshot.json", "ci-run-42", entries)
//
//	// Load a previous snapshot for comparison
//	snap, err := snapshot.Load(".envdiff-snapshot.json")
//
//	// Compare live diff entries against the snapshot baseline
//	delta := snapshot.Compare(liveEntries, snap)
package snapshot
