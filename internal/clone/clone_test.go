package clone_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-cli/envchain-cli/internal/chain"
	"github.com/envchain-cli/envchain-cli/internal/clone"
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
		t.Fatalf("seed Set(%q): %v", name, err)
	}
}

func TestClone_Basic(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "src", map[string]string{"FOO": "1", "BAR": "2"})
	c := clone.NewCloner(s)
	if err := c.Clone("src", "dst", clone.CloneOptions{}); err != nil {
		t.Fatalf("Clone: %v", err)
	}
	vars, _ := s.Get("dst")
	if vars["FOO"] != "1" || vars["BAR"] != "2" {
		t.Errorf("unexpected vars: %v", vars)
	}
}

func TestClone_SourceNotFound(t *testing.T) {
	s := tempStore(t)
	c := clone.NewCloner(s)
	if err := c.Clone("missing", "dst", clone.CloneOptions{}); err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestClone_DestinationExistsNoOverwrite(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "src", map[string]string{"A": "1"})
	seedChain(t, s, "dst", map[string]string{"B": "2"})
	c := clone.NewCloner(s)
	if err := c.Clone("src", "dst", clone.CloneOptions{Overwrite: false}); err == nil {
		t.Fatal("expected error when destination exists")
	}
}

func TestClone_DestinationExistsOverwrite(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "src", map[string]string{"X": "new"})
	seedChain(t, s, "dst", map[string]string{"X": "old"})
	c := clone.NewCloner(s)
	if err := c.Clone("src", "dst", clone.CloneOptions{Overwrite: true}); err != nil {
		t.Fatalf("Clone with overwrite: %v", err)
	}
	vars, _ := s.Get("dst")
	if vars["X"] != "new" {
		t.Errorf("expected overwritten value, got %q", vars["X"])
	}
}

func TestClone_PrefixFilter(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "src", map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "APP_NAME": "test"})
	c := clone.NewCloner(s)
	if err := c.Clone("src", "dst", clone.CloneOptions{Prefix: "DB_"}); err != nil {
		t.Fatalf("Clone with prefix: %v", err)
	}
	vars, _ := s.Get("dst")
	if len(vars) != 2 {
		t.Errorf("expected 2 keys, got %d: %v", len(vars), vars)
	}
	if _, ok := vars["APP_NAME"]; ok {
		t.Error("APP_NAME should have been filtered out")
	}
}

func TestClone_PrefixNoMatch(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "src", map[string]string{"FOO": "1"})
	c := clone.NewCloner(s)
	if err := c.Clone("src", "dst", clone.CloneOptions{Prefix: "NOPE_"}); err == nil {
		t.Fatal("expected error when prefix matches nothing")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
