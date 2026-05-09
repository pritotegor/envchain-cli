package diff_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-cli/internal/chain"
	"github.com/envchain-cli/internal/diff"
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
		if err := s.Add(name, k, v); err != nil {
			t.Fatalf("seed Add(%s, %s): %v", name, k, err)
		}
	}
}

func TestDiff_NoChanges(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "a", map[string]string{"FOO": "1", "BAR": "2"})
	seedChain(t, s, "b", map[string]string{"FOO": "1", "BAR": "2"})

	d := diff.NewDiffer(s)
	res, err := d.Compare("a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.HasChanges() {
		t.Errorf("expected no changes, got %+v", res.Changes)
	}
}

func TestDiff_Added(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "a", map[string]string{"FOO": "1"})
	seedChain(t, s, "b", map[string]string{"FOO": "1", "BAR": "2"})

	d := diff.NewDiffer(s)
	res, err := d.Compare("a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changes) != 1 || res.Changes[0].Kind != diff.Added || res.Changes[0].Key != "BAR" {
		t.Errorf("expected one Added change for BAR, got %+v", res.Changes)
	}
}

func TestDiff_Removed(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "a", map[string]string{"FOO": "1", "BAR": "2"})
	seedChain(t, s, "b", map[string]string{"FOO": "1"})

	d := diff.NewDiffer(s)
	res, err := d.Compare("a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changes) != 1 || res.Changes[0].Kind != diff.Removed || res.Changes[0].Key != "BAR" {
		t.Errorf("expected one Removed change for BAR, got %+v", res.Changes)
	}
}

func TestDiff_Modified(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "a", map[string]string{"FOO": "old"})
	seedChain(t, s, "b", map[string]string{"FOO": "new"})

	d := diff.NewDiffer(s)
	res, err := d.Compare("a", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(res.Changes))
	}
	c := res.Changes[0]
	if c.Kind != diff.Modified || c.OldValue != "old" || c.NewValue != "new" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_SourceNotFound(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "b", map[string]string{"FOO": "1"})

	d := diff.NewDiffer(s)
	_, err := d.Compare("missing", "b")
	if err == nil {
		t.Error("expected error for missing source chain")
	}
	_ = os.Getenv("CI") // suppress unused import warning
}

func TestDiff_DestinationNotFound(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "a", map[string]string{"FOO": "1"})

	d := diff.NewDiffer(s)
	_, err := d.Compare("a", "missing")
	if err == nil {
		t.Error("expected error for missing destination chain")
	}
}
