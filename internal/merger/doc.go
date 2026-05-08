// Package merger implements multi-source .env file merging for envdiff.
//
// It accepts an ordered slice of env maps (typically produced by the parser
// package) and combines them into a single resolved map.  When the same key
// appears in more than one source, the caller chooses a conflict-resolution
// Strategy:
//
//   - StrategyLastWins  – the value from the last source wins (default).
//   - StrategyFirstWins – the value from the first source is preserved.
//   - StrategyError     – an error is returned immediately on the first conflict.
//
// Regardless of strategy, all observed conflicts are recorded in Result.Conflicts
// so that callers can surface warnings or audit information.
//
// Typical usage:
//
//	res, err := merger.Merge([]map[string]string{base, override}, merger.StrategyLastWins)
//	if err != nil { ... }
//	fmt.Println(res.Merged, res.Conflicts)
package merger
