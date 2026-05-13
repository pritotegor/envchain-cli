package validate_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envchain-cli/internal/validate"
)

func TestChainName_Valid(t *testing.T) {
	valid := []string{"myproject", "my-project", "my_project", "Project123", "_hidden"}
	for _, name := range valid {
		if err := validate.ChainName(name); err != nil {
			t.Errorf("expected %q to be valid, got: %v", name, err)
		}
	}
}

func TestChainName_Empty(t *testing.T) {
	if err := validate.ChainName(""); !errors.Is(err, validate.ErrEmptyChainName) {
		t.Errorf("expected ErrEmptyChainName, got: %v", err)
	}
	if err := validate.ChainName("   "); !errors.Is(err, validate.ErrEmptyChainName) {
		t.Errorf("expected ErrEmptyChainName for whitespace, got: %v", err)
	}
}

func TestChainName_Invalid(t *testing.T) {
	invalid := []string{"my project", "my/project", "my.project", "my@project"}
	for _, name := range invalid {
		err := validate.ChainName(name)
		if !errors.Is(err, validate.ErrInvalidChainName) {
			t.Errorf("expected ErrInvalidChainName for %q, got: %v", name, err)
		}
	}
}

func TestKey_Valid(t *testing.T) {
	valid := []string{"FOO", "_BAR", "my_var", "VAR123", "_"}
	for _, k := range valid {
		if err := validate.Key(k); err != nil {
			t.Errorf("expected %q to be valid, got: %v", k, err)
		}
	}
}

func TestKey_Empty(t *testing.T) {
	if err := validate.Key(""); !errors.Is(err, validate.ErrEmptyKey) {
		t.Errorf("expected ErrEmptyKey, got: %v", err)
	}
}

func TestKey_Invalid(t *testing.T) {
	invalid := []string{"1VAR", "my-var", "my var", "VAR.NAME", "VAR=1"}
	for _, k := range invalid {
		err := validate.Key(k)
		if !errors.Is(err, validate.ErrInvalidKey) {
			t.Errorf("expected ErrInvalidKey for %q, got: %v", k, err)
		}
	}
}
