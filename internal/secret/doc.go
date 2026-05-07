// Package secret defines the Store interface for pluggable secret backends
// and ships a MemoryStore implementation for use in tests.
//
// A secret store allows envchain to resolve sensitive environment variable
// values from an external source (such as the OS keyring or a vault) rather
// than persisting them in the plain-text chain file.
//
// Usage:
//
//	var store secret.Store = secret.NewMemoryStore()
//	_ = store.Set("myproject", "API_KEY", "supersecret")
//	v, err := store.Get("myproject", "API_KEY")
//
// Implement the Store interface to integrate with any backend:
//
//	type KeyringStore struct{}
//	func (k *KeyringStore) Get(service, key string) (string, error)  { ... }
//	func (k *KeyringStore) Set(service, key, value string) error      { ... }
//	func (k *KeyringStore) Delete(service, key string) error          { ... }
//	func (k *KeyringStore) List(service string) ([]string, error)     { ... }
package secret
