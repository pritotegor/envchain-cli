package chain_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envchain-cli/internal/chain"
)

func tempStore(t *testing.T) *chain.Store {
	t.Helper()
	dir := t.TempDir()
	return chain.NewStore(filepath.Join(dir, "chains.json"))
}

func TestAddAndGet(t *testing.T) {
	s := tempStore(t)
	c := chain.EnvChain{
		Name:    "dev",
		Project: "myapp",
		Vars:    map[string]string{"FOO": "bar", "PORT": "8080"},
	}
	if err := s.Add(c); err != nil {
		t.Fatalf("Add: %v", err)
	}
	got, err := s.Get("dev", "myapp")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", got.Vars["FOO"])
	}
}

func TestAddOverwritesExisting(t *testing.T) {
	s := tempStore(t)
	s.Add(chain.EnvChain{Name: "dev", Project: "myapp", Vars: map[string]string{"X": "1"}})
	s.Add(chain.EnvChain{Name: "dev", Project: "myapp", Vars: map[string]string{"X": "2"}})

	got, _ := s.Get("dev", "myapp")
	if got.Vars["X"] != "2" {
		t.Errorf("expected X=2 after overwrite, got %s", got.Vars["X"])
	}
}

func TestGetNotFound(t *testing.T) {
	s := tempStore(t)
	_, err := s.Get("missing", "proj")
	if err == nil {
		t.Error("expected error for missing chain")
	}
}

func TestList(t *testing.T) {
	s := tempStore(t)
	s.Add(chain.EnvChain{Name: "dev", Project: "alpha", Vars: map[string]string{}})
	s.Add(chain.EnvChain{Name: "prod", Project: "alpha", Vars: map[string]string{}})
	s.Add(chain.EnvChain{Name: "dev", Project: "beta", Vars: map[string]string{}})

	all, _ := s.List("")
	if len(all) != 3 {
		t.Errorf("expected 3 chains, got %d", len(all))
	}

	alpha, _ := s.List("alpha")
	if len(alpha) != 2 {
		t.Errorf("expected 2 alpha chains, got %d", len(alpha))
	}
}

func TestDelete(t *testing.T) {
	s := tempStore(t)
	s.Add(chain.EnvChain{Name: "dev", Project: "myapp", Vars: map[string]string{}})

	if err := s.Delete("dev", "myapp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Get("dev", "myapp")
	if err == nil {
		t.Error("expected chain to be deleted")
	}
}

func TestDeleteNotFound(t *testing.T) {
	s := tempStore(t)
	if err := s.Delete("ghost", "proj"); err == nil {
		t.Error("expected error when deleting non-existent chain")
	}
}

func TestStoreCreatesParentDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "deep", "chains.json")
	s := chain.NewStore(path)
	s.Add(chain.EnvChain{Name: "x", Project: "p", Vars: map[string]string{}})
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
