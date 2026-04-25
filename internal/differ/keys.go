package differ

import "sort"

// unionKeys returns a sorted slice of all unique keys from both maps.
func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// FilterByType returns only the entries matching the given DiffType.
func FilterByType(entries []Entry, dt DiffType) []Entry {
	out := make([]Entry, 0)
	for _, e := range entries {
		if e.Type == dt {
			out = append(out, e)
		}
	}
	return out
}

// Summary returns counts of each diff type.
func Summary(entries []Entry) map[DiffType]int {
	counts := map[DiffType]int{
		DiffAdded:   0,
		DiffRemoved: 0,
		DiffChanged: 0,
		DiffSame:    0,
	}
	for _, e := range entries {
		counts[e.Type]++
	}
	return counts
}
