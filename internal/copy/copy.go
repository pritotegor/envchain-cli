// Package copy provides functionality to duplicate environment variable
// chains within the envchain store, optionally overwriting an existing target.
package copy

import (
	"fmt"

	"github.com/envchain-cli/envchain/internal/chain"
)

// Copier duplicates chains within a Store.
type Copier struct {
	store *chain.Store
}

// NewCopier returns a Copier backed by the given Store.
func NewCopier(store *chain.Store) *Copier {
	return &Copier{store: store}
}

// Copy duplicates the chain named src into a new chain named dst.
// If overwrite is false and dst already exists, Copy returns an error.
// If overwrite is true, any existing dst chain is replaced.
func (c *Copier) Copy(src, dst string, overwrite bool) error {
	vars, err := c.store.Get(src)
	if err != nil {
		return fmt.Errorf("copy: source chain %q not found: %w", src, err)
	}

	if !overwrite {
		if _, err := c.store.Get(dst); err == nil {
			return fmt.Errorf("copy: destination chain %q already exists (use --overwrite to replace)", dst)
		}
	}

	// Build a fresh map so mutations to the copy don't affect the source.
	copied := make(map[string]string, len(vars))
	for k, v := range vars {
		copied[k] = v
	}

	if err := c.store.Add(dst, copied); err != nil {
		return fmt.Errorf("copy: failed to write destination chain %q: %w", dst, err)
	}

	return nil
}

// Rename copies src to dst and then removes src.
// If overwrite is false and dst already exists, Rename returns an error
// without modifying either chain.
func (c *Copier) Rename(src, dst string, overwrite bool) error {
	if err := c.Copy(src, dst, overwrite); err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	if err := c.store.Delete(src); err != nil {
		return fmt.Errorf("rename: could not remove source chain %q after copy: %w", src, err)
	}
	return nil
}
