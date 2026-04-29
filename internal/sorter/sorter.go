// Package sorter provides utilities for sorting and grouping diff entries
// by various criteria such as key name, diff type, or prefix.
package sorter

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// SortBy defines the sort strategy for diff entries.
type SortBy int

const (
	// SortByKey sorts entries alphabetically by key name.
	SortByKey SortBy = iota
	// SortByType sorts entries by diff type (added, removed, changed, same).
	SortByType
	// SortByPrefix groups entries by their prefix (e.g., DB_, APP_).
	SortByPrefix
)

// typeOrder defines the display priority for diff types.
var typeOrder = map[differ.DiffType]int{
	differ.Added:   0,
	differ.Removed: 1,
	differ.Changed: 2,
	differ.Same:    3,
}

// Sort returns a sorted copy of the given diff entries using the specified strategy.
func Sort(entries []differ.Entry, by SortBy) []differ.Entry {
	result := make([]differ.Entry, len(entries))
	copy(result, entries)

	switch by {
	case SortByKey:
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Key < result[j].Key
		})
	case SortByType:
		sort.SliceStable(result, func(i, j int) bool {
			oi := typeOrder[result[i].Type]
			oj := typeOrder[result[j].Type]
			if oi != oj {
				return oi < oj
			}
			return result[i].Key < result[j].Key
		})
	case SortByPrefix:
		sort.SliceStable(result, func(i, j int) bool {
			pi := extractPrefix(result[i].Key)
			pj := extractPrefix(result[j].Key)
			if pi != pj {
				return pi < pj
			}
			return result[i].Key < result[j].Key
		})
	}

	return result
}

// GroupByPrefix returns a map of prefix -> entries for the given diff entries.
func GroupByPrefix(entries []differ.Entry) map[string][]differ.Entry {
	groups := make(map[string][]differ.Entry)
	for _, e := range entries {
		prefix := extractPrefix(e.Key)
		groups[prefix] = append(groups[prefix], e)
	}
	return groups
}

// extractPrefix returns the prefix of a key (the part before the first underscore).
// If no underscore is found, the full key is returned as the prefix.
func extractPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
