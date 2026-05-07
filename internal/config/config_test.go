package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envchain-cli/internal/config"
)

// overrideConfigDir temporarily redirects the OS config dir to a temp path.
func overrideConfigDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp) // Linux
	t.Setenv("AppData", tmp)          // Windows (best-effort)
	return tmp
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.SecretBackend != "memory" {
		t.Errorf("expected SecretBackend=memory, got %q", cfg.SecretBackend)
	}
	if cfg.AutoExportFormat != "shell" {
		t.Errorf("expected AutoExportFormat=shell, got %q", cfg.AutoExportFormat)
	}
}

func TestLoad_MissingFile_ReturnsDefault(t *testing.T) {
	overrideConfigDir(t)
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SecretBackend != "memory" {
		t.Errorf("expected default SecretBackend, got %q", cfg.SecretBackend)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	overrideConfigDir(t)

	want := &config.Config{
		DefaultStore:     "/tmp/mystore.json",
		SecretBackend:    "keyring",
		AutoExportFormat: "dotenv",
	}
	if err := config.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := config.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.DefaultStore != want.DefaultStore {
		t.Errorf("DefaultStore: got %q, want %q", got.DefaultStore, want.DefaultStore)
	}
	if got.SecretBackend != want.SecretBackend {
		t.Errorf("SecretBackend: got %q, want %q", got.SecretBackend, want.SecretBackend)
	}
	if got.AutoExportFormat != want.AutoExportFormat {
		t.Errorf("AutoExportFormat: got %q, want %q", got.AutoExportFormat, want.AutoExportFormat)
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	tmp := overrideConfigDir(t)
	cfg := config.DefaultConfig()
	if err := config.Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}

	dir, err := config.ConfigDir()
	if err != nil {
		t.Fatalf("ConfigDir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "config.json")); err != nil {
		t.Errorf("config file not created under %s: %v", tmp, err)
	}
}
