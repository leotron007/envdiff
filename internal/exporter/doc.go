// Package exporter provides functionality to serialise diff entries into
// various output formats suitable for writing merged or reconciled .env files.
//
// Supported formats:
//   - env / dotenv: standard KEY=VALUE lines, values with whitespace are quoted
//   - json: a JSON array of objects with key, value, and type fields
//   - yaml: a YAML mapping of key/value pairs
//
// The OmitRemoved option skips entries that were deleted in the diff, which is
// useful when generating a clean merged output rather than a full audit trail.
//
// Usage:
//
//	ex := exporter.New(os.Stdout, exporter.Options{
//	    Format:      exporter.FormatEnv,
//	    OmitRemoved: true,
//	})
//	if err := ex.Write(entries); err != nil {
//	    log.Fatal(err)
//	}
package exporter
