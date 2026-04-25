package differ_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/differ"
)

func TestDiff_Added(t *testing.T) {
	base := map[string]string{"A": "1"}
	other := map[string]string{"A": "1", "B": "2"}

	result := differ.Diff(base, other)
	if !result.HasChanges() {
		t.Fatal("expected changes")
	}
	added := differ.FilterByType(result.Entries, differ.DiffAdded)
	if len(added) != 1 || added[0].Key != "B" {
		t.Errorf("expected B to be added, got %+v", added)
	}
}

func TestDiff_Removed(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	other := map[string]string{"A": "1"}

	result := differ.Diff(base, other)
	removed := differ.FilterByType(result.Entries, differ.DiffRemoved)
	if len(removed) != 1 || removed[0].Key != "B" {
		t.Errorf("expected B removed, got %+v", removed)
	}
}

func TestDiff_Changed(t *testing.T) {
	base := map[string]string{"A": "old"}
	other := map[string]string{"A": "new"}

	result := differ.Diff(base, other)
	changed := differ.FilterByType(result.Entries, differ.DiffChanged)
	if len(changed) != 1 {
		t.Fatalf("expected 1 changed entry, got %d", len(changed))
	}
	if changed[0].BaseVal != "old" || changed[0].OtherVal != "new" {
		t.Errorf("unexpected values: %+v", changed[0])
	}
}

func TestDiff_Same(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	other := map[string]string{"A": "1", "B": "2"}

	result := differ.Diff(base, other)
	if result.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestDiff_StableOrder(t *testing.T) {
	base := map[string]string{"Z": "1", "A": "2", "M": "3"}
	other := map[string]string{"Z": "1", "A": "2", "M": "3"}

	result := differ.Diff(base, other)
	expected := []string{"A", "M", "Z"}
	for i, e := range result.Entries {
		if e.Key != expected[i] {
			t.Errorf("position %d: expected %s got %s", i, expected[i], e.Key)
		}
	}
}

func TestSummary(t *testing.T) {
	base := map[string]string{"A": "1", "B": "old"}
	other := map[string]string{"B": "new", "C": "3"}

	result := differ.Diff(base, other)
	summary := differ.Summary(result.Entries)

	if summary[differ.DiffAdded] != 1 {
		t.Errorf("added: want 1 got %d", summary[differ.DiffAdded])
	}
	if summary[differ.DiffRemoved] != 1 {
		t.Errorf("removed: want 1 got %d", summary[differ.DiffRemoved])
	}
	if summary[differ.DiffChanged] != 1 {
		t.Errorf("changed: want 1 got %d", summary[differ.DiffChanged])
	}
}
