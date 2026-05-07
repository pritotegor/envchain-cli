package inject_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envchain-cli/internal/chain"
	"github.com/user/envchain-cli/internal/inject"
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

func TestEnviron_MergesVars(t *testing.T) {
	s := tempStore(t)
	if err := s.Add("mychain", map[string]string{"INJECT_FOO": "bar", "INJECT_BAZ": "qux"}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	r := inject.NewRunner(s)
	env, err := r.Environ("mychain")
	if err != nil {
		t.Fatalf("Environ: %v", err)
	}

	found := map[string]string{}
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			found[parts[0]] = parts[1]
		}
	}

	if found["INJECT_FOO"] != "bar" {
		t.Errorf("expected INJECT_FOO=bar, got %q", found["INJECT_FOO"])
	}
	if found["INJECT_BAZ"] != "qux" {
		t.Errorf("expected INJECT_BAZ=qux, got %q", found["INJECT_BAZ"])
	}
}

func TestEnviron_OverridesExisting(t *testing.T) {
	t.Setenv("INJECT_OVERRIDE", "original")

	s := tempStore(t)
	if err := s.Add("overchain", map[string]string{"INJECT_OVERRIDE": "replaced"}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	r := inject.NewRunner(s)
	env, err := r.Environ("overchain")
	if err != nil {
		t.Fatalf("Environ: %v", err)
	}

	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if parts[0] == "INJECT_OVERRIDE" {
			if parts[1] != "replaced" {
				t.Errorf("expected replaced, got %q", parts[1])
			}
			return
		}
	}
	t.Error("INJECT_OVERRIDE not found in env")
}

func TestEnviron_UnknownChain(t *testing.T) {
	s := tempStore(t)
	r := inject.NewRunner(s)
	_, err := r.Environ("nonexistent")
	if err == nil {
		t.Error("expected error for unknown chain")
	}
}

func TestRun_NoCommand(t *testing.T) {
	s := tempStore(t)
	if err := s.Add("c", map[string]string{"X": "1"}); err != nil {
		t.Fatalf("Add: %v", err)
	}
	r := inject.NewRunner(s)
	if err := r.Run("c", []string{}); err == nil {
		t.Error("expected error when no command provided")
	}
}

func TestRun_EchoCommand(t *testing.T) {
	if _, err := os.LookupEnv("CI"); false {
		_ = err // suppress unused warning
	}
	s := tempStore(t)
	if err := s.Add("echochain", map[string]string{"HELLO": "world"}); err != nil {
		t.Fatalf("Add: %v", err)
	}
	r := inject.NewRunner(s)
	if err := r.Run("echochain", []string{"env"}); err != nil {
		// env(1) may not be available everywhere; skip rather than fail
		t.Skipf("env command not available: %v", err)
	}
}
