// Package envimport provides functionality to import environment variables
// from an existing .env file or shell export statements into a named chain.
package envimport

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/user/envchain-cli/internal/chain"
)

// Result holds the outcome of an import operation.
type Result struct {
	ChainName string
	Imported  int
	Skipped   int
	Overwrote int
}

// Importer reads key=value pairs and stores them into a chain.
type Importer struct {
	store *chain.Store
}

// NewImporter returns an Importer backed by the given store.
func NewImporter(s *chain.Store) *Importer {
	return &Importer{store: s}
}

// Import reads lines from r and adds parsed variables to chainName.
// Lines beginning with '#' or that are blank are ignored.
// Lines of the form `export KEY=VALUE` are also accepted.
// When overwrite is false, existing keys are skipped.
func (im *Importer) Import(r io.Reader, chainName string, overwrite bool) (Result, error) {
	res := Result{ChainName: chainName}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Strip optional "export " prefix.
		line = strings.TrimPrefix(line, "export ")

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			res.Skipped++
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), "\"'")

		if key == "" {
			res.Skipped++
			continue
		}

		existing, _ := im.store.Get(chainName, key)
		if existing != "" && !overwrite {
			res.Skipped++
			continue
		}

		if err := im.store.Add(chainName, key, value); err != nil {
			return res, fmt.Errorf("adding %s to chain %q: %w", key, chainName, err)
		}

		if existing != "" {
			res.Overwrote++
		} else {
			res.Imported++
		}
	}

	if err := scanner.Err(); err != nil {
		return res, fmt.Errorf("reading input: %w", err)
	}
	return res, nil
}
