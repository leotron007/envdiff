// Package ci provides CI/CD integration helpers for envdiff.
// It detects the current CI environment and formats output accordingly.
package ci

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Platform represents a detected CI platform.
type Platform string

const (
	PlatformUnknown  Platform = "unknown"
	PlatformGitHub   Platform = "github"
	PlatformGitLab   Platform = "gitlab"
	PlatformCircleCI Platform = "circleci"
)

// Detector detects the current CI environment.
type Detector struct {
	env func(string) string
}

// NewDetector creates a new Detector using os.Getenv.
func NewDetector() *Detector {
	return &Detector{env: os.Getenv}
}

// Detect returns the detected CI platform.
func (d *Detector) Detect() Platform {
	switch {
	case d.env("GITHUB_ACTIONS") == "true":
		return PlatformGitHub
	case d.env("GITLAB_CI") == "true":
		return PlatformGitLab
	case d.env("CIRCLECI") == "true":
		return PlatformCircleCI
	default:
		return PlatformUnknown
	}
}

// IsCI returns true if running inside any known CI environment.
func (d *Detector) IsCI() bool {
	return d.Detect() != PlatformUnknown
}

// Annotator writes CI-specific annotations for diff entries.
type Annotator struct {
	platform Platform
	out      io.Writer
}

// NewAnnotator creates an Annotator for the given platform and writer.
func NewAnnotator(platform Platform, out io.Writer) *Annotator {
	if out == nil {
		out = os.Stdout
	}
	return &Annotator{platform: platform, out: out}
}

// Annotate writes platform-specific annotations for each diff entry.
func (a *Annotator) Annotate(entries []differ.DiffEntry) {
	for _, e := range entries {
		if e.Type == differ.Same {
			continue
		}
		switch a.platform {
		case PlatformGitHub:
			level := "warning"
			if e.Type == differ.Removed {
				level = "error"
			}
			fmt.Fprintf(a.out, "::%s title=envdiff::%s\n", level, formatMessage(e))
		default:
			fmt.Fprintf(a.out, "[envdiff] %s\n", formatMessage(e))
		}
	}
}

// Summary returns a human-readable summary of the diff entries,
// reporting the count of added, removed, and changed keys.
func Summary(entries []differ.DiffEntry) string {
	var added, removed, changed int
	for _, e := range entries {
		switch e.Type {
		case differ.Added:
			added++
		case differ.Removed:
			removed++
		case differ.Changed:
			changed++
		}
	}
	return fmt.Sprintf("envdiff: %d added, %d removed, %d changed", added, removed, changed)
}

func formatMessage(e differ.DiffEntry) string {
	switch e.Type {
	case differ.Added:
		return fmt.Sprintf("Key %q added", e.Key)
	case differ.Removed:
		return fmt.Sprintf("Key %q removed", e.Key)
	case differ.Changed:
		return fmt.Sprintf("Key %q changed", e.Key)
	default:
		return strings.ToLower(string(e.Type)) + ": " + e.Key
	}
}
