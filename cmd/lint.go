package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envchain-cli/internal/chain"
	"envchain-cli/internal/lint"
)

func init() {
	var chainName string

	lintCmd := &cobra.Command{
		Use:   "lint [chain]",
		Short: "Check chains for variable naming and value issues",
		Long: `Lint inspects environment variable chains and reports style or
correctness issues such as non-uppercase keys or empty values.

If a chain name is provided only that chain is checked; otherwise
all chains in the store are inspected.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				chainName = args[0]
			}

			st, err := chain.NewStore(storePath())
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}

			l := lint.NewLinter(st)
			issues, err := l.Lint(chainName)
			if err != nil {
				return err
			}

			if len(issues) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No issues found.")
				return nil
			}

			for _, iss := range issues {
				fmt.Fprintln(cmd.OutOrStdout(), iss)
			}

			// Exit with a non-zero code so CI pipelines can gate on lint.
			os.Exit(1)
			return nil
		},
	}

	rootCmd.AddCommand(lintCmd)
}
