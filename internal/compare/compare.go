package compare

import (
	"fmt"
	"sort"

	"github.com/envchain-cli/internal/chain"
)

// Result holds the comparison outcome between two chains.
type Result struct {
	OnlyInA    []string          // keys present only in chain A
	OnlyInB    []string          // keys present only in chain B
	SameValue  []string          // keys with identical values in both
	DiffValue  map[string][2]string // keys with differing values: [A, B]
}

// Comparer compares variable sets across two chains.
type Comparer struct {
	store *chain.Store
}

// NewComparer returns a Comparer backed by the given store.
func NewComparer(s *chain.Store) *Comparer {
	return &Comparer{store: s}
}

// Compare inspects two chains and returns a Result describing their differences.
func (c *Comparer) Compare(chainA, chainB string) (*Result, error) {
	vsA, err := c.store.Get(chainA)
	if err != nil {
		return nil, fmt.Errorf("compare: chain %q not found: %w", chainA, err)
	}
	vsB, err := c.store.Get(chainB)
	if err != nil {
		return nil, fmt.Errorf("compare: chain %q not found: %w", chainB, err)
	}

	res := &Result{
		DiffValue: make(map[string][2]string),
	}

	setA := make(map[string]string, len(vsA))
	for _, kv := range vsA {
		setA[kv.Key] = kv.Value
	}
	setB := make(map[string]string, len(vsB))
	for _, kv := range vsB {
		setB[kv.Key] = kv.Value
	}

	allKeys := unionKeys(setA, setB)
	for _, k := range allKeys {
		va, inA := setA[k]
		vb, inB := setB[k]
		switch {
		case inA && !inB:
			res.OnlyInA = append(res.OnlyInA, k)
		case !inA && inB:
			res.OnlyInB = append(res.OnlyInB, k)
		case va == vb:
			res.SameValue = append(res.SameValue, k)
		default:
			res.DiffValue[k] = [2]string{va, vb}
		}
	}
	return res, nil
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
