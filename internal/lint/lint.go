// Package lint provides validation checks for chains stored in a Store,
// reporting variables that shadow each other or keys that violate naming
// conventions without modifying any data.
package lint

import (
	"fmt"
	"strings"

	"github.com/envchain-cli/envchain-cli/internal/chain"
	"github.com/envchain-cli/envchain-cli/internal/validate"
)

// Issue describes a single lint finding.
type Issue struct {
	Chain   string
	Key     string
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("%s.%s: %s", i.Chain, i.Key, i.Message)
}

// Linter runs lint checks against a chain store.
type Linter struct {
	store *chain.Store
}

// NewLinter returns a Linter backed by the given Store.
func NewLinter(s *chain.Store) *Linter {
	return &Linter{store: s}
}

// Lint checks every key in every chain (or only the named chains when names
// are provided) and returns all issues found.
func (l *Linter) Lint(names ...string) ([]Issue, error) {
	targets := names
	if len(targets) == 0 {
		all, err := l.store.List()
		if err != nil {
			return nil, fmt.Errorf("lint: list chains: %w", err)
		}
		targets = all
	}

	var issues []Issue
	for _, name := range targets {
		vars, err := l.store.Get(name)
		if err != nil {
			return nil, fmt.Errorf("lint: get chain %q: %w", name, err)
		}
		for k := range vars {
			if err := validate.Key(k); err != nil {
				issues = append(issues, Issue{Chain: name, Key: k, Message: err.Error()})
				continue
			}
			if k != strings.ToUpper(k) {
				issues = append(issues, Issue{
					Chain:   name,
					Key:     k,
					Message: "key is not upper-case",
				})
			}
		}
	}
	return issues, nil
}
