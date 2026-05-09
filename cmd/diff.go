package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envchain-cli/internal/chain"
	"envchain-cli/internal/diff"
)

func init() {
	var flagColor bool

	diffCmd := &cobra.Command{
		Use:   "diff <source> <destination>",
		Short: "Show differences between two environment chains",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			store, err := chain.NewStore(cfg.StorePath)
			if err != nil {
				return err
			}

			d := diff.NewDiffer(store)
			result, err := d.Diff(args[0], args[1])
			if err != nil {
				return err
			}

			changes := result.Changes()
			if len(changes) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No differences found.")
				return nil
			}

			for _, e := range changes {
				line := e.String()
				if flagColor {
					line = e.Colored()
				}
				fmt.Fprintln(cmd.OutOrStdout(), line)
			}
			return nil
		},
	}

	diffCmd.Flags().BoolVar(&flagColor, "color", false, "Colorize output")
	rootCmd.AddCommand(diffCmd)

	_ = os.Stderr // keep import
}
