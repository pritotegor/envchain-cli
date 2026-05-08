// Package rename provides functionality to rename an existing chain within a store.
package rename

import (
	"errors"
	"fmt"

	"github.com/yourorg/envchain-cli/internal/chain"
)

// ErrSourceNotFound is returned when the source chain does not exist.
var ErrSourceNotFound = errors.New("source chain not found")

// ErrDestinationExists is returned when the destination chain already exists
// and overwrite is not requested.
var ErrDestinationExists = errors.New("destination chain already exists")

// Renamer renames chains within a Store.
type Renamer struct {
	store *chain.Store
}

// NewRenamer creates a new Renamer backed by the given Store.
func NewRenamer(store *chain.Store) *Renamer {
	return &Renamer{store: store}
}

// Rename renames the chain identified by src to dst.
// If overwrite is false and dst already exists, ErrDestinationExists is returned.
// The src chain is removed after its variables are copied to dst.
func (r *Renamer) Rename(src, dst string, overwrite bool) error {
	vars, err := r.store.Get(src)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrSourceNotFound, src)
	}

	if !overwrite {
		if _, err := r.store.Get(dst); err == nil {
			return fmt.Errorf("%w: %s", ErrDestinationExists, dst)
		}
	}

	if err := r.store.Add(dst, vars); err != nil {
		return fmt.Errorf("rename: writing destination chain %q: %w", dst, err)
	}

	if err := r.store.Delete(src); err != nil {
		return fmt.Errorf("rename: removing source chain %q after copy: %w", src, err)
	}

	return nil
}
