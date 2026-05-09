// Package merge provides functionality for merging environment variable chains.
//
// A merge operation reads all variables from one or more source chains and
// writes them into a destination chain. The caller controls whether existing
// destination keys are overwritten via MergeOptions.
//
// Basic usage:
//
//	st, _ := chain.NewStore(path)
//	m := merge.NewMerger(st)
//	err := m.Merge("production", []string{"base", "overrides"}, merge.MergeOptions{
//		Overwrite: true,
//	})
//
// Sources are applied left-to-right, so later sources take precedence when
// Overwrite is enabled.
package merge
