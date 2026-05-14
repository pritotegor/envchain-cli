// Package watch provides functionality to watch a chain for changes
// and re-execute a command when environment variables are updated.
package watch

import (
	"context"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/envchain-cli/internal/chain"
)

// Watcher polls a chain for changes and re-runs a command on update.
type Watcher struct {
	store    *chain.Store
	interval time.Duration
	logger   *log.Logger
}

// NewWatcher creates a Watcher that polls at the given interval.
func NewWatcher(store *chain.Store, interval time.Duration) *Watcher {
	return &Watcher{
		store:    store,
		interval: interval,
		logger:   log.New(os.Stderr, "[watch] ", log.LstdFlags),
	}
}

// Watch polls chainName for changes, re-running args when vars change.
// It blocks until ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context, chainName string, args []string) error {
	if len(args) == 0 {
		return ErrNoCommand
	}

	last, err := w.snapshot(chainName)
	if err != nil {
		return err
	}

	var proc *exec.Cmd
	proc = w.launch(ctx, args, last)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if proc != nil && proc.Process != nil {
				_ = proc.Process.Kill()
			}
			return ctx.Err()
		case <-ticker.C:
			current, err := w.snapshot(chainName)
			if err != nil {
				w.logger.Printf("poll error: %v", err)
				continue
			}
			if !equal(last, current) {
				w.logger.Printf("chain %q changed, restarting command", chainName)
				if proc != nil && proc.Process != nil {
					_ = proc.Process.Kill()
					_ = proc.Wait()
				}
				last = current
				proc = w.launch(ctx, args, current)
			}
		}
	}
}

func (w *Watcher) snapshot(chainName string) (map[string]string, error) {
	vars, err := w.store.Get(chainName)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = v
	}
	return out, nil
}

func (w *Watcher) launch(ctx context.Context, args []string, env map[string]string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = mergeEnv(os.Environ(), env)
	if err := cmd.Start(); err != nil {
		w.logger.Printf("failed to start command: %v", err)
		return nil
	}
	return cmd
}

func mergeEnv(base []string, overrides map[string]string) []string {
	out := make([]string, len(base))
	copy(out, base)
	for k, v := range overrides {
		out = append(out, k+"="+v)
	}
	return out
}

func equal(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
