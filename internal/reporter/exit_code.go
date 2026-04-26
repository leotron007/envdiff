package reporter

import "github.com/user/envdiff/internal/differ"

// ExitCode returns a CI-friendly exit code based on diff results.
// Returns 0 if all entries are Same, 1 if any differences exist.
func ExitCode(entries []differ.DiffEntry) int {
	for _, e := range entries {
		if e.Type != differ.Same {
			return 1
		}
	}
	return 0
}

// HasDiff reports whether any entry represents a real difference.
func HasDiff(entries []differ.DiffEntry) bool {
	return ExitCode(entries) == 1
}
