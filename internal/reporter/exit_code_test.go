package reporter_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/reporter"
)

func TestExitCode_NoDiff_ReturnsZero(t *testing.T) {
	ents := []differ.DiffEntry{
		{Key: "A", Type: differ.Same, OldValue: "x", NewValue: "x"},
		{Key: "B", Type: differ.Same, OldValue: "y", NewValue: "y"},
	}
	if code := reporter.ExitCode(ents); code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestExitCode_WithAdded_ReturnsOne(t *testing.T) {
	ents := []differ.DiffEntry{
		{Key: "NEW_KEY", Type: differ.Added, NewValue: "value"},
	}
	if code := reporter.ExitCode(ents); code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestExitCode_WithRemoved_ReturnsOne(t *testing.T) {
	ents := []differ.DiffEntry{
		{Key: "OLD_KEY", Type: differ.Removed, OldValue: "gone"},
	}
	if code := reporter.ExitCode(ents); code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestExitCode_WithChanged_ReturnsOne(t *testing.T) {
	ents := []differ.DiffEntry{
		{Key: "PORT", Type: differ.Changed, OldValue: "3000", NewValue: "8080"},
	}
	if code := reporter.ExitCode(ents); code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestExitCode_Empty_ReturnsZero(t *testing.T) {
	if code := reporter.ExitCode([]differ.DiffEntry{}); code != 0 {
		t.Errorf("expected exit code 0 for empty entries, got %d", code)
	}
}

func TestHasDiff_Mixed_ReturnsTrue(t *testing.T) {
	ents := []differ.DiffEntry{
		{Key: "A", Type: differ.Same},
		{Key: "B", Type: differ.Added, NewValue: "new"},
	}
	if !reporter.HasDiff(ents) {
		t.Error("expected HasDiff to return true")
	}
}

func TestHasDiff_AllSame_ReturnsFalse(t *testing.T) {
	ents := []differ.DiffEntry{
		{Key: "A", Type: differ.Same},
	}
	if reporter.HasDiff(ents) {
		t.Error("expected HasDiff to return false")
	}
}
