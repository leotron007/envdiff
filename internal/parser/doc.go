// Package parser provides functionality for reading and parsing .env files.
//
// It supports the common .env format:
//   - KEY=VALUE pairs (one per line)
//   - Lines beginning with '#' are treated as comments and skipped
//   - Inline comments are stripped (e.g. KEY=value # comment)
//   - Values wrapped in single or double quotes are unquoted automatically
//   - Empty lines are ignored
//
// Basic usage:
//
//	ef, err := parser.LoadFile(".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	val, ok := ef.Get("DATABASE_URL")
//	if ok {
//		fmt.Println("DB URL:", val)
//	}
//
// The parser is intentionally strict: any line that is not a comment,
// blank, or a valid KEY=VALUE pair will return a parse error with the
// offending line number included in the message.
package parser
