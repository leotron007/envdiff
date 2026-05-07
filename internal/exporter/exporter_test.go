package exporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/exporter"
)

func sampleEntries() []differ.Entry {
	return []differ.Entry{
		{Key: "APP_NAME", Value: "myapp", Type: differ.Same},
		{Key: "DB_HOST", Value: "localhost", Type: differ.Added},
		{Key: "OLD_KEY", Value: "gone", Type: differ.Removed},
		{Key: "PORT", Value: "8080", Type: differ.Changed},
	}
}

func TestWrite_EnvFormat_AllEntries(t *testing.T) {
	var buf bytes.Buffer
	ex := exporter.New(&buf, exporter.Options{Format: exporter.FormatEnv})
	if err := ex.Write(sampleEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME in output, got:\n%s", out)
	}
	if !strings.Contains(out, "OLD_KEY=gone") {
		t.Errorf("expected OLD_KEY in output, got:\n%s", out)
	}
}

func TestWrite_EnvFormat_OmitRemoved(t *testing.T) {
	var buf bytes.Buffer
	ex := exporter.New(&buf, exporter.Options{Format: exporter.FormatEnv, OmitRemoved: true})
	if err := ex.Write(sampleEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "OLD_KEY") {
		t.Errorf("removed key should be omitted, got:\n%s", out)
	}
}

func TestWrite_JSONFormat_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	ex := exporter.New(&buf, exporter.Options{Format: exporter.FormatJSON})
	if err := ex.Write(sampleEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 4 {
		t.Errorf("expected 4 entries, got %d", len(result))
	}
	if result[0]["key"] != "APP_NAME" {
		t.Errorf("unexpected first key: %s", result[0]["key"])
	}
}

func TestWrite_JSONFormat_OmitRemoved(t *testing.T) {
	var buf bytes.Buffer
	ex := exporter.New(&buf, exporter.Options{Format: exporter.FormatJSON, OmitRemoved: true})
	if err := ex.Write(sampleEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 entries after omitting removed, got %d", len(result))
	}
}

func TestNew_NilWriter_DefaultsToStdout(t *testing.T) {
	ex := exporter.New(nil, exporter.Options{})
	if ex == nil {
		t.Fatal("expected non-nil exporter")
	}
}

func TestWrite_UnknownFormat_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	ex := exporter.New(&buf, exporter.Options{Format: "xml"})
	if err := ex.Write(sampleEntries()); err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestWrite_ValueWithSpaces_IsQuoted(t *testing.T) {
	var buf bytes.Buffer
	ex := exporter.New(&buf, exporter.Options{Format: exporter.FormatEnv})
	entries := []differ.Entry{
		{Key: "MSG", Value: "hello world", Type: differ.Same},
	}
	if err := ex.Write(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"hello world"`) {
		t.Errorf("expected quoted value in output, got: %s", out)
	}
}
