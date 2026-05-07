package inject

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/user/envchain-cli/internal/chain"
)

// Runner executes a command with injected environment variables from a named chain.
type Runner struct {
	store *chain.Store
}

// NewRunner creates a new Runner backed by the given Store.
func NewRunner(store *chain.Store) *Runner {
	return &Runner{store: store}
}

// Run executes the provided command with the environment variables from the
// named chain merged into the current process environment.
func (r *Runner) Run(chainName string, command []string) error {
	if len(command) == 0 {
		return fmt.Errorf("inject: no command provided")
	}

	vars, err := r.store.Get(chainName)
	if err != nil {
		return fmt.Errorf("inject: chain %q not found: %w", chainName, err)
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Start with the current environment.
	env := os.Environ()

	// Build a set of keys already present so we can override them.
	existing := make(map[string]int, len(env))
	for i, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			existing[parts[0]] = i
		}
	}

	for k, v := range vars {
		entry := k + "=" + v
		if idx, ok := existing[k]; ok {
			env[idx] = entry
		} else {
			env = append(env, entry)
		}
	}

	cmd.Env = env
	return cmd.Run()
}

// Environ returns the merged environment as a slice of KEY=VALUE strings
// without executing any command. Useful for inspection or export.
func (r *Runner) Environ(chainName string) ([]string, error) {
	vars, err := r.store.Get(chainName)
	if err != nil {
		return nil, fmt.Errorf("inject: chain %q not found: %w", chainName, err)
	}

	env := os.Environ()
	existing := make(map[string]int, len(env))
	for i, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			existing[parts[0]] = i
		}
	}

	for k, v := range vars {
		entry := k + "=" + v
		if idx, ok := existing[k]; ok {
			env[idx] = entry
		} else {
			env = append(env, entry)
		}
	}

	return env, nil
}
