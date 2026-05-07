package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/snapshot"
)

func sampleEntries() []differ.Entry {
	return []differ.Entry{
		{Key: "APP_ENV", Type: differ.Same, BaseValue: "production", OtherValue: "production"},
		{Key: "DB_HOST", Type: differ.Changed, BaseValue: "localhost", OtherValue: "db.prod.internal"},
		{Key: "NEW_KEY", Type: differ.Added, BaseValue: "", OtherValue: "newval"},
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	err := snapshot.Save(path, "test-label", sampleEntries())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected snapshot file to exist")
	}
}

func TestLoad_ReturnsSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	original := sampleEntries()
	if err := snapshot.Save(path, "my-label", original); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if snap.Label != "my-label" {
		t.Errorf("expected label 'my-label', got %q", snap.Label)
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if snap.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("CreatedAt is in the future")
	}
	if len(snap.Entries) != len(original) {
		t.Errorf("expected %d entries, got %d", len(original), len(snap.Entries))
	}
}

func TestLoad_NotFound_ReturnsError(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCompare_ReturnsDelta(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	base := sampleEntries()
	if err := snapshot.Save(path, "base", base); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	snap, _ := snapshot.Load(path)

	live := []differ.Entry{
		{Key: "APP_ENV", Type: differ.Same, BaseValue: "production", OtherValue: "production"},
		{Key: "DB_HOST", Type: differ.Changed, BaseValue: "localhost", OtherValue: "db.new.internal"},
		{Key: "BRAND_NEW", Type: differ.Added, BaseValue: "", OtherValue: "surprise"},
	}

	delta := snapshot.Compare(live, snap)
	if len(delta) != 2 {
		t.Errorf("expected 2 delta entries, got %d", len(delta))
	}
}

func TestCompare_NoChanges_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	entries := sampleEntries()
	_ = snapshot.Save(path, "base", entries)
	snap, _ := snapshot.Load(path)

	delta := snapshot.Compare(entries, snap)
	if len(delta) != 0 {
		t.Errorf("expected 0 delta entries, got %d", len(delta))
	}
}
