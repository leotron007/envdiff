package validator_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/differ"
	"github.com/your-org/envdiff/internal/validator"
)

func TestValidate_ValidEnv_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	issues := validator.Validate(env, nil)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestValidate_InvalidKeyFormat_ReturnsError(t *testing.T) {
	env := map[string]string{
		"123INVALID": "value",
	}
	issues := validator.Validate(env, nil)
	if len(issues) == 0 {
		t.Fatal("expected at least one issue for invalid key format")
	}
	if issues[0].Severity != validator.SeverityError {
		t.Errorf("expected error severity, got %s", issues[0].Severity)
	}
}

func TestValidate_EmptyValue_ReturnsWarn(t *testing.T) {
	env := map[string]string{
		"EMPTY_KEY": "",
	}
	issues := validator.Validate(env, nil)
	if len(issues) == 0 {
		t.Fatal("expected at least one issue for empty value")
	}
	found := false
	for _, i := range issues {
		if i.Key == "EMPTY_KEY" && i.Severity == validator.SeverityWarn {
			found = true
		}
	}
	if !found {
		t.Error("expected warn issue for EMPTY_KEY")
	}
}

func TestValidate_AddedEntryEmptyValue_ReturnsWarn(t *testing.T) {
	entries := []differ.Entry{
		{Key: "NEW_KEY", Type: differ.Added, NewValue: ""},
	}
	issues := validator.Validate(map[string]string{}, entries)
	if len(issues) == 0 {
		t.Fatal("expected issue for added entry with empty value")
	}
	if issues[0].Severity != validator.SeverityWarn {
		t.Errorf("expected warn, got %s", issues[0].Severity)
	}
}

func TestHasErrors_WithErrorIssue_ReturnsTrue(t *testing.T) {
	issues := []validator.Issue{
		{Key: "BAD", Message: "invalid", Severity: validator.SeverityError},
	}
	if !validator.HasErrors(issues) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_OnlyWarnings_ReturnsFalse(t *testing.T) {
	issues := []validator.Issue{
		{Key: "WARN_KEY", Message: "empty value", Severity: validator.SeverityWarn},
	}
	if validator.HasErrors(issues) {
		t.Error("expected HasErrors to return false for warnings only")
	}
}

func TestIssue_String_Format(t *testing.T) {
	i := validator.Issue{Key: "MY_KEY", Message: "something wrong", Severity: validator.SeverityError}
	s := i.String()
	if s != "[ERROR] MY_KEY: something wrong" {
		t.Errorf("unexpected string format: %s", s)
	}
}
