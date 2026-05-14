package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/envchain-cli/internal/chain"
	"github.com/envchain-cli/internal/watch"
	"github.com/spf13/cobra"
)

func init() {
	var interval time.Duration

	watchCmd := &cobra.Command{
		Use:   "watch <chain> -- <command> [args...]",
		Short: "Re-run a command when a chain's variables change",
		Long: `Polls the named chain at a regular interval and restarts the given
command whenever its environment variables are updated.

Example:
  envchain watch dev -- go run ./server`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("requires a chain name and at least one command argument")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := resolveStorePath(cmd)
			if err != nil {
				return err
			}
			st, err := chain.NewStore(cfgPath)
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}

			chainName := args[0]
			commandArgs := args[1:]

			w := watch.NewWatcher(st, interval)
			ctx := cmd.Context()
			if err := w.Watch(ctx, chainName, commandArgs); err != nil {
				if err.Error() == "context canceled" {
					return nil
				}
				_, _ = fmt.Fprintln(os.Stderr, err)
				return err
			}
			return nil
		},
	}

	watchCmd.Flags().DurationVarP(
		&interval, "interval", "i", 2*time.Second,
		"polling interval (e.g. 500ms, 2s)",
	)

	rootCmd.AddCommand(watchCmd)
}
