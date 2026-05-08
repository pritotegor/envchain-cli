package delete_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envchain-cli/internal/chain"
	"github.com/yourorg/envchain-cli/internal/delete"
)

func tempStore(t *testing.T) *chain.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := chain.NewStore(filepath.Join(dir, "chains.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func seedChain(t *testing.T, s *chain.Store, name string, vars map[string]string) {
	t.Helper()
	if err := s.Set(name, vars); err != nil {
		t.Fatalf("seed Set: %v", err)
	}
}

func TestDeleteChain_Basic(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "mychain", map[string]string{"FOO": "bar"})

	d := delete.NewDeleter(s)
	if err := d.DeleteChain("mychain"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := s.Get("mychain"); err == nil {
		t.Fatal("expected chain to be gone, but Get succeeded")
	}
}

func TestDeleteChain_NotFound(t *testing.T) {
	s := tempStore(t)
	d := delete.NewDeleter(s)

	err := d.DeleteChain("ghost")
	if !errors.Is(err, delete.ErrChainNotFound) {
		t.Fatalf("expected ErrChainNotFound, got %v", err)
	}
}

func TestDeleteVar_Basic(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "mychain", map[string]string{"FOO": "bar", "BAZ": "qux"})

	d := delete.NewDeleter(s)
	if err := d.DeleteVar("mychain", "FOO"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vars, err := s.Get("mychain")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if _, ok := vars["FOO"]; ok {
		t.Error("FOO should have been deleted")
	}
	if vars["BAZ"] != "qux" {
		t.Errorf("BAZ should remain; got %q", vars["BAZ"])
	}
}

func TestDeleteVar_ChainNotFound(t *testing.T) {
	s := tempStore(t)
	d := delete.NewDeleter(s)

	err := d.DeleteVar("ghost", "KEY")
	if !errors.Is(err, delete.ErrChainNotFound) {
		t.Fatalf("expected ErrChainNotFound, got %v", err)
	}
}

func TestDeleteVar_KeyNotFound(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "mychain", map[string]string{"FOO": "bar"})

	d := delete.NewDeleter(s)
	err := d.DeleteVar("mychain", "MISSING")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
	_ = os.Stderr // suppress unused import warning
}
