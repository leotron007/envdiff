package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/envdiff/envdiff/internal/audit"
	"github.com/envdiff/envdiff/internal/differ"
)

func TestLogDiff_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	entries := []differ.Entry{
		{Key: "FOO", Type: differ.Added},
		{Key: "BAR", Type: differ.Removed},
		{Key: "BAZ", Type: differ.Same},
	}

	if err := l.LogDiff([]string{".env", ".env.prod"}, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var ev audit.Event
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &ev); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if ev.Type != audit.EventDiff {
		t.Errorf("expected type %q, got %q", audit.EventDiff, ev.Type)
	}
	if ev.Entries != 3 {
		t.Errorf("expected 3 entries, got %d", ev.Entries)
	}
	if !strings.Contains(ev.Summary, "+1") {
		t.Errorf("summary should mention added count, got %q", ev.Summary)
	}
}

func TestLogReconcile_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.LogReconcile([]string{".env"}, 7); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var ev audit.Event
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &ev); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if ev.Type != audit.EventReconcile {
		t.Errorf("expected type %q, got %q", audit.EventReconcile, ev.Type)
	}
	if ev.Entries != 7 {
		t.Errorf("expected 7 patch lines, got %d", ev.Entries)
	}
}

func TestLogSecretRead_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	keys := []string{"API_KEY", "DB_PASSWORD"}
	if err := l.LogSecretRead(keys); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var ev audit.Event
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &ev); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if ev.Type != audit.EventSecretRead {
		t.Errorf("expected type %q, got %q", audit.EventSecretRead, ev.Type)
	}
	if ev.Entries != 2 {
		t.Errorf("expected 2 entries, got %d", ev.Entries)
	}
}

func TestNew_NilWriter_DoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked: %v", r)
		}
	}()
	// Should default to os.Stderr without panicking.
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestTimestamp_IsUTC(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	_ = l.LogSecretRead([]string{"SECRET"})

	var ev audit.Event
	_ = json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &ev)

	if ev.Timestamp.Location().String() != "UTC" {
		t.Errorf("expected UTC timestamp, got %v", ev.Timestamp.Location())
	}
}
