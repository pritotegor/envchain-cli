package validate

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// validKeyPattern matches POSIX-compliant environment variable names.
var validKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ErrEmptyChainName is returned when a chain name is empty.
var ErrEmptyChainName = errors.New("chain name must not be empty")

// ErrInvalidChainName is returned when a chain name contains illegal characters.
var ErrInvalidChainName = errors.New("chain name must contain only letters, digits, hyphens, and underscores")

// ErrEmptyKey is returned when an environment variable key is empty.
var ErrEmptyKey = errors.New("variable key must not be empty")

// ErrInvalidKey is returned when an environment variable key is not a valid identifier.
var ErrInvalidKey = errors.New("variable key must start with a letter or underscore and contain only letters, digits, and underscores")

// validChainPattern matches chain names: letters, digits, hyphens, underscores.
var validChainPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// ChainName validates that name is a non-empty, safe chain identifier.
func ChainName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrEmptyChainName
	}
	if !validChainPattern.MatchString(name) {
		return fmt.Errorf("%w: %q", ErrInvalidChainName, name)
	}
	return nil
}

// Key validates that key is a non-empty, POSIX-compliant variable name.
func Key(key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	if !validKeyPattern.MatchString(key) {
		return fmt.Errorf("%w: %q", ErrInvalidKey, key)
	}
	return nil
}
