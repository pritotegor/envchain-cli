// Package clone provides functionality to deep-clone an entire chain store
// into a new named chain, optionally filtering keys by a prefix.
package clone

import (
	"fmt"

	"github.com/envchain-cli/envchain-cli/internal/chain"
)

// Cloner clones one chain into a new chain.
type Cloner struct {
	store *chain.Store
}

// NewCloner returns a Cloner backed by the given store.
func NewCloner(store *chain.Store) *Cloner {
	return &Cloner{store: store}
}

// CloneOptions controls the behaviour of a clone operation.
type CloneOptions struct {
	// Overwrite allows the destination chain to be replaced if it already exists.
	Overwrite bool
	// Prefix, when non-empty, restricts cloned keys to those that start with Prefix.
	Prefix string
}

// Clone copies all matching keys from src into dst.
// It returns an error if src does not exist, or if dst already exists and
// Overwrite is false.
func (c *Cloner) Clone(src, dst string, opts CloneOptions) error {
	vars, err := c.store.Get(src)
	if err != nil {
		return fmt.Errorf("clone: source chain %q not found: %w", src, err)
	}

	if !opts.Overwrite {
		if _, err := c.store.Get(dst); err == nil {
			return fmt.Errorf("clone: destination chain %q already exists (use --overwrite to replace)", dst)
		}
	}

	filtered := make(map[string]string, len(vars))
	for k, v := range vars {
		if opts.Prefix == "" || len(k) >= len(opts.Prefix) && k[:len(opts.Prefix)] == opts.Prefix {
			filtered[k] = v
		}
	}

	if len(filtered) == 0 {
		return fmt.Errorf("clone: no keys matched prefix %q in chain %q", opts.Prefix, src)
	}

	return c.store.Set(dst, filtered)
}
