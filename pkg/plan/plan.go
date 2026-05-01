// Package plan builds and executes safe Docker prune plans with guardrails,
// dry-run by default, and audited execution.
package plan

import (
	"github.com/tyutyutyu/dodu/pkg/group"
)

// Mark identifies an object the user wants to delete.
type Mark struct {
	Kind group.Kind
	ID   string // image ID, container ID, volume name, or build cache ID
}

// Item is a planned (or blocked) deletion.
type Item struct {
	Kind       group.Kind
	ID         string
	Name       string
	EstReclaim int64
	Reason     string // human-readable rationale
}

// Plan is the result of Build: items to delete, plus warnings and blocked
// items that will never be deleted.
type Plan struct {
	Items      []Item
	EstReclaim int64
	Warnings   []string
	Blocked    []Item
}
