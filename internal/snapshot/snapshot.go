// Package snapshot provides functionality to save and load .env diff snapshots
// for baseline comparison across CI runs.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/differ"
)

// Snapshot represents a saved state of a diff result.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label"`
	Entries   []differ.Entry    `json:"entries"`
}

// Save writes a snapshot of the given entries to the specified file path.
func Save(path, label string, entries []differ.Entry) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Entries:   entries,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}

	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot: file not found: %s", path)
		}
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}

	return &snap, nil
}

// Compare returns entries that differ between a live set and a saved snapshot.
// It returns only entries whose Type or values changed relative to the snapshot.
func Compare(live []differ.Entry, snap *Snapshot) []differ.Entry {
	index := make(map[string]differ.Entry, len(snap.Entries))
	for _, e := range snap.Entries {
		index[e.Key] = e
	}

	var delta []differ.Entry
	for _, e := range live {
		prev, found := index[e.Key]
		if !found || prev.Type != e.Type || prev.BaseValue != e.BaseValue || prev.OtherValue != e.OtherValue {
			delta = append(delta, e)
		}
	}

	return delta
}
