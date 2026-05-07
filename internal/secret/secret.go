// Package secret provides an interface and implementations for
// integrating external secret stores (e.g. system keyring) with envchain.
package secret

import (
	"errors"
	"fmt"
)

// ErrNotFound is returned when a secret key does not exist in the store.
var ErrNotFound = errors.New("secret not found")

// Store defines the interface for a secret backend.
type Store interface {
	// Get retrieves the secret value for the given service and key.
	Get(service, key string) (string, error)
	// Set stores a secret value under the given service and key.
	Set(service, key, value string) error
	// Delete removes a secret from the store.
	Delete(service, key string) error
	// List returns all keys stored under the given service.
	List(service string) ([]string, error)
}

// MemoryStore is an in-memory implementation of Store, useful for testing.
type MemoryStore struct {
	data map[string]string
}

// NewMemoryStore creates a new empty MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]string)}
}

func memKey(service, key string) string {
	return fmt.Sprintf("%s/%s", service, key)
}

// Get retrieves a value from the in-memory store.
func (m *MemoryStore) Get(service, key string) (string, error) {
	v, ok := m.data[memKey(service, key)]
	if !ok {
		return "", fmt.Errorf("%w: %s/%s", ErrNotFound, service, key)
	}
	return v, nil
}

// Set stores a value in the in-memory store.
func (m *MemoryStore) Set(service, key, value string) error {
	m.data[memKey(service, key)] = value
	return nil
}

// Delete removes a value from the in-memory store.
func (m *MemoryStore) Delete(service, key string) error {
	k := memKey(service, key)
	if _, ok := m.data[k]; !ok {
		return fmt.Errorf("%w: %s/%s", ErrNotFound, service, key)
	}
	delete(m.data, k)
	return nil
}

// List returns all keys stored under the given service prefix.
func (m *MemoryStore) List(service string) ([]string, error) {
	prefix := service + "/"
	var keys []string
	for k := range m.data {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			keys = append(keys, k[len(prefix):])
		}
	}
	return keys, nil
}
