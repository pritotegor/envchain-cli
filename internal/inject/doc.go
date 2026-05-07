// Package inject provides utilities for running subprocesses with environment
// variables sourced from a named envchain store.
//
// A Runner merges the variables stored under a chain name into the current
// process environment, giving chain variables precedence over any identically
// named variables already present. The merged environment is then passed to
// the child process, leaving the parent process environment unchanged.
//
// Basic usage:
//
//	store, _ := chain.NewStore("/path/to/chains.json")
//	r := inject.NewRunner(store)
//
//	// Execute a command with injected vars:
//	err := r.Run("myproject", []string{"go", "test", "./..."})
//
//	// Inspect the merged environment without running anything:
//	env, err := r.Environ("myproject")
package inject
