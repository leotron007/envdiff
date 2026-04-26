// Package masker provides utilities for masking sensitive values
// in .env file entries before display or output.
package masker

import "strings"

// DefaultSensitiveKeys contains common key patterns considered sensitive.
var DefaultSensitiveKeys = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE_KEY",
	"AUTH",
	"CREDENTIAL",
	"DSN",
	"DATABASE_URL",
}

// Masker masks sensitive environment variable values.
type Masker struct {
	sensitiveKeys []string
	maskValue     string
}

// New creates a Masker with the given sensitive key patterns and mask string.
// If sensitiveKeys is nil, DefaultSensitiveKeys is used.
// If maskValue is empty, "***" is used.
func New(sensitiveKeys []string, maskValue string) *Masker {
	if sensitiveKeys == nil {
		sensitiveKeys = DefaultSensitiveKeys
	}
	if maskValue == "" {
		maskValue = "***"
	}
	return &Masker{
		sensitiveKeys: sensitiveKeys,
		maskValue:     maskValue,
	}
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range m.sensitiveKeys {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// Mask returns the masked value if the key is sensitive, otherwise the original value.
func (m *Masker) Mask(key, value string) string {
	if m.IsSensitive(key) {
		return m.maskValue
	}
	return value
}

// MaskMap returns a copy of the map with sensitive values replaced by the mask string.
func (m *Masker) MaskMap(env map[string]string) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = m.Mask(k, v)
	}
	return result
}
