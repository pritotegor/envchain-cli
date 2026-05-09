package merge

import (
	"fmt"

	"github.com/envchain-cli/internal/chain"
)

// Merger merges variables from one or more source chains into a destination chain.
type Merger struct {
	store *chain.Store
}

// NewMerger creates a new Merger backed by the given store.
func NewMerger(store *chain.Store) *Merger {
	return &Merger{store: store}
}

// MergeOptions controls the behaviour of a merge operation.
type MergeOptions struct {
	// Overwrite allows destination variables to be overwritten by source values.
	Overwrite bool
}

// Merge copies all variables from each source chain into the destination chain.
// Sources are processed left-to-right; later sources win when Overwrite is true.
// If Overwrite is false and a key already exists in the destination the source
// value is skipped (no error is returned).
func (m *Merger) Merge(dst string, sources []string, opts MergeOptions) error {
	if _, err := m.store.Get(dst); err != nil {
		// destination does not exist yet – create an empty chain
		if err2 := m.store.Add(dst, "", ""); err2 != nil {
			return fmt.Errorf("merge: create destination chain %q: %w", dst, err2)
		}
	}

	for _, src := range sources {
		vars, err := m.store.Get(src)
		if err != nil {
			return fmt.Errorf("merge: source chain %q not found", src)
		}

		dstVars, err := m.store.Get(dst)
		if err != nil {
			return fmt.Errorf("merge: read destination chain %q: %w", dst, err)
		}

		for k, v := range vars {
			if _, exists := dstVars[k]; exists && !opts.Overwrite {
				continue
			}
			if err := m.store.Add(dst, k, v); err != nil {
				return fmt.Errorf("merge: set %q in %q: %w", k, dst, err)
			}
		}
	}

	return nil
}
