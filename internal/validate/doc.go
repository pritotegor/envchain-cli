// Package validate provides input validation helpers for envchain-cli.
//
// It enforces naming rules for chain identifiers and environment variable keys
// before they are persisted to the store, preventing malformed data and
// potential injection issues.
//
// Chain names must consist of letters, digits, hyphens, and underscores.
// Variable keys must follow POSIX conventions: start with a letter or
// underscore, followed by letters, digits, or underscores.
package validate
