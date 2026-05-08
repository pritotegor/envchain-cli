// Package copy implements chain duplication and rename operations for
// envchain stores.
//
// # Overview
//
// A [Copier] wraps a [chain.Store] and exposes two high-level operations:
//
//   - Copy – duplicates a named chain into a new name, with an optional
//     overwrite flag to replace an existing destination.
//
//   - Rename – performs a Copy followed by deletion of the source chain,
//     providing an atomic-feeling move within the same store file.
//
// # Example
//
//	cp := copy.NewCopier(store)
//
//	// Duplicate "production" into "staging" without clobbering existing data.
//	if err := cp.Copy("production", "staging", false); err != nil {
//		log.Fatal(err)
//	}
package copy
