// Package differ implements environment variable diffing between two env maps.
//
// Basic usage:
//
//	base := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
//	other := map[string]string{"DB_HOST": "prod.db", "PORT": "5432", "DEBUG": "false"}
//
//	result := differ.Diff(base, other)
//	for _, entry := range result.Entries {
//		fmt.Println(entry.Key, entry.Type)
//	}
//
// DiffType values:
//   - DiffAdded:   key exists in other but not in base
//   - DiffRemoved: key exists in base but not in other
//   - DiffChanged: key exists in both but values differ
//   - DiffSame:    key exists in both with identical values
//
// Use FilterByType and Summary for targeted inspection of results.
package differ
