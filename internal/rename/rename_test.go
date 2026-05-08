package rename_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envchain-cli/internal/chain"
	"github.com/yourorg/envchain-cli/internal/rename"
)

func tempStore(t *testing.T) *chain.Store {
	t.Helper()
	dir := t.TempDir()
	store, err := chain.NewStore(filepath.Join(dir, "chains.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func seedChain(t *testing.T, store *chain.Store, name string, vars map[string]string) {
	t.Helper()
	if err := store.Add(name, vars); err != nil {
		t.Fatalf("seed Add(%q): %v", name, err)
	}
}

func TestRename_Basic(t *testing.T) {
	store := tempStore(t)
	seedChain(t, store, "old", map[string]string{"FOO": "bar"})

	r := rename.NewRenamer(store)
	if err := r.Rename("old", "new", false); err != nil {
		t.Fatalf("Rename: %v", err)
	}

	if _, err := store.Get("old"); err == nil {
		t.Error("expected old chain to be removed, but it still exists")
	}

	vars, err := store.Get("new")
	if err != nil {
		t.Fatalf("Get new: %v", err)
	}
	if vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", vars["FOO"])
	}
}

func TestRename_SourceNotFound(t *testing.T) {
	store := tempStore(t)
	r := rename.NewRenamer(store)

	err := r.Rename("missing", "new", false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, rename.ErrSourceNotFound) {
		t.Errorf("expected ErrSourceNotFound, got %v", err)
	}
}

func TestRename_DestinationExistsNoOverwrite(t *testing.T) {
	store := tempStore(t)
	seedChain(t, store, "old", map[string]string{"A": "1"})
	seedChain(t, store, "existing", map[string]string{"B": "2"})

	r := rename.NewRenamer(store)
	err := r.Rename("old", "existing", false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, rename.ErrDestinationExists) {
		t.Errorf("expected ErrDestinationExists, got %v", err)
	}
}

func TestRename_DestinationExistsWithOverwrite(t *testing.T) {
	store := tempStore(t)
	seedChain(t, store, "old", map[string]string{"KEY": "new-value"})
	seedChain(t, store, "target", map[string]string{"KEY": "old-value"})

	r := rename.NewRenamer(store)
	if err := r.Rename("old", "target", true); err != nil {
		t.Fatalf("Rename with overwrite: %v", err)
	}

	vars, err := store.Get("target")
	if err != nil {
		t.Fatalf("Get target: %v", err)
	}
	if vars["KEY"] != "new-value" {
		t.Errorf("expected KEY=new-value after overwrite, got %q", vars["KEY"])
	}

	if _, err := store.Get("old"); err == nil {
		t.Error("expected old chain to be removed after overwrite rename")
	}
}

// Ensure the test file compiles even without explicit os import usage.
var _ = os.DevNull
var _ = filepath.Separator
