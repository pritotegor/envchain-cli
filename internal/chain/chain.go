package chain

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// EnvChain represents a named set of environment variables for a project.
type EnvChain struct {
	Name    string            `json:"name"`
	Project string            `json:"project"`
	Vars    map[string]string `json:"vars"`
}

// Store manages persistence of EnvChain entries on disk.
type Store struct {
	Path string
}

// NewStore creates a Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{Path: path}
}

// load reads all chains from disk. Returns empty slice if file doesn't exist.
func (s *Store) load() ([]EnvChain, error) {
	data, err := os.ReadFile(s.Path)
	if errors.Is(err, os.ErrNotExist) {
		return []EnvChain{}, nil
	}
	if err != nil {
		return nil, err
	}
	var chains []EnvChain
	if err := json.Unmarshal(data, &chains); err != nil {
		return nil, err
	}
	return chains, nil
}

// save writes all chains to disk, creating parent directories as needed.
func (s *Store) save(chains []EnvChain) error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(chains, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.Path, data, 0o600)
}

// Add inserts or replaces a chain with the same name+project combination.
func (s *Store) Add(chain EnvChain) error {
	chains, err := s.load()
	if err != nil {
		return err
	}
	for i, c := range chains {
		if c.Name == chain.Name && c.Project == chain.Project {
			chains[i] = chain
			return s.save(chains)
		}
	}
	chains = append(chains, chain)
	return s.save(chains)
}

// Get retrieves a chain by name and project.
func (s *Store) Get(name, project string) (*EnvChain, error) {
	chains, err := s.load()
	if err != nil {
		return nil, err
	}
	for _, c := range chains {
		if c.Name == name && c.Project == project {
			return &c, nil
		}
	}
	return nil, errors.New("chain not found")
}

// List returns all chains, optionally filtered by project (empty string = all).
func (s *Store) List(project string) ([]EnvChain, error) {
	chains, err := s.load()
	if err != nil {
		return nil, err
	}
	if project == "" {
		return chains, nil
	}
	var filtered []EnvChain
	for _, c := range chains {
		if c.Project == project {
			filtered = append(filtered, c)
		}
	}
	return filtered, nil
}

// Delete removes a chain by name and project. Returns error if not found.
func (s *Store) Delete(name, project string) error {
	chains, err := s.load()
	if err != nil {
		return err
	}
	for i, c := range chains {
		if c.Name == name && c.Project == project {
			chains = append(chains[:i], chains[i+1:]...)
			return s.save(chains)
		}
	}
	return errors.New("chain not found")
}
