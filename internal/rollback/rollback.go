package rollback

import (
	"fmt"

	"github.com/envchain-cli/internal/chain"
	"github.com/envchain-cli/internal/snapshot"
)

// Rollbacker restores a chain to a previously captured snapshot.
type Rollbacker struct {
	store     *chain.Store
	snapper   *snapshot.Snapshotter
}

// NewRollbacker creates a Rollbacker backed by the given store.
func NewRollbacker(store *chain.Store) *Rollbacker {
	return &Rollbacker{
		store:   store,
		snapper: snapshot.NewSnapshotter(store),
	}
}

// Rollback replaces the current contents of chainName with the snapshot
// identified by snapshotName. If overwrite is false and chainName already
// contains keys the operation is aborted.
func (r *Rollbacker) Rollback(chainName, snapshotName string, overwrite bool) error {
	current, err := r.store.Get(chainName)
	if err == nil && len(current) > 0 && !overwrite {
		return fmt.Errorf("rollback: chain %q already has keys; use overwrite=true to force", chainName)
	}

	// Restore snapshot into a temporary chain name then copy to target.
	tmpName := "__rollback_tmp_" + chainName
	if err := r.snapper.Restore(snapshotName, tmpName, true); err != nil {
		return fmt.Errorf("rollback: restore snapshot %q: %w", snapshotName, err)
	}

	restored, err := r.store.Get(tmpName)
	if err != nil {
		return fmt.Errorf("rollback: read restored chain: %w", err)
	}

	// Delete existing keys in target chain.
	if err := r.store.Delete(chainName); err != nil && !isNotFound(err) {
		_ = r.store.Delete(tmpName)
		return fmt.Errorf("rollback: clear target chain: %w", err)
	}

	for k, v := range restored {
		if err := r.store.Set(chainName, k, v); err != nil {
			_ = r.store.Delete(tmpName)
			return fmt.Errorf("rollback: set %s=%s: %w", k, v, err)
		}
	}

	_ = r.store.Delete(tmpName)
	return nil
}

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == chain.ErrChainNotFound.Error()
}
