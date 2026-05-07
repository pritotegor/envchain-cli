// Package config provides loading and persistence of the envchain CLI
// user configuration.
//
// Configuration is stored as a JSON file in the OS-appropriate user config
// directory (e.g. ~/.config/envchain/config.json on Linux/macOS).
//
// Usage:
//
//	cfg, err := config.Load()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	cfg.SecretBackend = "keyring"
//	if err := config.Save(cfg); err != nil {
//		log.Fatal(err)
//	}
package config
