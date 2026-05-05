package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/config"
)

func writeTempConfig(t *testing.T, v any) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	p := filepath.Join(t.TempDir(), "envdiff.json")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}

func TestDefault_HasSensibleValues(t *testing.T) {
	cfg := config.Default()
	if cfg.OutputFormat != config.FormatText {
		t.Errorf("expected FormatText, got %q", cfg.OutputFormat)
	}
	if !cfg.MaskSecrets {
		t.Error("expected MaskSecrets to be true by default")
	}
	if cfg.FailOnDiff {
		t.Error("expected FailOnDiff to be false by default")
	}
}

func TestLoadFile_MergesFields(t *testing.T) {
	raw := map[string]any{
		"files":         []string{"a.env", "b.env"},
		"output_format": "json",
		"fail_on_diff":  true,
	}
	p := writeTempConfig(t, raw)

	cfg := config.Default()
	if err := config.LoadFile(p, cfg); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}

	if cfg.OutputFormat != config.FormatJSON {
		t.Errorf("expected FormatJSON, got %q", cfg.OutputFormat)
	}
	if !cfg.FailOnDiff {
		t.Error("expected FailOnDiff to be true after load")
	}
	if len(cfg.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(cfg.Files))
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	cfg := config.Default()
	err := config.LoadFile("/nonexistent/path.json", cfg)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidate_TooFewFiles(t *testing.T) {
	cfg := config.Default()
	cfg.Files = []string{"only-one.env"}
	if err := config.Validate(cfg); err == nil {
		t.Error("expected validation error for fewer than 2 files")
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	cfg := config.Default()
	cfg.Files = []string{"a.env", "b.env"}
	cfg.OutputFormat = "xml"
	if err := config.Validate(cfg); err == nil {
		t.Error("expected validation error for unknown format")
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := config.Default()
	cfg.Files = []string{"a.env", "b.env"}
	if err := config.Validate(cfg); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}
