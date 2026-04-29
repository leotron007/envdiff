package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/sorter"
)

var testEntries = []differ.Entry{
	{Key: "DB_HOST", Type: differ.Same},
	{Key: "APP_SECRET", Type: differ.Added},
	{Key: "DB_PORT", Type: differ.Changed},
	{Key: "APP_NAME", Type: differ.Removed},
	{Key: "LOG_LEVEL", Type: differ.Added},
}

func TestSort_ByKey(t *testing.T) {
	result := sorter.Sort(testEntries, sorter.SortByKey)
	expected := []string{"APP_NAME", "APP_SECRET", "DB_HOST", "DB_PORT", "LOG_LEVEL"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestSort_ByType(t *testing.T) {
	result := sorter.Sort(testEntries, sorter.SortByType)
	// Added entries should come first
	if result[0].Type != differ.Added {
		t.Errorf("expected first entry to be Added, got %v", result[0].Type)
	}
	// Same entries should come last
	if result[len(result)-1].Type != differ.Same {
		t.Errorf("expected last entry to be Same, got %v", result[len(result)-1].Type)
	}
}

func TestSort_ByPrefix(t *testing.T) {
	result := sorter.Sort(testEntries, sorter.SortByPrefix)
	// APP_ entries should be first
	if !startsWith(result[0].Key, "APP_") {
		t.Errorf("expected first entry to have APP_ prefix, got %q", result[0].Key)
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	original := make([]differ.Entry, len(testEntries))
	copy(original, testEntries)
	sorter.Sort(testEntries, sorter.SortByKey)
	for i, e := range testEntries {
		if e.Key != original[i].Key {
			t.Errorf("original slice was mutated at index %d", i)
		}
	}
}

func TestGroupByPrefix(t *testing.T) {
	groups := sorter.GroupByPrefix(testEntries)
	if len(groups["APP"]) != 2 {
		t.Errorf("expected 2 APP entries, got %d", len(groups["APP"]))
	}
	if len(groups["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(groups["DB"]))
	}
	if len(groups["LOG"]) != 1 {
		t.Errorf("expected 1 LOG entry, got %d", len(groups["LOG"]))
	}
}

func TestGroupByPrefix_NoUnderscore(t *testing.T) {
	entries := []differ.Entry{
		{Key: "PORT", Type: differ.Same},
		{Key: "HOST", Type: differ.Added},
	}
	groups := sorter.GroupByPrefix(entries)
	if _, ok := groups["PORT"]; !ok {
		t.Error("expected group with key PORT")
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
