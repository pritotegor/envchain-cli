// Package snapshot provides functionality to capture and restore
// the state of an environment chain at a point in time.
package snapshot

import (
	"fmt"
	"time"

	"github.com/envchain-cli/internal/chain"
)

// Snapshot represents a point-in-time capture of a chain's variables.
type Snapshot struct {
	ChainName string            `json:"chain_name"`
	CapturedAt time.Time        `json:"captured_at"`
	Vars       map[string]string `json:"vars"`
}

// Snapshotter captures and restores chain snapshots.
type Snapshotter struct {
	store *chain.Store
}

// NewSnapshotter returns a Snapshotter backed by the given store.
func NewSnapshotter(store *chain.Store) *Snapshotter {
	return &Snapshotter{store: store}
}

// Capture creates a snapshot of the named chain's current variables.
func (s *Snapshotter) Capture(chainName string) (*Snapshot, error) {
	vars, err := s.store.Get(chainName)
	if err != nil {
		return nil, fmt.Errorf("snapshot: capture %q: %w", chainName, err)
	}

	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}

	return &Snapshot{
		ChainName:  chainName,
		CapturedAt: time.Now().UTC(),
		Vars:       copy,
	}, nil
}

// Restore overwrites the named chain's variables with those from the snapshot.
// If overwrite is false and the chain already has keys, an error is returned.
func (s *Snapshotter) Restore(snap *Snapshot, overwrite bool) error {
	existing, err := s.store.Get(snap.ChainName)
	if err == nil && len(existing) > 0 && !overwrite {
		return fmt.Errorf("snapshot: restore %q: chain already exists; use overwrite to replace", snap.ChainName)
	}

	if err := s.store.Add(snap.ChainName, snap.Vars); err != nil {
		return fmt.Errorf("snapshot: restore %q: %w", snap.ChainName, err)
	}
	return nil
}

// Age returns the duration elapsed since the snapshot was captured.
func (snap *Snapshot) Age() time.Duration {
	return time.Since(snap.CapturedAt)
}

// Summary returns a human-readable one-line description of the snapshot,
// including the chain name, number of variables, and time since capture.
func (snap *Snapshot) Summary() string {
	return fmt.Sprintf("chain=%q vars=%d age=%s", snap.ChainName, len(snap.Vars), snap.Age().Round(time.Second))
}

// Diff returns the keys that differ between this snapshot and another.
// It returns three slices: keys added in other, keys removed in other,
// and keys whose values changed between the two snapshots.
func (snap *Snapshot) Diff(other *Snapshot) (added, removed, changed []string) {
	for k, v := range other.Vars {
		if orig, ok := snap.Vars[k]; !ok {
			added = append(added, k)
		} else if orig != v {
			changed = append(changed, k)
		}
	}
	for k := range snap.Vars {
		if _, ok := other.Vars[k]; !ok {
			removed = append(removed, k)
		}
	}
	return added, removed, changed
}
