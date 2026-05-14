package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-cli/internal/chain"
	"github.com/envchain-cli/internal/snapshot"
)

func tempStore(t *testing.T) *chain.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "snapshot-test-*")
	if err != nil {
		t.Fatalf("tempStore: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	store, err := chain.NewStore(filepath.Join(dir, "chains.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func TestCapture_Basic(t *testing.T) {
	store := tempStore(t)
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := store.Add("mychain", vars); err != nil {
		t.Fatalf("Add: %v", err)
	}

	sn := snapshot.NewSnapshotter(store)
	snap, err := sn.Capture("mychain")
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}

	if snap.ChainName != "mychain" {
		t.Errorf("ChainName = %q; want %q", snap.ChainName, "mychain")
	}
	if snap.Vars["FOO"] != "bar" {
		t.Errorf("Vars[FOO] = %q; want %q", snap.Vars["FOO"], "bar")
	}
	if snap.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestCapture_UnknownChain(t *testing.T) {
	store := tempStore(t)
	sn := snapshot.NewSnapshotter(store)
	_, err := sn.Capture("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown chain, got nil")
	}
}

func TestRestore_Basic(t *testing.T) {
	store := tempStore(t)
	sn := snapshot.NewSnapshotter(store)

	snap := &snapshot.Snapshot{
		ChainName: "restored",
		Vars:      map[string]string{"KEY": "value"},
	}

	if err := sn.Restore(snap, false); err != nil {
		t.Fatalf("Restore: %v", err)
	}

	vars, err := store.Get("restored")
	if err != nil {
		t.Fatalf("Get after restore: %v", err)
	}
	if vars["KEY"] != "value" {
		t.Errorf("KEY = %q; want %q", vars["KEY"], "value")
	}
}

func TestRestore_NoOverwrite(t *testing.T) {
	store := tempStore(t)
	if err := store.Add("existing", map[string]string{"A": "1"}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	sn := snapshot.NewSnapshotter(store)
	snap := &snapshot.Snapshot{
		ChainName: "existing",
		Vars:      map[string]string{"B": "2"},
	}

	if err := sn.Restore(snap, false); err == nil {
		t.Fatal("expected error when overwrite=false and chain exists")
	}
}

func TestRestore_Overwrite(t *testing.T) {
	store := tempStore(t)
	if err := store.Add("chain", map[string]string{"OLD": "yes"}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	sn := snapshot.NewSnapshotter(store)
	snap := &snapshot.Snapshot{
		ChainName: "chain",
		Vars:      map[string]string{"NEW": "yes"},
	}

	if err := sn.Restore(snap, true); err != nil {
		t.Fatalf("Restore with overwrite: %v", err)
	}

	vars, _ := store.Get("chain")
	if vars["NEW"] != "yes" {
		t.Errorf("NEW = %q; want %q", vars["NEW"], "yes")
	}
}
