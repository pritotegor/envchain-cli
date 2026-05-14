// Package template provides functionality for rendering environment variable
// values that contain template placeholders referencing other variables within
// the same chain or from a base chain.
package template

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/envchain-cli/envchain-cli/internal/chain"
)

// placeholder matches ${VAR_NAME} style references.
var placeholder = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)

// Renderer resolves template placeholders in chain variable values.
type Renderer struct {
	store *chain.Store
}

// NewRenderer returns a Renderer backed by the given store.
func NewRenderer(store *chain.Store) *Renderer {
	return &Renderer{store: store}
}

// Render returns a copy of the variables for chainName with all ${VAR}
// placeholders expanded. References that cannot be resolved are left as-is.
// Circular references are detected and returned as an error.
func (r *Renderer) Render(chainName string) (map[string]string, error) {
	vars, err := r.store.Get(chainName)
	if err != nil {
		return nil, fmt.Errorf("template: chain %q not found: %w", chainName, err)
	}

	resolved := make(map[string]string, len(vars))
	for k, v := range vars {
		expanded, err := expand(v, vars, []string{k})
		if err != nil {
			return nil, fmt.Errorf("template: expanding %q in chain %q: %w", k, chainName, err)
		}
		resolved[k] = expanded
	}
	return resolved, nil
}

// expand recursively resolves placeholders within value using the provided
// variable map, tracking the resolution path to detect cycles.
func expand(value string, vars map[string]string, path []string) (string, error) {
	var expandErr error
	result := placeholder.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		name := match[2 : len(match)-1] // strip ${ and }
		for _, p := range path {
			if p == name {
				expandErr = fmt.Errorf("circular reference detected: %s", strings.Join(append(path, name), " -> "))
				return match
			}
		}
		ref, ok := vars[name]
		if !ok {
			return match // leave unresolved placeholders intact
		}
		inner, err := expand(ref, vars, append(path, name))
		if err != nil {
			expandErr = err
			return match
		}
		return inner
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}
