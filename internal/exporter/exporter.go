package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Format represents the output format for exported env files.
type Format string

const (
	FormatEnv    Format = "env"
	FormatJSON   Format = "json"
	FormatDotenv Format = "dotenv"
)

// Options controls export behaviour.
type Options struct {
	Format      Format
	MaskSecrets bool
	OmitRemoved bool
}

// Exporter writes reconciled env entries to a destination.
type Exporter struct {
	w   io.Writer
	opt Options
}

// New creates an Exporter writing to w. If w is nil, os.Stdout is used.
func New(w io.Writer, opt Options) *Exporter {
	if w == nil {
		w = os.Stdout
	}
	if opt.Format == "" {
		opt.Format = FormatEnv
	}
	return &Exporter{w: w, opt: opt}
}

// Write serialises entries into the chosen format.
func (e *Exporter) Write(entries []differ.Entry) error {
	switch e.opt.Format {
	case FormatJSON:
		return e.writeJSON(entries)
	case FormatEnv, FormatDotenv:
		return e.writeEnv(entries)
	default:
		return fmt.Errorf("exporter: unknown format %q", e.opt.Format)
	}
}

func (e *Exporter) writeEnv(entries []differ.Entry) error {
	for _, en := range entries {
		if e.opt.OmitRemoved && en.Type == differ.Removed {
			continue
		}
		val := en.Value
		if needsQuoting(val) {
			val = fmt.Sprintf("%q", val)
		}
		if _, err := fmt.Fprintf(e.w, "%s=%s\n", en.Key, val); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeJSON(entries []differ.Entry) error {
	type jsonEntry struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Type  string `json:"type"`
	}
	out := make([]jsonEntry, 0, len(entries))
	for _, en := range entries {
		if e.opt.OmitRemoved && en.Type == differ.Removed {
			continue
		}
		out = append(out, jsonEntry{Key: en.Key, Value: en.Value, Type: string(en.Type)})
	}
	enc := json.NewEncoder(e.w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func needsQuoting(s string) bool {
	return strings.ContainsAny(s, " \t\n#")
}
