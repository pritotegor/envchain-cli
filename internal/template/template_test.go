package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-cli/envchain-cli/internal/chain"
	"github.com/envchain-cli/envchain-cli/internal/template"
)

func tempStore(t *testing.T) *chain.Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "template-test-*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	store, err := chain.NewStore(filepath.Join(dir, "chains.json"))
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	return store
}

func TestRender_NoPlaceholders(t *testing.T) {
	store := tempStore(t)
	_ = store.Add("app", map[string]string{"HOST": "localhost", "PORT": "8080"})
	r := template.NewRenderer(store)
	out, err := r.Render("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" || out["PORT"] != "8080" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestRender_ResolvesPlaceholder(t *testing.T) {
	store := tempStore(t)
	_ = store.Add("app", map[string]string{
		"BASE_URL": "http://${HOST}:${PORT}",
		"HOST":     "example.com",
		"PORT":     "443",
	})
	r := template.NewRenderer(store)
	out, err := r.Render("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["BASE_URL"] != "http://example.com:443" {
		t.Errorf("got BASE_URL=%q, want http://example.com:443", out["BASE_URL"])
	}
}

func TestRender_UnresolvedPlaceholderLeftIntact(t *testing.T) {
	store := tempStore(t)
	_ = store.Add("app", map[string]string{"GREETING": "Hello, ${NAME}!"})
	r := template.NewRenderer(store)
	out, err := r.Render("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["GREETING"] != "Hello, ${NAME}!" {
		t.Errorf("got %q, want original placeholder", out["GREETING"])
	}
}

func TestRender_CircularReference(t *testing.T) {
	store := tempStore(t)
	_ = store.Add("app", map[string]string{
		"A": "${B}",
		"B": "${A}",
	})
	r := template.NewRenderer(store)
	_, err := r.Render("app")
	if err == nil {
		t.Fatal("expected circular reference error, got nil")
	}
}

func TestRender_UnknownChain(t *testing.T) {
	store := tempStore(t)
	r := template.NewRenderer(store)
	_, err := r.Render("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown chain, got nil")
	}
}

func TestRender_NestedPlaceholders(t *testing.T) {
	store := tempStore(t)
	_ = store.Add("app", map[string]string{
		"PROTO":    "https",
		"HOST":     "api.example.com",
		"BASE":     "${PROTO}://${HOST}",
		"ENDPOINT": "${BASE}/v1",
	})
	r := template.NewRenderer(store)
	out, err := r.Render("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ENDPOINT"] != "https://api.example.com/v1" {
		t.Errorf("got ENDPOINT=%q, want https://api.example.com/v1", out["ENDPOINT"])
	}
}
