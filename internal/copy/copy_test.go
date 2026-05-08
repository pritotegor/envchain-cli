package copy_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-cli/envchain/internal/chain"
	envcopy "github.com/envchain-cli/envchain/internal/copy"
)

func tempStore(t *testing.T) *chain.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := chain.NewStore(filepath.Join(dir, "chains.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return st
}

func seedChain(t *testing.T, st *chain.Store, name string, vars map[string]string) {
	t.Helper()
	if err := st.Add(name, vars); err != nil {
		t.Fatalf("seed Add(%q): %v", name, err)
	}
}

func TestCopy_Basic(t *testing.T) {
	st := tempStore(t)
	seedChain(t, st, "src", map[string]string{"FOO": "bar", "BAZ": "qux"})

	cp := envcopy.NewCopier(st)
	if err := cp.Copy("src", "dst", false); err != nil {
		t.Fatalf("Copy: %v", err)
	}

	vars, err := st.Get("dst")
	if err != nil {
		t.Fatalf("Get dst: %v", err)
	}
	if vars["FOO"] != "bar" || vars["BAZ"] != "qux" {
		t.Errorf("unexpected vars: %v", vars)
	}
}

func TestCopy_SourceNotFound(t *testing.T) {
	st := tempStore(t)
	cp := envcopy.NewCopier(st)
	if err := cp.Copy("missing", "dst", false); err == nil {
		t.Fatal("expected error for missing source, got nil")
	}
}

func TestCopy_DestinationExistsNoOverwrite(t *testing.T) {
	st := tempStore(t)
	seedChain(t, st, "src", map[string]string{"A": "1"})
	seedChain(t, st, "dst", map[string]string{"B": "2"})

	cp := envcopy.NewCopier(st)
	if err := cp.Copy("src", "dst", false); err == nil {
		t.Fatal("expected error when dst exists and overwrite=false")
	}
}

func TestCopy_DestinationExistsWithOverwrite(t *testing.T) {
	st := tempStore(t)
	seedChain(t, st, "src", map[string]string{"A": "1"})
	seedChain(t, st, "dst", map[string]string{"B": "2"})

	cp := envcopy.NewCopier(st)
	if err := cp.Copy("src", "dst", true); err != nil {
		t.Fatalf("Copy with overwrite: %v", err)
	}

	vars, _ := st.Get("dst")
	if vars["A"] != "1" {
		t.Errorf("expected dst to have A=1 after overwrite, got %v", vars)
	}
}

func TestRename_MovesChain(t *testing.T) {
	st := tempStore(t)
	seedChain(t, st, "old", map[string]string{"X": "y"})

	cp := envcopy.NewCopier(st)
	if err := cp.Rename("old", "new", false); err != nil {
		t.Fatalf("Rename: %v", err)
	}

	if _, err := st.Get("old"); err == nil {
		t.Error("expected old chain to be deleted after rename")
	}
	vars, err := st.Get("new")
	if err != nil {
		t.Fatalf("Get new: %v", err)
	}
	if vars["X"] != "y" {
		t.Errorf("unexpected vars in new chain: %v", vars)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
