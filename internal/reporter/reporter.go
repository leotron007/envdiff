// Package reporter formats and outputs diff results for human and CI consumption.
package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Format represents the output format for the reporter.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes diff output to a writer.
type Reporter struct {
	w      io.Writer
	format Format
}

// New creates a new Reporter writing to w with the given format.
func New(w io.Writer, format Format) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{w: w, format: format}
}

// Write outputs the diff entries to the configured writer.
func (r *Reporter) Write(entries []differ.DiffEntry) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(entries)
	default:
		return r.writeText(entries)
	}
}

func (r *Reporter) writeText(entries []differ.DiffEntry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(r.w, "No differences found.")
		return err
	}
	for _, e := range entries {
		var line string
		switch e.Type {
		case differ.Added:
			line = fmt.Sprintf("+ %s=%s", e.Key, e.NewValue)
		case differ.Removed:
			line = fmt.Sprintf("- %s=%s", e.Key, e.OldValue)
		case differ.Changed:
			line = fmt.Sprintf("~ %s: %s -> %s", e.Key, e.OldValue, e.NewValue)
		case differ.Same:
			line = fmt.Sprintf("  %s=%s", e.Key, e.NewValue)
		}
		if _, err := fmt.Fprintln(r.w, line); err != nil {
			return err
		}
	}
	return nil
}

func (r *Reporter) writeJSON(entries []differ.DiffEntry) error {
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, e := range entries {
		sb.WriteString(fmt.Sprintf(
			"  {\"key\": %q, \"type\": %q, \"old\": %q, \"new\": %q}",
			e.Key, string(e.Type), e.OldValue, e.NewValue,
		))
		if i < len(entries)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")
	_, err := fmt.Fprint(r.w, sb.String())
	return err
}
