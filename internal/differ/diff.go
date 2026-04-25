// Package differ provides functionality to compare two sets of environment
// variables and produce a structured diff result.
package differ

// DiffType represents the type of difference for a key.
type DiffType string

const (
	DiffAdded   DiffType = "added"
	DiffRemoved DiffType = "removed"
	DiffChanged DiffType = "changed"
	DiffSame    DiffType = "same"
)

// Entry represents a single diff entry for an environment key.
type Entry struct {
	Key      string
	BaseVal  string
	OtherVal string
	Type     DiffType
}

// Result holds the complete diff between two env maps.
type Result struct {
	Entries []Entry
}

// HasChanges returns true if there are any non-same entries.
func (r *Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Type != DiffSame {
			return true
		}
	}
	return false
}

// Diff compares base and other env maps and returns a Result.
// Keys are reported in a stable order: all keys from both maps, sorted.
func Diff(base, other map[string]string) *Result {
	keys := unionKeys(base, other)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		bv, inBase := base[k]
		ov, inOther := other[k]

		var dt DiffType
		switch {
		case inBase && !inOther:
			dt = DiffRemoved
		case !inBase && inOther:
			dt = DiffAdded
		case bv != ov:
			dt = DiffChanged
		default:
			dt = DiffSame
		}

		entries = append(entries, Entry{
			Key:      k,
			BaseVal:  bv,
			OtherVal: ov,
			Type:     dt,
		})
	}

	return &Result{Entries: entries}
}
