package merger_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/merger"
)

func TestMerge_NoConflict_CombinesAll(t *testing.T) {
	a := map[string]string{"APP_NAME": "foo", "PORT": "8080"}
	b := map[string]string{"DB_HOST": "localhost"}

	res, err := merger.Merge([]map[string]string{a, b}, merger.StrategyLastWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Merged) != 3 {
		t.Errorf("expected 3 keys, got %d", len(res.Merged))
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
}

func TestMerge_LastWins_OverwritesValue(t *testing.T) {
	a := map[string]string{"PORT": "8080"}
	b := map[string]string{"PORT": "9090"}

	res, err := merger.Merge([]map[string]string{a, b}, merger.StrategyLastWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["PORT"] != "9090" {
		t.Errorf("expected 9090, got %s", res.Merged["PORT"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMerge_FirstWins_KeepsOriginal(t *testing.T) {
	a := map[string]string{"PORT": "8080"}
	b := map[string]string{"PORT": "9090"}

	res, err := merger.Merge([]map[string]string{a, b}, merger.StrategyFirstWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["PORT"] != "8080" {
		t.Errorf("expected 8080, got %s", res.Merged["PORT"])
	}
}

func TestMerge_StrategyError_ReturnsError(t *testing.T) {
	a := map[string]string{"SECRET": "abc"}
	b := map[string]string{"SECRET": "xyz"}

	_, err := merger.Merge([]map[string]string{a, b}, merger.StrategyError)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMerge_EmptySources_ReturnsEmpty(t *testing.T) {
	res, err := merger.Merge([]map[string]string{}, merger.StrategyLastWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Merged) != 0 {
		t.Errorf("expected empty map, got %d keys", len(res.Merged))
	}
}

func TestMerge_MultipleConflicts_RecordsAll(t *testing.T) {
	a := map[string]string{"A": "1", "B": "x"}
	b := map[string]string{"A": "2", "B": "y"}
	c := map[string]string{"A": "3"}

	res, err := merger.Merge([]map[string]string{a, b, c}, merger.StrategyLastWins)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["A"] != "3" {
		t.Errorf("expected A=3, got %s", res.Merged["A"])
	}
	if len(res.Conflicts) != 2 {
		t.Errorf("expected 2 conflicts, got %d", len(res.Conflicts))
	}
}
