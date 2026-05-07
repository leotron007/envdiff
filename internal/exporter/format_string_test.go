package exporter_test

import (
	"testing"

	"github.com/user/envdiff/internal/exporter"
)

func TestParseFormat_ValidFormats(t *testing.T) {
	cases := []struct {
		input    string
		wantFmt  exporter.Format
	}{
		{"env", exporter.FormatEnv},
		{"dotenv", exporter.FormatDotenv},
		{"json", exporter.FormatJSON},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := exporter.ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantFmt {
				t.Errorf("got %q, want %q", got, tc.wantFmt)
			}
		})
	}
}

func TestParseFormat_InvalidFormat_ReturnsError(t *testing.T) {
	_, err := exporter.ParseFormat("yaml")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestFormat_String(t *testing.T) {
	if exporter.FormatJSON.String() != "json" {
		t.Errorf("unexpected String() value: %s", exporter.FormatJSON.String())
	}
}

func TestFormat_IsText(t *testing.T) {
	if !exporter.FormatEnv.IsText() {
		t.Error("env format should be text")
	}
	if !exporter.FormatDotenv.IsText() {
		t.Error("dotenv format should be text")
	}
	if exporter.FormatJSON.IsText() {
		t.Error("json format should not be text")
	}
}

func TestFormat_IsStructured(t *testing.T) {
	if !exporter.FormatJSON.IsStructured() {
		t.Error("json format should be structured")
	}
	if exporter.FormatEnv.IsStructured() {
		t.Error("env format should not be structured")
	}
}
