// Package diff provides functionality to compare two environment variable
// chains and report additions, removals, and modifications between them.
package diff

import (
	"fmt"
	"sort"

	"github.com/envchain-cli/internal/chain"
)

// ChangeKind describes the type of change between two chains.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
)

// Change represents a single variable difference between two chains.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Result holds all changes between a source and destination chain.
type Result struct {
	Source      string
	Destination string
	Changes     []Change
}

// HasChanges returns true when at least one difference was found.
func (r *Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Differ compares chains stored in a Store.
type Differ struct {
	store *chain.Store
}

// NewDiffer creates a Differ backed by the given store.
func NewDiffer(s *chain.Store) *Differ {
	return &Differ{store: s}
}

// Compare returns the diff between srcName and dstName chains.
// An error is returned if either chain does not exist.
func (d *Differ) Compare(srcName, dstName string) (*Result, error) {
	srcVars, err := d.store.Get(srcName)
	if err != nil {
		return nil, fmt.Errorf("source chain %q: %w", srcName, err)
	}

	dstVars, err := d.store.Get(dstName)
	if err != nil {
		return nil, fmt.Errorf("destination chain %q: %w", dstName, err)
	}

	result := &Result{Source: srcName, Destination: dstName}

	// Keys present in src but not in dst, or modified.
	for k, sv := range srcVars {
		if dv, ok := dstVars[k]; !ok {
			result.Changes = append(result.Changes, Change{Key: k, Kind: Removed, OldValue: sv})
		} else if sv != dv {
			result.Changes = append(result.Changes, Change{Key: k, Kind: Modified, OldValue: sv, NewValue: dv})
		}
	}

	// Keys present in dst but not in src.
	for k, dv := range dstVars {
		if _, ok := srcVars[k]; !ok {
			result.Changes = append(result.Changes, Change{Key: k, Kind: Added, NewValue: dv})
		}
	}

	sort.Slice(result.Changes, func(i, j int) bool {
		return result.Changes[i].Key < result.Changes[j].Key
	})

	return result, nil
}
