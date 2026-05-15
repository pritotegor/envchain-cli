package cmd

import (
	"fmt"
	"os"

	"github.com/envchain-cli/internal/chain"
	"github.com/envchain-cli/internal/rollback"
	"github.com/spf13/cobra"
)

func init() {
	var overwrite bool

	cmd := &cobra.Command{
		Use:   "rollback <chain> <snapshot>",
		Short: "Restore a chain to a previously captured snapshot",
		Long: `Rollback replaces the current contents of <chain> with the
environment variables recorded in <snapshot>.

By default the command refuses to overwrite a chain that already contains
keys. Pass --overwrite to force the replacement.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainName := args[0]
			snapshotName := args[1]

			store, err := chain.NewStore(chainDir())
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}

			rb := rollback.NewRollbacker(store)
			if err := rb.Rollback(chainName, snapshotName, overwrite); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return err
			}

			fmt.Printf("Chain %q restored from snapshot %q\n", chainName, snapshotName)
			return nil
		},
	}

	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys in the target chain")
	rootCmd.AddCommand(cmd)
}
