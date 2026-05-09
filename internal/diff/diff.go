package diff

import (
	"errors"
	"fmt"
	"sort"
)

// ChangeKind describes how a key differs between two chains.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
	Unchanged ChangeKind = "unchanged"
)

// Entry represents a single key comparison result.
type Entry struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// String returns a human-readable representation of the entry.
func (e Entry) String() string {
	switch e.Kind {
	case Added:
		return fmt.Sprintf("+ %s=%s", e.Key, e.NewValue)
	case Removed:
		return fmt.Sprintf("- %s=%s", e.Key, e.OldValue)
	case Modified:
		return fmt.Sprintf("~ %s: %s → %s", e.Key, e.OldValue, e.NewValue)
	default:
		return fmt.Sprintf("  %s=%s", e.Key, e.NewValue)
	}
}

// Colored returns an ANSI-colored version of String.
func (e Entry) Colored() string {
	const (
		green  = "\033[32m"
		red    = "\033[31m"
		yellow = "\033[33m"
		reset  = "\033[0m"
	)
	switch e.Kind {
	case Added:
		return green + e.String() + reset
	case Removed:
		return red + e.String() + reset
	case Modified:
		return yellow + e.String() + reset
	default:
		return e.String()
	}
}

// Result holds all entries from a diff operation.
type Result struct {
	entries []Entry
}

// Changes returns only entries that are not Unchanged.
func (r Result) Changes() []Entry {
	var out []Entry
	for _, e := range r.entries {
		if e.Kind != Unchanged {
			out = append(out, e)
		}
	}
	return out
}

// All returns every entry including unchanged ones.
func (r Result) All() []Entry { return r.entries }

// Store is the minimal interface required by Differ.
type Store interface {
	Get(chain, key string) (string, error)
	List(chain string) (map[string]string, error)
}

// Differ computes differences between two chains.
type Differ struct {
	store Store
}

// NewDiffer creates a Differ backed by the given store.
func NewDiffer(s Store) *Differ { return &Differ{store: s} }

// Diff compares src and dst chains and returns the Result.
func (d *Differ) Diff(src, dst string) (Result, error) {
	srcVars, err := d.store.List(src)
	if err != nil {
		return Result{}, fmt.Errorf("diff: source chain %q: %w", src, err)
	}
	dstVars, err := d.store.List(dst)
	if err != nil {
		return Result{}, fmt.Errorf("diff: destination chain %q: %w", dst, err)
	}

	keys := make(map[string]struct{})
	for k := range srcVars {
		keys[k] = struct{}{}
	}
	for k := range dstVars {
		keys[k] = struct{}{}
	}

	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	var entries []Entry
	for _, k := range sorted {
		sv, inSrc := srcVars[k]
		dv, inDst := dstVars[k]
		switch {
		case inSrc && !inDst:
			entries = append(entries, Entry{Key: k, Kind: Removed, OldValue: sv})
		case !inSrc && inDst:
			entries = append(entries, Entry{Key: k, Kind: Added, NewValue: dv})
		case sv != dv:
			entries = append(entries, Entry{Key: k, Kind: Modified, OldValue: sv, NewValue: dv})
		default:
			entries = append(entries, Entry{Key: k, Kind: Unchanged, OldValue: sv, NewValue: dv})
		}
	}

	_ = errors.New // keep import tidy
	return Result{entries: entries}, nil
}
