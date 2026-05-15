package grep_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-cli/envchain-cli/internal/chain"
	"github.com/envchain-cli/envchain-cli/internal/grep"
)

func tempStore(t *testing.T) *chain.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "grep-test-*")
	if err != nil {
		t.Fatalf("tempStore: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
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
			t.Fatalf("seed %s/%s: %v", name, k, err)
		}
	}
}

func TestSearch_MatchesKey(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "web", map[string]string{"DATABASE_URL": "postgres://", "PORT": "8080"})
	g := grep.NewGrepper(s)

	matches, err := g.Search("DATABASE", nil, false)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(matches) != 1 || matches[0].Key != "DATABASE_URL" {
		t.Fatalf("expected DATABASE_URL match, got %+v", matches)
	}
}

func TestSearch_MatchesValue(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "web", map[string]string{"DB": "postgres://localhost", "CACHE": "redis://"})
	g := grep.NewGrepper(s)

	matches, err := g.Search("postgres", nil, false)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(matches) != 1 || matches[0].Key != "DB" {
		t.Fatalf("expected DB match, got %+v", matches)
	}
}

func TestSearch_KeysOnly_SkipsValueMatch(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "web", map[string]string{"FOO": "postgres://", "BAR": "redis://"})
	g := grep.NewGrepper(s)

	matches, err := g.Search("postgres", nil, true)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected no matches in keysOnly mode, got %+v", matches)
	}
}

func TestSearch_AcrossMultipleChains(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "web", map[string]string{"API_KEY": "abc"})
	seedChain(t, s, "worker", map[string]string{"API_KEY": "xyz", "TIMEOUT": "30"})
	g := grep.NewGrepper(s)

	matches, err := g.Search("API_KEY", nil, false)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
}

func TestSearch_InvalidPattern(t *testing.T) {
	s := tempStore(t)
	g := grep.NewGrepper(s)
	_, err := g.Search("[invalid", nil, false)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestSearch_FilterByChain(t *testing.T) {
	s := tempStore(t)
	seedChain(t, s, "web", map[string]string{"SECRET": "web-secret"})
	seedChain(t, s, "worker", map[string]string{"SECRET": "worker-secret"})
	g := grep.NewGrepper(s)

	matches, err := g.Search("SECRET", []string{"web"}, false)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(matches) != 1 || matches[0].Chain != "web" {
		t.Fatalf("expected only web chain match, got %+v", matches)
	}
}
