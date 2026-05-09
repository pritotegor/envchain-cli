package cmd

import (
	"fmt"
	"os"

	"github.com/envchain-cli/internal/chain"
	"github.com/envchain-cli/internal/config"
	"github.com/envchain-cli/internal/merge"
	"github.com/spf13/cobra"
)

var mergeOverwrite bool

var mergeCmd = &cobra.Command{
	Use:   "merge <destination> <source>...",
	Short: "Merge variables from one or more source chains into a destination chain",
	Long: `Merge copies all variables from each source chain into the destination chain.

Sources are processed left-to-right. Use --overwrite to allow source values to
replace existing keys in the destination chain.`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dst := args[0]
		sources := args[1:]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		st, err := chain.NewStore(cfg.StorePath)
		if err != nil {
			return fmt.Errorf("open store: %w", err)
		}

		m := merge.NewMerger(st)
		if err := m.Merge(dst, sources, merge.MergeOptions{Overwrite: mergeOverwrite}); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			return err
		}

		fmt.Printf("Merged %v → %s\n", sources, dst)
		return nil
	},
}

func init() {
	mergeCmd.Flags().BoolVarP(&mergeOverwrite, "overwrite", "f", false,
		"overwrite existing keys in the destination chain")
	rootCmd.AddCommand(mergeCmd)
}
