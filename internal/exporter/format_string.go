package exporter

import "fmt"

// ParseFormat converts a raw string into a Format constant.
// Returns an error if the value is not recognised.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatEnv, FormatDotenv, FormatJSON:
		return Format(s), nil
	default:
		return "", fmt.Errorf("exporter: unrecognised format %q; valid options are env, dotenv, json", s)
	}
}

// String implements fmt.Stringer.
func (f Format) String() string {
	return string(f)
}

// IsText reports whether the format produces plain-text line-oriented output.
func (f Format) IsText() bool {
	return f == FormatEnv || f == FormatDotenv
}

// IsStructured reports whether the format produces structured (machine-readable) output.
func (f Format) IsStructured() bool {
	return f == FormatJSON
}
