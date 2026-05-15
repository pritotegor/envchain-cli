package compare_test

import (
	"os"
	"testing"

	"github.com/envchain-cli/internal/chain"
	"github.com/envchain-cli/internal/compare"
)

func tempStore(t *testing.T) *chain.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "compare-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := chain.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func seedChain(t *testing.T, s *chain.Store, name string, vars map[string]string) {
	t.Helper()
	for k, v := range vars {
		if err := s.Add(name, k, v); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCompare_NoChanges(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "a", map[string]string{"FOO": "1", "BAR": "2"})
	seedChain(t, s, "b", map[string]string{"FOO": "1", "BAR": "2"})

	res, err := compare.NewComparer(s).Compare("a", "b")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.OnlyInA) != 0 || len(res.OnlyInB) != 0 || len(res.DiffValue) != 0 {
		t.Errorf("expected no differences, got %+v", res)
	}
	if len(res.SameValue) != 2 {
		t.Errorf("expected 2 same keys, got %d", len(res.SameValue))
	}
}

func TestCompare_OnlyInA(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "a", map[string]string{"FOO": "1", "EXTRA": "x"})
	seedChain(t, s, "b", map[string]string{"FOO": "1"})

	res, err := compare.NewComparer(s).Compare("a", "b")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.OnlyInA) != 1 || res.OnlyInA[0] != "EXTRA" {
		t.Errorf("expected EXTRA only in A, got %v", res.OnlyInA)
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "a", map[string]string{"FOO": "1"})
	seedChain(t, s, "b", map[string]string{"FOO": "1", "NEW": "n"})

	res, err := compare.NewComparer(s).Compare("a", "b")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.OnlyInB) != 1 || res.OnlyInB[0] != "NEW" {
		t.Errorf("expected NEW only in B, got %v", res.OnlyInB)
	}
}

func TestCompare_DiffValue(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "a", map[string]string{"FOO": "old"})
	seedChain(t, s, "b", map[string]string{"FOO": "new"})

	res, err := compare.NewComparer(s).Compare("a", "b")
	if err != nil {
		t.Fatal(err)
	}
	v, ok := res.DiffValue["FOO"]
	if !ok {
		t.Fatal("expected FOO in DiffValue")
	}
	if v[0] != "old" || v[1] != "new" {
		t.Errorf("unexpected diff values: %v", v)
	}
}

func TestCompare_UnknownChain(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "a", map[string]string{"X": "1"})

	_, err := compare.NewComparer(s).Compare("a", "missing")
	if err == nil {
		t.Fatal("expected error for missing chain")
	}
}
