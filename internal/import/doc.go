// Package envimport provides an Importer that reads key=value pairs from
// a reader (e.g. a .env file or shell export script) and stores them into
// a named envchain chain.
//
// Supported input formats:
//
//	# plain dotenv
//	KEY=value
//
//	# shell export style
//	export KEY=value
//
//	# quoted values (single or double quotes are stripped)
//	KEY="some value"
//	KEY='some value'
//
// Comments (lines beginning with '#') and blank lines are silently ignored.
// When overwrite is false, keys that already exist in the target chain are
// skipped and counted in Result.Skipped.
package envimport
