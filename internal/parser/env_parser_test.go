package parser

import (
	"strings"
	"testing"
)

func TestParse_BasicKeyValue(t *testing.T) {
	input := `APP_NAME=myapp
DEBUG=true
PORT=8080
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(ef.Entries))
	}
	if ef.Index["APP_NAME"].Value != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %q", ef.Index["APP_NAME"].Value)
	}
}

func TestParse_SkipsComments(t *testing.T) {
	input := `# This is a comment
KEY=value
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(ef.Entries))
	}
}

func TestParse_StripQuotes(t *testing.T) {
	input := `SECRET="my secret value"
TOKEN='another token'
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Index["SECRET"].Value != "my secret value" {
		t.Errorf("expected unquoted value, got %q", ef.Index["SECRET"].Value)
	}
	if ef.Index["TOKEN"].Value != "another token" {
		t.Errorf("expected unquoted value, got %q", ef.Index["TOKEN"].Value)
	}
}

func TestParse_InlineComment(t *testing.T) {
	input := `HOST=localhost # the host
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := ef.Index["HOST"]
	if entry.Value != "localhost" {
		t.Errorf("expected value 'localhost', got %q", entry.Value)
	}
	if entry.Comment != "# the host" {
		t.Errorf("expected comment '# the host', got %q", entry.Comment)
	}
}

func TestParse_InvalidLine(t *testing.T) {
	input := `INVALID_LINE
`
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParse_EmptyFile(t *testing.T) {
	ef, err := Parse(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(ef.Entries))
	}
}
