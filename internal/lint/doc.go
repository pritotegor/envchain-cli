// Package lint provides chain variable linting utilities for envchain-cli.
//
// A Linter inspects one or more chains stored in a chain.Store and reports
// style or correctness issues found in variable names and values.
//
// Supported checks:
//   - Keys must match the canonical ALL_CAPS_UNDERSCORE pattern
//   - Keys that are all-lowercase are flagged as a softer warning
//   - Empty values are flagged as potential oversights
//
// Usage:
//
//	l := lint.NewLinter(store)
//	issues, err := l.Lint("") // empty string = all chains
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, issue := range issues {
//	    fmt.Println(issue)
//	}
package lint
