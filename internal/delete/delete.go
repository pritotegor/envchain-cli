// Package delete provides functionality to remove an environment chain
// or individual variables from a chain within the envchain store.
package delete

import (
	"errors"
	"fmt"
)

// ErrChainNotFound is returned when the target chain does not exist.
var ErrChainNotFound = errors.New("chain not found")

// Store is the interface required by the Deleter to read and write chains.
type Store interface {
	Get(chain string) (map[string]string, error)
	Delete(chain string) error
	Set(chain string, vars map[string]string) error
}

// Deleter handles deletion of chains and individual variables.
type Deleter struct {
	store Store
}

// NewDeleter returns a Deleter backed by the given Store.
func NewDeleter(s Store) *Deleter {
	return &Deleter{store: s}
}

// DeleteChain removes an entire chain from the store.
// Returns ErrChainNotFound if the chain does not exist.
func (d *Deleter) DeleteChain(chain string) error {
	_, err := d.store.Get(chain)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrChainNotFound, chain)
	}
	return d.store.Delete(chain)
}

// DeleteVar removes a single variable from a chain.
// Returns ErrChainNotFound if the chain does not exist.
// Returns an error if the key is not present in the chain.
func (d *Deleter) DeleteVar(chain, key string) error {
	vars, err := d.store.Get(chain)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrChainNotFound, chain)
	}
	if _, ok := vars[key]; !ok {
		return fmt.Errorf("variable %q not found in chain %q", key, chain)
	}
	delete(vars, key)
	return d.store.Set(chain, vars)
}
