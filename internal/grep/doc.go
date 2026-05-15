// Package grep implements pattern-based search across environment variable
// chains stored by envchain-cli.
//
// It supports:
//   - Regular expression matching against key names and/or values
//   - Scoping searches to a specific subset of chains
//   - A keys-only mode that restricts matching to key names, ignoring values
//
// Example usage:
//
//	g := grep.NewGrepper(store)
//	matches, err := g.Search("DATABASE", nil, false)
//	for _, m := range matches {
//	    fmt.Printf("%s  %s=%s\n", m.Chain, m.Key, m.Value)
//	}
package grep
