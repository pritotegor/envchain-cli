package merge_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-cli/internal/chain"
	"github.com/envchain-cli/internal/merge"
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
	for k, v := range vars {
		if err := st.Add(name, k, v); err != nil {
			t.Fatalf("seed %s/%s: %v", name, k, err)
		}
	}
}

func TestMerge_Basic(t *testing.T) {
	st := tempStore(t)
	seedChain(t, st, "alpha", map[string]string{"FOO": "1", "BAR": "2"})
	seedChain(t, st, "beta", map[string]string{"BAZ": "3"})

	m := merge.NewMerger(st)
	if err := m.Merge("gamma", []string{"alpha", "beta"}, merge.MergeOptions{}); err != nil {
		t.Fatalf("Merge: %v", err)
	}

	got, err := st.Get("gamma")
	if err != nil {
		t.Fatalf("Get gamma: %v", err)
	}
	for _, key := range []string{"FOO", "BAR", "BAZ"} {
		if _, ok := got[key]; !ok {
			t.Errorf("expected key %q in merged chain", key)
		}
	}
}

func TestMerge_NoOverwrite(t *testing.T) {
	st := tempStore(t)
	seedChain(t, st, "src", map[string]string{"KEY": "new"})
	seedChain(t, st, "dst", map[string]string{"KEY": "original"})

	m := merge.NewMerger(st)
	if err := m.Merge("dst", []string{"src"}, merge.MergeOptions{Overwrite: false}); err != nil {
		t.Fatalf("Merge: %v", err)
	}

	got, _ := st.Get("dst")
	if got["KEY"] != "original" {
		t.Errorf("expected original value, got %q", got["KEY"])
	}
}

func TestMerge_Overwrite(t *testing.T) {
	st := tempStore(t)
	seedChain(t, st, "src", map[string]string{"KEY": "new"})
	seedChain(t, st, "dst", map[string]string{"KEY": "original"})

	m := merge.NewMerger(st)
	if err := m.Merge("dst", []string{"src"}, merge.MergeOptions{Overwrite: true}); err != nil {
		t.Fatalf("Merge: %v", err)
	}

	got, _ := st.Get("dst")
	if got["KEY"] != "new" {
		t.Errorf("expected overwritten value, got %q", got["KEY"])
	}
}

func TestMerge_SourceNotFound(t *testing.T) {
	st := tempStore(t)
	m := merge.NewMerger(st)
	err := m.Merge("dst", []string{"nonexistent"}, merge.MergeOptions{})
	if err == nil {
		t.Fatal("expected error for missing source chain")
	}
}

func TestMerge_CreatesDestination(t *testing.T) {
	st := tempStore(t)
	seedChain(t, st, "src", map[string]string{"X": "42"})

	m := merge.NewMerger(st)
	if err := m.Merge("brand-new", []string{"src"}, merge.MergeOptions{}); err != nil {
		t.Fatalf("Merge: %v", err)
	}

	got, err := st.Get("brand-new")
	if err != nil {
		t.Fatalf("Get brand-new: %v", err)
	}
	if got["X"] != "42" {
		t.Errorf("expected X=42, got %q", got["X"])
	}
}

func init() {
	_ = os.Getenv // suppress unused import
}
