package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/envchain-cli/internal/chain"
	"github.com/envchain-cli/internal/rename"
)

func init() {
	var overwrite bool

	renameCmd := &cobra.Command{
		Use:   "rename <src> <dst>",
		Short: "Rename an environment chain",
		Long: `Rename moves all variables from the source chain to the destination chain.

By default the command fails if the destination chain already exists.
Pass --overwrite to replace it.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			store, err := chain.NewStore(cfg.StorePath)
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}

			r := rename.NewRenamer(store)
			if err := r.Rename(cmd.Context(), rename.Options{
				Src:       args[0],
				Dst:       args[1],
				Overwrite: overwrite,
			}); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "renamed chain %q → %q\n", args[0], args[1])
			return nil
		},
	}

	renameCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite destination chain if it already exists")
	rootCmd.AddCommand(renameCmd)
}
