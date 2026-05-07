package secret_test

import (
	"errors"
	"sort"
	"testing"

	"github.com/yourorg/envchain-cli/internal/secret"
)

func TestMemoryStore_SetAndGet(t *testing.T) {
	s := secret.NewMemoryStore()
	if err := s.Set("myapp", "DB_PASS", "s3cr3t"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	v, err := s.Get("myapp", "DB_PASS")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if v != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %q", v)
	}
}

func TestMemoryStore_GetNotFound(t *testing.T) {
	s := secret.NewMemoryStore()
	_, err := s.Get("myapp", "MISSING")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, secret.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	s := secret.NewMemoryStore()
	_ = s.Set("myapp", "TOKEN", "abc")
	if err := s.Delete("myapp", "TOKEN"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Get("myapp", "TOKEN")
	if !errors.Is(err, secret.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestMemoryStore_DeleteNotFound(t *testing.T) {
	s := secret.NewMemoryStore()
	err := s.Delete("myapp", "GHOST")
	if !errors.Is(err, secret.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryStore_List(t *testing.T) {
	s := secret.NewMemoryStore()
	_ = s.Set("proj", "KEY_A", "1")
	_ = s.Set("proj", "KEY_B", "2")
	_ = s.Set("other", "KEY_C", "3")

	keys, err := s.List("proj")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "KEY_A" || keys[1] != "KEY_B" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestMemoryStore_ListEmpty(t *testing.T) {
	s := secret.NewMemoryStore()
	keys, err := s.List("nonexistent")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("expected empty list, got %v", keys)
	}
}
