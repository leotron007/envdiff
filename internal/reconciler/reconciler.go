// Package reconciler provides functionality to generate reconciliation
// patches between two .env files, producing the minimal set of changes
// needed to bring a source environment in sync with a target environment.
package reconciler

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Action represents a reconciliation action type.
type Action string

const (
	ActionAdd    Action = "add"
	ActionRemove Action = "remove"
	ActionUpdate Action = "update"
)

// Patch represents a single reconciliation operation.
type Patch struct {
	Key    string
	Action Action
	Value  string // empty for remove actions
}

// Result holds the full reconciliation output.
type Result struct {
	Patches []Patch
}

// Generate produces a Result containing all patches needed to reconcile
// the source env map to match the target env map.
func Generate(entries []differ.DiffEntry) Result {
	var patches []Patch

	for _, e := range entries {
		switch e.Type {
		case differ.Added:
			patches = append(patches, Patch{Key: e.Key, Action: ActionAdd, Value: e.NewValue})
		case differ.Removed:
			patches = append(patches, Patch{Key: e.Key, Action: ActionRemove})
		case differ.Changed:
			patches = append(patches, Patch{Key: e.Key, Action: ActionUpdate, Value: e.NewValue})
		}
	}

	sort.Slice(patches, func(i, j int) bool {
		return patches[i].Key < patches[j].Key
	})

	return Result{Patches: patches}
}

// RenderPatch formats a Result as lines suitable for appending to or
// replacing values in a .env file. Remove actions are rendered as comments.
func RenderPatch(r Result) string {
	var sb strings.Builder
	for _, p := range r.Patches {
		switch p.Action {
		case ActionAdd, ActionUpdate:
			fmt.Fprintf(&sb, "%s=%s\n", p.Key, p.Value)
		case ActionRemove:
			fmt.Fprintf(&sb, "# REMOVE: %s\n", p.Key)
		}
	}
	return sb.String()
}
