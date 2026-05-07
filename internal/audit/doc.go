// Package audit implements structured, append-only audit logging for the
// envdiff CLI. Every significant operation — diffing environment files,
// generating reconciliation patches, and accessing masked secrets — is
// recorded as a newline-delimited JSON event.
//
// # Usage
//
//	l := audit.New(os.Stderr)          // write to stderr
//	l.LogDiff(files, diffEntries)      // record a diff run
//	l.LogReconcile(files, patchLines)  // record a patch generation
//	l.LogSecretRead(sensitiveKeys)     // record secret access
//
// For persistent on-disk logs use NewFileLogger:
//
//	logger, fw, err := audit.NewFileLogger("/var/log/envdiff/audit.log")
//	if err != nil { ... }
//	defer fw.Close()
//
// All timestamps are recorded in UTC. Each event is a self-contained JSON
// object suitable for ingestion by log-aggregation systems such as Loki,
// Splunk, or CloudWatch Logs.
package audit
