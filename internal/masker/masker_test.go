package masker_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/masker"
)

func TestIsSensitive_MatchesKnownPatterns(t *testing.T) {
	m := masker.New(nil, "")

	sensitive := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "STRIPE_SECRET", "PRIVATE_KEY", "DATABASE_URL"}
	for _, key := range sensitive {
		if !m.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_SafeKeys(t *testing.T) {
	m := masker.New(nil, "")

	safe := []string{"APP_ENV", "PORT", "LOG_LEVEL", "FEATURE_FLAG"}
	for _, key := range safe {
		if m.IsSensitive(key) {
			t.Errorf("expected %q to not be sensitive", key)
		}
	}
}

func TestMask_SensitiveValue(t *testing.T) {
	m := masker.New(nil, "***")
	got := m.Mask("DB_PASSWORD", "supersecret")
	if got != "***" {
		t.Errorf("expected *** got %q", got)
	}
}

func TestMask_SafeValue(t *testing.T) {
	m := masker.New(nil, "***")
	got := m.Mask("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected production got %q", got)
	}
}

func TestMaskMap_MasksOnlySensitive(t *testing.T) {
	m := masker.New(nil, "[REDACTED]")
	input := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "s3cr3t",
		"PORT":        "8080",
		"API_KEY":     "key-abc-123",
	}

	result := m.MaskMap(input)

	if result["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be masked")
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should not be masked")
	}
	if result["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if result["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY should be masked")
	}
}

func TestNew_CustomSensitiveKeys(t *testing.T) {
	m := masker.New([]string{"CUSTOM_FIELD"}, "XXX")
	if !m.IsSensitive("MY_CUSTOM_FIELD") {
		t.Error("expected MY_CUSTOM_FIELD to be sensitive with custom keys")
	}
	if m.IsSensitive("API_KEY") {
		t.Error("API_KEY should not be sensitive with custom keys only")
	}
}
