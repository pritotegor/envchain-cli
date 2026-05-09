package diff_test

import (
	"os"
	"path/filepath"
	"testing"

	"envchain-cli/internal/chain"
	"envchain-cli/internal/diff"
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
	for k, v := range vars {
		if err := s.Add(name, k, v, false); err != nil {
			t.Fatalf("seed Add %s/%s: %v", name, k, err)
		}
	}
}

func TestDiff_NoChanges(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "a", map[string]string{"FOO": "1", "BAR": "2"})
	seedChain(t, s, "b", map[string]string{"FOO": "1", "BAR": "2"})

	r, err := diff.NewDiffer(s).Diff("a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Changes()) != 0 {
		t.Errorf("expected 0 changes, got %d", len(r.Changes()))
	}
}

func TestDiff_Added(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "src", map[string]string{"FOO": "1"})
	seedChain(t, s, "dst", map[string]string{"FOO": "1", "BAR": "2"})

	r, err := diff.NewDiffer(s).Diff("src", "dst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	changes := r.Changes()
	if len(changes) != 1 || changes[0].Kind != diff.Added || changes[0].Key != "BAR" {
		t.Errorf("expected one Added BAR entry, got %+v", changes)
	}
}

func TestDiff_Removed(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "src", map[string]string{"FOO": "1", "BAR": "2"})
	seedChain(t, s, "dst", map[string]string{"FOO": "1"})

	r, err := diff.NewDiffer(s).Diff("src", "dst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	changes := r.Changes()
	if len(changes) != 1 || changes[0].Kind != diff.Removed || changes[0].Key != "BAR" {
		t.Errorf("expected one Removed BAR entry, got %+v", changes)
	}
}

func TestDiff_Modified(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "src", map[string]string{"FOO": "old"})
	seedChain(t, s, "dst", map[string]string{"FOO": "new"})

	r, err := diff.NewDiffer(s).Diff("src", "dst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	changes := r.Changes()
	if len(changes) != 1 || changes[0].Kind != diff.Modified {
		t.Errorf("expected one Modified entry, got %+v", changes)
	}
	if changes[0].OldValue != "old" || changes[0].NewValue != "new" {
		t.Errorf("wrong values: %+v", changes[0])
	}
}

func TestDiff_SourceNotFound(t *testing.T) {
	s := tempStore(t)
	_, err := diff.NewDiffer(s).Diff("ghost", "also-ghost")
	if err == nil {
		t.Fatal("expected error for missing chain")
	}
}

func TestDiff_EntryString(t *testing.T) {
	e := diff.Entry{Key: "X", Kind: diff.Added, NewValue: "42"}
	if got := e.String(); got != "+ X=42" {
		t.Errorf("unexpected String: %q", got)
	}
}

func TestDiff_EntryColored(t *testing.T) {
	e := diff.Entry{Key: "X", Kind: diff.Removed, OldValue: "old"}
	colored := e.Colored()
	if colored == e.String() {
		t.Error("expected colored output to differ from plain string")
	}
}

func init() { _ = os.Stderr }
