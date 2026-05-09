// Package diff provides functionality for comparing two environment variable
// chains and reporting the differences between them.
//
// A Differ computes the symmetric difference between two named chains stored
// in a chain.Store, returning a Result that categorises every key as Added,
// Removed, Modified, or Unchanged relative to the source chain.
//
// Example usage:
//
//	d := diff.NewDiffer(store)
//	result, err := d.Diff("staging", "production")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, entry := range result.Changes() {
//		fmt.Println(entry)
//	}
package diff
