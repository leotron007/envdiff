// Package reconciler generates and renders reconciliation patches between
// two .env environments.
//
// Given a slice of [differ.DiffEntry] values produced by the differ package,
// the reconciler computes the minimal set of operations (add, update, remove)
// required to bring a source environment into alignment with a target.
//
// # Usage
//
//	entries := differ.Diff(source, target)
//	result := reconciler.Generate(entries)
//	fmt.Print(reconciler.RenderPatch(result))
//
// The rendered patch is a plain-text representation suitable for manual
// application or further processing in CI pipelines. Remove operations are
// emitted as comments so that they are visible but not automatically applied.
package reconciler
