// Package config provides configuration loading and validation for envdiff CLI.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Config holds the runtime configuration for envdiff.
type Config struct {
	// Files is the list of .env file paths to compare.
	Files []string `json:"files"`

	// OutputFormat controls report output (text or json).
	OutputFormat Format `json:"output_format"`

	// MaskSecrets enables secret masking in output.
	MaskSecrets bool `json:"mask_secrets"`

	// Prefix filters entries to only those matching this prefix.
	Prefix string `json:"prefix,omitempty"`

	// ExcludeKeys is a list of keys to omit from the diff.
	ExcludeKeys []string `json:"exclude_keys,omitempty"`

	// FailOnDiff causes a non-zero exit code when differences are found.
	FailOnDiff bool `json:"fail_on_diff"`

	// CIAnnotations enables CI-specific annotation output.
	CIAnnotations bool `json:"ci_annotations"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		OutputFormat: FormatText,
		MaskSecrets:  true,
		FailOnDiff:   false,
		CIAnnotations: false,
	}
}

// LoadFile reads a JSON config file from the given path and merges it into cfg.
func LoadFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("config: read %q: %w", path, err)
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("config: parse %q: %w", path, err)
	}
	return nil
}

// Validate checks the config for logical errors and returns a combined error.
func Validate(cfg *Config) error {
	var errs []string

	if len(cfg.Files) < 2 {
		errs = append(errs, "at least two files must be specified")
	}

	if cfg.OutputFormat != FormatText && cfg.OutputFormat != FormatJSON {
		errs = append(errs, fmt.Sprintf("unknown output format %q; must be \"text\" or \"json\"", cfg.OutputFormat))
	}

	if len(errs) > 0 {
		return fmt.Errorf("config validation: %s", strings.Join(errs, "; "))
	}
	return nil
}
