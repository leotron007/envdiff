package reporter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/reporter"
)

func entries() []differ.DiffEntry {
	return []differ.DiffEntry{
		{Key: "APP_ENV", Type: differ.Added, NewValue: "production"},
		{Key: "DB_HOST", Type: differ.Removed, OldValue: "localhost"},
		{Key: "PORT", Type: differ.Changed, OldValue: "3000", NewValue: "8080"},
		{Key: "LOG_LEVEL", Type: differ.Same, OldValue: "info", NewValue: "info"},
	}
}

func TestWrite_TextFormat_ContainsSymbols(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Write(entries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ APP_ENV=production") {
		t.Errorf("expected added line, got:\n%s", out)
	}
	if !strings.Contains(out, "- DB_HOST=localhost") {
		t.Errorf("expected removed line, got:\n%s", out)
	}
	if !strings.Contains(out, "~ PORT: 3000 -> 8080") {
		t.Errorf("expected changed line, got:\n%s", out)
	}
	if !strings.Contains(out, "  LOG_LEVEL=info") {
		t.Errorf("expected same line, got:\n%s", out)
	}
}

func TestWrite_TextFormat_EmptyEntries(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Write([]differ.DiffEntry{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found.") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestWrite_JSONFormat_ValidStructure(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatJSON)
	if err := r.Write(entries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "[") || !strings.Contains(out, "]") {
		t.Errorf("expected JSON array, got:\n%s", out)
	}
	if !strings.Contains(out, `"key": "APP_ENV"`) {
		t.Errorf("expected APP_ENV key in JSON, got:\n%s", out)
	}
	if !strings.Contains(out, `"type": "added"`) {
		t.Errorf("expected type added in JSON, got:\n%s", out)
	}
}

func TestNew_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic
	r := reporter.New(nil, reporter.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
