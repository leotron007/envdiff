package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/filter"
)

var sampleEntries = []differ.DiffEntry{
	{Key: "APP_NAME", Type: differ.Same, BaseValue: "app", CompareValue: "app"},
	{Key: "APP_ENV", Type: differ.Changed, BaseValue: "dev", CompareValue: "prod"},
	{Key: "DB_HOST", Type: differ.Added, BaseValue: "", CompareValue: "localhost"},
	{Key: "SECRET_KEY", Type: differ.Removed, BaseValue: "abc", CompareValue: ""},
	{Key: "APP_VERSION", Type: differ.Added, BaseValue: "", CompareValue: "1.2.0"},
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{})
	if len(result) != len(sampleEntries) {
		t.Errorf("expected %d entries, got %d", len(sampleEntries), len(result))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{Prefix: "APP_"})
	if len(result) != 3 {
		t.Errorf("expected 3 entries with prefix APP_, got %d", len(result))
	}
	for _, e := range result {
		if e.Key[:4] != "APP_" {
			t.Errorf("unexpected key without prefix: %s", e.Key)
		}
	}
}

func TestApply_ExcludeKeys(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{Exclude: []string{"SECRET_KEY", "DB_HOST"}})
	for _, e := range result {
		if e.Key == "SECRET_KEY" || e.Key == "DB_HOST" {
			t.Errorf("excluded key %s still present in results", e.Key)
		}
	}
	if len(result) != 3 {
		t.Errorf("expected 3 entries after exclusion, got %d", len(result))
	}
}

func TestApply_OnlyTypes_Added(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{OnlyTypes: []differ.DiffType{differ.Added}})
	if len(result) != 2 {
		t.Errorf("expected 2 added entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Type != differ.Added {
			t.Errorf("expected only Added type, got %s", e.Type)
		}
	}
}

func TestApply_CombinedPrefixAndType(t *testing.T) {
	result := filter.Apply(sampleEntries, filter.Options{
		Prefix:    "APP_",
		OnlyTypes: []differ.DiffType{differ.Added, differ.Changed},
	})
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	result := filter.Apply(nil, filter.Options{Prefix: "APP_"})
	if result != nil && len(result) != 0 {
		t.Errorf("expected empty result for nil input, got %v", result)
	}
}
