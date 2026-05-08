// Package merger provides functionality to merge multiple .env file maps
// into a single resolved map, with configurable conflict resolution strategies.
package merger

import (
	"fmt"
)

// Strategy defines how conflicting keys are resolved during a merge.
type Strategy int

const (
	// StrategyLastWins keeps the value from the last file that defines the key.
	StrategyLastWins Strategy = iota
	// StrategyFirstWins keeps the value from the first file that defines the key.
	StrategyFirstWins
	// StrategyError returns an error when a conflict is detected.
	StrategyError
)

// Conflict records a key that appeared in more than one source.
type Conflict struct {
	Key    string
	Values []string // one per source, in order
}

// Result holds the merged environment map and any conflicts that were observed.
type Result struct {
	Merged    map[string]string
	Conflicts []Conflict
}

// Merge combines the provided env maps according to the given strategy.
// Sources are processed in order; index 0 is considered the "first" source.
func Merge(sources []map[string]string, strategy Strategy) (*Result, error) {
	merged := make(map[string]string)
	conflictMap := make(map[string][]string)

	for _, src := range sources {
		for k, v := range src {
			existing, exists := merged[k]
			if !exists {
				merged[k] = v
				continue
			}
			// Record the conflict regardless of strategy.
			if len(conflictMap[k]) == 0 {
				conflictMap[k] = []string{existing}
			}
			conflictMap[k] = append(conflictMap[k], v)

			switch strategy {
			case StrategyLastWins:
				merged[k] = v
			case StrategyFirstWins:
				// keep existing — do nothing
			case StrategyError:
				return nil, fmt.Errorf("merger: conflict on key %q (%q vs %q)", k, existing, v)
			}
		}
	}

	var conflicts []Conflict
	for k, vals := range conflictMap {
		conflicts = append(conflicts, Conflict{Key: k, Values: vals})
	}

	return &Result{Merged: merged, Conflicts: conflicts}, nil
}
