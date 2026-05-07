// Package audit provides structured audit logging for envdiff operations,
// recording diffs, reconciliations, and secret access events to a persistent log.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/envdiff/envdiff/internal/differ"
)

// EventType classifies the kind of audit event.
type EventType string

const (
	EventDiff       EventType = "diff"
	EventReconcile  EventType = "reconcile"
	EventSecretRead EventType = "secret_read"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Type      EventType `json:"type"`
	User      string    `json:"user,omitempty"`
	Files     []string  `json:"files,omitempty"`
	Summary   string    `json:"summary"`
	Entries   int       `json:"entries,omitempty"`
}

// Logger writes audit events as newline-delimited JSON.
type Logger struct {
	w io.Writer
}

// New creates a Logger writing to w. If w is nil, os.Stderr is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{w: w}
}

// LogDiff records a diff event given the source files and resulting entries.
func (l *Logger) LogDiff(files []string, entries []differ.Entry) error {
	s := differ.Summary(entries)
	return l.write(Event{
		Timestamp: time.Now().UTC(),
		Type:      EventDiff,
		Files:     files,
		Summary:   fmt.Sprintf("+%d -%d ~%d =%d", s.Added, s.Removed, s.Changed, s.Same),
		Entries:   len(entries),
	})
}

// LogReconcile records a reconcile/patch-generation event.
func (l *Logger) LogReconcile(files []string, patchLines int) error {
	return l.write(Event{
		Timestamp: time.Now().UTC(),
		Type:      EventReconcile,
		Files:     files,
		Summary:   fmt.Sprintf("%d patch lines generated", patchLines),
		Entries:   patchLines,
	})
}

// LogSecretRead records that sensitive keys were accessed/masked.
func (l *Logger) LogSecretRead(keys []string) error {
	return l.write(Event{
		Timestamp: time.Now().UTC(),
		Type:      EventSecretRead,
		Summary:   fmt.Sprintf("%d sensitive key(s) accessed", len(keys)),
		Entries:   len(keys),
	})
}

func (l *Logger) write(e Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}
