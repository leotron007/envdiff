// Package filter provides functionality to filter diff entries
// based on key patterns, prefixes, or explicit exclusion lists.
package filter

import (
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Options holds configuration for filtering diff entries.
type Options struct {
	// Prefix restricts results to keys starting with this prefix.
	Prefix string

	// Exclude is a list of exact key names to omit from results.
	Exclude []string

	// OnlyTypes restricts results to specific diff types (e.g. "added", "removed", "changed").
	OnlyTypes []differ.DiffType
}

// Apply returns a filtered slice of DiffEntry values based on the given Options.
func Apply(entries []differ.DiffEntry, opts Options) []differ.DiffEntry {
	excludeSet := make(map[string]struct{}, len(opts.Exclude))
	for _, k := range opts.Exclude {
		excludeSet[k] = struct{}{}
	}

	typeSet := make(map[differ.DiffType]struct{}, len(opts.OnlyTypes))
	for _, t := range opts.OnlyTypes {
		typeSet[t] = struct{}{}
	}

	var result []differ.DiffEntry
	for _, e := range entries {
		if _, excluded := excludeSet[e.Key]; excluded {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if len(typeSet) > 0 {
			if _, ok := typeSet[e.Type]; !ok {
				continue
			}
		}
		result = append(result, e)
	}
	return result
}
