// Package clone implements deep-cloning of environment variable chains.
//
// A clone operation copies every key-value pair from a source chain into a
// newly named destination chain. An optional key prefix can be supplied to
// restrict which variables are included in the clone, making it easy to
// extract a logical sub-set of a large chain (e.g. all keys that start with
// "DB_") into its own dedicated chain.
//
// Basic usage:
//
//	cloner := clone.NewCloner(store)
//	err := cloner.Clone("production", "staging", clone.CloneOptions{
//		Overwrite: false,
//		Prefix:    "DB_",
//	})
package clone
