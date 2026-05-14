package cmd

import (
	"fmt"
	"os"

	"github.com/envchain-cli/envchain-cli/internal/chain"
	"github.com/envchain-cli/envchain-cli/internal/clone"
	"github.com/envchain-cli/envchain-cli/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	var overwrite bool
	var prefix string

	cloneCmd := &cobra.Command{
		Use:   "clone <source> <destination>",
		Short: "Clone a chain into a new chain",
		Long: `Clone copies all key-value pairs from SOURCE into a new chain named
DESTINATION. Use --prefix to restrict which keys are copied.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			s, err := chain.NewStore(cfg.StorePath)
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}

			c := clone.NewCloner(s)
			err = c.Clone(src, dst, clone.CloneOptions{
				Overwrite: overwrite,
				Prefix:    prefix,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "cloned chain %q → %q\n", src, dst)
			return nil
		},
	}

	cloneCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite destination chain if it already exists")
	cloneCmd.Flags().StringVar(&prefix, "prefix", "", "only clone keys that start with this prefix")

	rootCmd.AddCommand(cloneCmd)
}
