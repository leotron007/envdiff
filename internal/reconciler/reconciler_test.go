package reconciler_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/reconciler"
)

func entries(diffs ...differ.DiffEntry) []differ.DiffEntry { return diffs }

func TestGenerate_AddedEntry(t *testing.T) {
	result := reconciler.Generate(entries(
		differ.DiffEntry{Key: "NEW_KEY", Type: differ.Added, NewValue: "value1"},
	))
	if len(result.Patches) != 1 {
		t.Fatalf("expected 1 patch, got %d", len(result.Patches))
	}
	if result.Patches[0].Action != reconciler.ActionAdd {
		t.Errorf("expected ActionAdd, got %s", result.Patches[0].Action)
	}
	if result.Patches[0].Value != "value1" {
		t.Errorf("unexpected value: %s", result.Patches[0].Value)
	}
}

func TestGenerate_RemovedEntry(t *testing.T) {
	result := reconciler.Generate(entries(
		differ.DiffEntry{Key: "OLD_KEY", Type: differ.Removed, OldValue: "old"},
	))
	if result.Patches[0].Action != reconciler.ActionRemove {
		t.Errorf("expected ActionRemove, got %s", result.Patches[0].Action)
	}
	if result.Patches[0].Value != "" {
		t.Errorf("remove patch should have empty value")
	}
}

func TestGenerate_ChangedEntry(t *testing.T) {
	result := reconciler.Generate(entries(
		differ.DiffEntry{Key: "KEY", Type: differ.Changed, OldValue: "a", NewValue: "b"},
	))
	if result.Patches[0].Action != reconciler.ActionUpdate {
		t.Errorf("expected ActionUpdate, got %s", result.Patches[0].Action)
	}
	if result.Patches[0].Value != "b" {
		t.Errorf("expected new value 'b', got %s", result.Patches[0].Value)
	}
}

func TestGenerate_SameEntrySkipped(t *testing.T) {
	result := reconciler.Generate(entries(
		differ.DiffEntry{Key: "KEY", Type: differ.Same, OldValue: "x", NewValue: "x"},
	))
	if len(result.Patches) != 0 {
		t.Errorf("expected no patches for Same entries, got %d", len(result.Patches))
	}
}

func TestGenerate_StableOrder(t *testing.T) {
	result := reconciler.Generate(entries(
		differ.DiffEntry{Key: "Z_KEY", Type: differ.Added, NewValue: "z"},
		differ.DiffEntry{Key: "A_KEY", Type: differ.Added, NewValue: "a"},
	))
	if result.Patches[0].Key != "A_KEY" {
		t.Errorf("expected A_KEY first, got %s", result.Patches[0].Key)
	}
}

func TestRenderPatch_ContainsExpectedLines(t *testing.T) {
	result := reconciler.Generate(entries(
		differ.DiffEntry{Key: "ADD_ME", Type: differ.Added, NewValue: "hello"},
		differ.DiffEntry{Key: "DROP_ME", Type: differ.Removed, OldValue: "bye"},
		differ.DiffEntry{Key: "UPD_ME", Type: differ.Changed, OldValue: "old", NewValue: "new"},
	))
	output := reconciler.RenderPatch(result)
	if !strings.Contains(output, "ADD_ME=hello") {
		t.Errorf("missing ADD_ME line in output")
	}
	if !strings.Contains(output, "# REMOVE: DROP_ME") {
		t.Errorf("missing REMOVE comment in output")
	}
	if !strings.Contains(output, "UPD_ME=new") {
		t.Errorf("missing UPD_ME update line in output")
	}
}
