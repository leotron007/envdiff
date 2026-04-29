package ci_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/ci"
	"github.com/user/envdiff/internal/differ"
)

func newDetector(env map[string]string) *ci.Detector {
	d := &ci.Detector{}
	// Use exported constructor but override via test helper below.
	_ = d
	// Re-implement via a local wrapper since Detector.env is unexported.
	// We test Detect indirectly through a thin exported API.
	return ci.NewDetectorWithEnv(func(key string) string {
		return env[key]
	})
}

func TestDetect_GitHub(t *testing.T) {
	d := newDetector(map[string]string{"GITHUB_ACTIONS": "true"})
	if got := d.Detect(); got != ci.PlatformGitHub {
		t.Errorf("expected github, got %q", got)
	}
}

func TestDetect_GitLab(t *testing.T) {
	d := newDetector(map[string]string{"GITLAB_CI": "true"})
	if got := d.Detect(); got != ci.PlatformGitLab {
		t.Errorf("expected gitlab, got %q", got)
	}
}

func TestDetect_CircleCI(t *testing.T) {
	d := newDetector(map[string]string{"CIRCLECI": "true"})
	if got := d.Detect(); got != ci.PlatformCircleCI {
		t.Errorf("expected circleci, got %q", got)
	}
}

func TestDetect_Unknown(t *testing.T) {
	d := newDetector(map[string]string{})
	if got := d.Detect(); got != ci.PlatformUnknown {
		t.Errorf("expected unknown, got %q", got)
	}
}

func TestIsCI_True(t *testing.T) {
	d := newDetector(map[string]string{"GITHUB_ACTIONS": "true"})
	if !d.IsCI() {
		t.Error("expected IsCI to return true")
	}
}

func TestIsCI_False(t *testing.T) {
	d := newDetector(map[string]string{})
	if d.IsCI() {
		t.Error("expected IsCI to return false")
	}
}

func TestAnnotate_GitHub_Warning(t *testing.T) {
	var buf bytes.Buffer
	a := ci.NewAnnotator(ci.PlatformGitHub, &buf)
	entries := []differ.DiffEntry{
		{Key: "FOO", Type: differ.Added},
	}
	a.Annotate(entries)
	if !strings.Contains(buf.String(), "::warning") {
		t.Errorf("expected GitHub warning annotation, got: %q", buf.String())
	}
}

func TestAnnotate_GitHub_Error_OnRemoved(t *testing.T) {
	var buf bytes.Buffer
	a := ci.NewAnnotator(ci.PlatformGitHub, &buf)
	entries := []differ.DiffEntry{
		{Key: "BAR", Type: differ.Removed},
	}
	a.Annotate(entries)
	if !strings.Contains(buf.String(), "::error") {
		t.Errorf("expected GitHub error annotation, got: %q", buf.String())
	}
}

func TestAnnotate_SkipsSame(t *testing.T) {
	var buf bytes.Buffer
	a := ci.NewAnnotator(ci.PlatformGitHub, &buf)
	entries := []differ.DiffEntry{
		{Key: "SAME_KEY", Type: differ.Same},
	}
	a.Annotate(entries)
	if buf.Len() != 0 {
		t.Errorf("expected no output for Same entries, got: %q", buf.String())
	}
}

func TestAnnotate_DefaultPlatform(t *testing.T) {
	var buf bytes.Buffer
	a := ci.NewAnnotator(ci.PlatformUnknown, &buf)
	entries := []differ.DiffEntry{
		{Key: "X", Type: differ.Changed},
	}
	a.Annotate(entries)
	if !strings.Contains(buf.String(), "[envdiff]") {
		t.Errorf("expected generic annotation, got: %q", buf.String())
	}
}
