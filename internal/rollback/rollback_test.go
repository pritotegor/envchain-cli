package rollback_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-cli/internal/chain"
	"github.com/envchain-cli/internal/rollback"
	"github.com/envchain-cli/internal/snapshot"
)

func tempStore(t *testing.T) *chain.Store {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "chains")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	store, err := chain.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func seedChain(t *testing.T, store *chain.Store, name string, vars map[string]string) {
	t.Helper()
	for k, v := range vars {
		if err := store.Set(name, k, v); err != nil {
			t.Fatalf("seed: set %s/%s: %v", name, k, err)
		}
	}
}

func TestRollback_Basic(t *testing.T) {
	store := tempStore(t)
	seedChain(t, store, "app", map[string]string{"FOO": "bar", "BAZ": "qux"})

	snapper := snapshot.NewSnapshotter(store)
	if err := snapper.Capture("app", "snap1"); err != nil {
		t.Fatalf("capture: %v", err)
	}

	// Mutate the chain after snapshot.
	if err := store.Set("app", "FOO", "changed"); err != nil {
		t.Fatal(err)
	}

	rb := rollback.NewRollbacker(store)
	if err := rb.Rollback("app", "snap1", true); err != nil {
		t.Fatalf("rollback: %v", err)
	}

	vars, err := store.Get("app")
	if err != nil {
		t.Fatal(err)
	}
	if vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", vars["FOO"])
	}
}

func TestRollback_NoOverwrite_Fails(t *testing.T) {
	store := tempStore(t)
	seedChain(t, store, "app", map[string]string{"A": "1"})

	snapper := snapshot.NewSnapshotter(store)
	if err := snapper.Capture("app", "snap2"); err != nil {
		t.Fatalf("capture: %v", err)
	}

	rb := rollback.NewRollbacker(store)
	err := rb.Rollback("app", "snap2", false)
	if err == nil {
		t.Fatal("expected error when overwrite=false and chain has keys")
	}
}

func TestRollback_UnknownSnapshot(t *testing.T) {
	store := tempStore(t)
	rb := rollback.NewRollbacker(store)
	err := rb.Rollback("app", "no-such-snap", true)
	if err == nil {
		t.Fatal("expected error for unknown snapshot")
	}
}
