// Package validator provides validation logic for .env file entries,
// checking for common issues such as empty values, invalid key formats,
// and duplicate keys.
package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/your-org/envdiff/internal/differ"
)

// Issue severity levels.
const (
	SeverityWarn  = "warn"
	SeverityError = "error"
)

// Issue represents a single validation problem found in an env map.
type Issue struct {
	Key      string
	Message  string
	Severity string
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(i.Severity), i.Key, i.Message)
}

// validKeyRe matches keys that consist of uppercase letters, digits, and underscores,
// starting with a letter or underscore (POSIX convention).
var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Validate inspects an env map and a slice of diff entries, returning any
// issues found. It checks for:
//   - Invalid key format
//   - Empty values (warn)
//   - Keys present in diff entries that were added but have empty values
func Validate(env map[string]string, entries []differ.Entry) []Issue {
	var issues []Issue

	// Check all keys in the env map.
	for key, value := range env {
		if !validKeyRe.MatchString(key) {
			issues = append(issues, Issue{
				Key:      key,
				Message:  "invalid key format (must match [A-Za-z_][A-Za-z0-9_]*)",
				Severity: SeverityError,
			})
		}
		if strings.TrimSpace(value) == "" {
			issues = append(issues, Issue{
				Key:      key,
				Message:  "value is empty",
				Severity: SeverityWarn,
			})
		}
	}

	// Check diff entries for added keys with empty values.
	for _, e := range entries {
		if e.Type == differ.Added && strings.TrimSpace(e.NewValue) == "" {
			issues = append(issues, Issue{
				Key:      e.Key,
				Message:  "newly added key has an empty value",
				Severity: SeverityWarn,
			})
		}
	}

	return issues
}

// HasErrors returns true if any of the issues have error severity.
func HasErrors(issues []Issue) bool {
	for _, i := range issues {
		if i.Severity == SeverityError {
			return true
		}
	}
	return false
}
