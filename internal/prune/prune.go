// Package prune removes chains whose keys are all empty-valued,
// or removes individual keys that have an empty value from a chain.
package prune

import (
	"fmt"

	"github.com/envchain-cli/envchain-cli/internal/chain"
)

// Result summarises what was removed during a prune operation.
type Result struct {
	// RemovedChains holds the names of chains that were entirely deleted.
	RemovedChains []string
	// RemovedKeys maps chain name → list of keys that were deleted.
	RemovedKeys map[string][]string
}

// Pruner removes empty keys and optionally empty chains from a Store.
type Pruner struct {
	store *chain.Store
}

// NewPruner returns a Pruner backed by store.
func NewPruner(store *chain.Store) *Pruner {
	return &Pruner{store: store}
}

// Prune iterates over the given chain names (or all chains when names is
// empty) and deletes any key whose value is the empty string.  When
// removeEmptyChains is true, a chain that ends up with zero keys is also
// deleted entirely.
func (p *Pruner) Prune(names []string, removeEmptyChains bool) (Result, error) {
	result := Result{
		RemovedKeys: make(map[string][]string),
	}

	targets := names
	if len(targets) == 0 {
		all, err := p.store.List()
		if err != nil {
			return result, fmt.Errorf("prune: list chains: %w", err)
		}
		targets = all
	}

	for _, name := range targets {
		vars, err := p.store.Get(name)
		if err != nil {
			return result, fmt.Errorf("prune: get chain %q: %w", name, err)
		}

		for k, v := range vars {
			if v != "" {
				continue
			}
			if err := p.store.DeleteVar(name, k); err != nil {
				return result, fmt.Errorf("prune: delete key %q from %q: %w", k, name, err)
			}
			result.RemovedKeys[name] = append(result.RemovedKeys[name], k)
		}

		if !removeEmptyChains {
			continue
		}

		remaining, err := p.store.Get(name)
		if err != nil {
			return result, fmt.Errorf("prune: re-read chain %q: %w", name, err)
		}
		if len(remaining) == 0 {
			if err := p.store.DeleteChain(name); err != nil {
				return result, fmt.Errorf("prune: delete chain %q: %w", name, err)
			}
			result.RemovedChains = append(result.RemovedChains, name)
		}
	}

	return result, nil
}
