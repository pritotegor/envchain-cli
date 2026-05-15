package prune_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-cli/envchain-cli/internal/chain"
	"github.com/envchain-cli/envchain-cli/internal/prune"
)

func tempStore(t *testing.T) *chain.Store {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "chains.json")
	s, err := chain.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func seedChain(t *testing.T, s *chain.Store, name string, vars map[string]string) {
	t.Helper()
	if err := s.Add(name, vars); err != nil {
		t.Fatalf("seed Add(%q): %v", name, err)
	}
}

func TestPrune_RemovesEmptyKeys(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "app", map[string]string{"A": "hello", "B": "", "C": "world"})

	p := prune.NewPruner(s)
	res, err := p.Prune([]string{"app"}, false)
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}

	if len(res.RemovedKeys["app"]) != 1 || res.RemovedKeys["app"][0] != "B" {
		t.Errorf("expected [B] removed, got %v", res.RemovedKeys["app"])
	}

	vars, _ := s.Get("app")
	if _, ok := vars["B"]; ok {
		t.Error("key B should have been deleted")
	}
	if vars["A"] != "hello" || vars["C"] != "world" {
		t.Error("non-empty keys should be preserved")
	}
}

func TestPrune_RemovesEmptyChain(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "empty", map[string]string{"X": ""})

	p := prune.NewPruner(s)
	res, err := p.Prune([]string{"empty"}, true)
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}

	if len(res.RemovedChains) != 1 || res.RemovedChains[0] != "empty" {
		t.Errorf("expected chain 'empty' removed, got %v", res.RemovedChains)
	}
}

func TestPrune_KeepsChainWhenFlagFalse(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "solo", map[string]string{"K": ""})

	p := prune.NewPruner(s)
	res, err := p.Prune([]string{"solo"}, false)
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}

	if len(res.RemovedChains) != 0 {
		t.Errorf("chain should not be removed when flag is false")
	}
}

func TestPrune_AllChains(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "a", map[string]string{"X": "", "Y": "keep"})
	seedChain(t, s, "b", map[string]string{"Z": ""})

	p := prune.NewPruner(s)
	res, err := p.Prune(nil, true)
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}

	if len(res.RemovedChains) != 1 || res.RemovedChains[0] != "b" {
		t.Errorf("expected chain 'b' removed, got %v", res.RemovedChains)
	}
	if len(res.RemovedKeys["a"]) != 1 {
		t.Errorf("expected one key removed from 'a', got %v", res.RemovedKeys["a"])
	}
}

func TestPrune_NoEmptyKeys(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "full", map[string]string{"A": "1", "B": "2"})

	p := prune.NewPruner(s)
	res, err := p.Prune([]string{"full"}, true)
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}

	if len(res.RemovedKeys) != 0 || len(res.RemovedChains) != 0 {
		t.Errorf("nothing should be removed, got %+v", res)
	}
}

func init() {
	// Ensure os is imported (used implicitly via t.TempDir internals).
	_ = os.DevNull
}
