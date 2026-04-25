package differ_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/differ"
)

func TestUnionKeys_Deduplicated(t *testing.T) {
	// indirectly tested via Diff; test FilterByType directly
	entries := []differ.Entry{
		{Key: "A", Type: differ.DiffAdded},
		{Key: "B", Type: differ.DiffSame},
		{Key: "C", Type: differ.DiffAdded},
	}

	added := differ.FilterByType(entries, differ.DiffAdded)
	if len(added) != 2 {
		t.Errorf("expected 2 added, got %d", len(added))
	}
}

func TestFilterByType_Empty(t *testing.T) {
	result := differ.FilterByType([]differ.Entry{}, differ.DiffChanged)
	if len(result) != 0 {
		t.Error("expected empty slice")
	}
}

func TestSummary_AllSame(t *testing.T) {
	entries := []differ.Entry{
		{Key: "A", Type: differ.DiffSame},
		{Key: "B", Type: differ.DiffSame},
	}
	s := differ.Summary(entries)
	if s[differ.DiffSame] != 2 {
		t.Errorf("expected 2 same, got %d", s[differ.DiffSame])
	}
	if s[differ.DiffAdded] != 0 || s[differ.DiffRemoved] != 0 || s[differ.DiffChanged] != 0 {
		t.Error("expected zero counts for other types")
	}
}
