// Package rename provides functionality for renaming environment variable
// chains within the envchain store.
//
// A rename operation moves all variables from a source chain to a destination
// chain. By default the operation will fail if the destination chain already
// exists; pass Overwrite: true in the Options to replace it.
//
// Example usage:
//
//	r := rename.NewRenamer(store)
//	err := r.Rename(ctx, rename.Options{
//		Src:       "staging",
//		Dst:       "production",
//		Overwrite: false,
//	})
package rename
