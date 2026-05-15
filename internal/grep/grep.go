// Package grep provides functionality for searching environment variable
// keys and values across one or more chains in the store.
package grep

import (
	"regexp"
	"sort"

	"github.com/envchain-cli/envchain-cli/internal/chain"
)

// Match represents a single search result.
type Match struct {
	Chain string
	Key   string
	Value string
}

// Grepper searches chains for keys or values matching a pattern.
type Grepper struct {
	store *chain.Store
}

// NewGrepper returns a Grepper backed by the given store.
func NewGrepper(s *chain.Store) *Grepper {
	return &Grepper{store: s}
}

// Search scans the given chains (or all chains if none specified) for keys or
// values matching pattern. If keysOnly is true only key names are matched.
func (g *Grepper) Search(pattern string, chains []string, keysOnly bool) ([]Match, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	if len(chains) == 0 {
		chains, err = g.store.List()
		if err != nil {
			return nil, err
		}
	}

	var results []Match
	for _, name := range chains {
		vars, err := g.store.Get(name)
		if err != nil {
			continue
		}
		keys := make([]string, 0, len(vars))
		for k := range vars {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := vars[k]
			if re.MatchString(k) || (!keysOnly && re.MatchString(v)) {
				results = append(results, Match{Chain: name, Key: k, Value: v})
			}
		}
	}
	return results, nil
}
