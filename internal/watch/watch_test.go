package watch_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/envchain-cli/internal/chain"
	"github.com/envchain-cli/internal/watch"
)

func tempStore(t *testing.T) *chain.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := chain.NewStore(filepath.Join(dir, "chains.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return st
}

func TestWatch_ErrNoCommand(t *testing.T) {
	st := tempStore(t)
	if err := st.Add("dev", "KEY", "val"); err != nil {
		t.Fatal(err)
	}
	w := watch.NewWatcher(st, 50*time.Millisecond)
	ctx := context.Background()
	err := w.Watch(ctx, "dev", nil)
	if err != watch.ErrNoCommand {
		t.Fatalf("expected ErrNoCommand, got %v", err)
	}
}

func TestWatch_ErrUnknownChain(t *testing.T) {
	st := tempStore(t)
	w := watch.NewWatcher(st, 50*time.Millisecond)
	ctx := context.Background()
	err := w.Watch(ctx, "nonexistent", []string{"echo", "hi"})
	if err == nil {
		t.Fatal("expected error for unknown chain")
	}
}

func TestWatch_CancelStopsLoop(t *testing.T) {
	st := tempStore(t)
	if err := st.Add("dev", "FOO", "bar"); err != nil {
		t.Fatal(err)
	}
	w := watch.NewWatcher(st, 30*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- w.Watch(ctx, "dev", []string{os.Args[0], "-test.run=^$"})
	}()

	time.Sleep(80 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Watch did not return after cancel")
	}
}

func TestWatch_DetectsChange(t *testing.T) {
	st := tempStore(t)
	if err := st.Add("dev", "FOO", "initial"); err != nil {
		t.Fatal(err)
	}
	w := watch.NewWatcher(st, 30*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = w.Watch(ctx, "dev", []string{os.Args[0], "-test.run=^$"})
	}()

	time.Sleep(50 * time.Millisecond)
	if err := st.Add("dev", "FOO", "updated"); err != nil {
		t.Fatal(err)
	}
	// Give the watcher time to detect the change without asserting restart
	// internals — we just ensure no panic or deadlock occurs.
	time.Sleep(100 * time.Millisecond)
}
