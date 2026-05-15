package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/envchain-cli/internal/chain"
	"github.com/envchain-cli/internal/compare"
	"github.com/spf13/cobra"
)

func init() {
	var storeDir string

	cmd := &cobra.Command{
		Use:   "compare <chainA> <chainB>",
		Short: "Compare two chains and show key-level differences",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := chain.NewStore(storeDir)
			if err != nil {
				return err
			}
			res, err := compare.NewComparer(s).Compare(args[0], args[1])
			if err != nil {
				return err
			}

			w := cmd.OutOrStdout()

			if len(res.OnlyInA) > 0 {
				sort.Strings(res.OnlyInA)
				for _, k := range res.OnlyInA {
					fmt.Fprintf(w, "< %s\n", k)
				}
			}
			if len(res.OnlyInB) > 0 {
				sort.Strings(res.OnlyInB)
				for _, k := range res.OnlyInB {
					fmt.Fprintf(w, "> %s\n", k)
				}
			}
			keys := make([]string, 0, len(res.DiffValue))
			for k := range res.DiffValue {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				v := res.DiffValue[k]
				fmt.Fprintf(w, "~ %s: %q → %q\n", k, v[0], v[1])
			}

			if len(res.OnlyInA) == 0 && len(res.OnlyInB) == 0 && len(res.DiffValue) == 0 {
				fmt.Fprintln(w, "chains are identical")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&storeDir, "store", os.Getenv("ENVCHAIN_STORE"), "path to chain store")
	rootCmd.AddCommand(cmd)
}
