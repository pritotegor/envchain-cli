// Package config handles loading and saving envchain CLI configuration
// from a TOML file stored in the user's config directory.
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	appName    = "envchain"
	configFile = "config.json"
)

// Config holds the global envchain CLI configuration.
type Config struct {
	// DefaultStore is the path to the default chain store file.
	DefaultStore string `json:"default_store,omitempty"`
	// SecretBackend selects which secret backend to use (e.g. "memory", "keyring").
	SecretBackend string `json:"secret_backend,omitempty"`
	// AutoExportFormat sets the default export format (shell, dotenv, json).
	AutoExportFormat string `json:"auto_export_format,omitempty"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		SecretBackend:    "memory",
		AutoExportFormat: "shell",
	}
}

// ConfigDir returns the OS-appropriate config directory for envchain.
func ConfigDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, appName), nil
}

// Load reads the config file from the default config directory.
// If the file does not exist, DefaultConfig is returned without error.
func Load() (*Config, error) {
	dir, err := ConfigDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, configFile)

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return DefaultConfig(), nil
	}
	if err != nil {
		return nil, err
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the config to the default config directory, creating it if needed.
func Save(cfg *Config) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, configFile), data, 0o600)
}
