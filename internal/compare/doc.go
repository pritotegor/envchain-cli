// Package compare provides side-by-side comparison of two envchain variable
// sets, classifying each key as unique to one side, shared with the same
// value, or shared with differing values.
//
// Usage:
//
//	cmp := compare.NewComparer(store)
//	res, err := cmp.Compare("staging", "production")
//	if err != nil { ... }
//	// res.OnlyInA  — keys only in "staging"
//	// res.OnlyInB  — keys only in "production"
//	// res.DiffValue — keys whose values differ
//	// res.SameValue — keys that are identical
package compare
