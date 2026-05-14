package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourorg/envchain-cli/internal/audit"
	"github.com/yourorg/envchain-cli/internal/config"
)

func init() {
	var chainFilter string
	var last int

	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Show the audit log of chain mutations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			logPath := cfg.AuditLogPath
			if logPath == "" {
				logPath = config.ConfigDir() + "/audit.log"
			}

			l, err := audit.NewLogger(logPath)
			if err != nil {
				return err
			}

			entries, err := l.ReadAll()
			if err != nil {
				return err
			}

			if chainFilter != "" {
				filtered := entries[:0]
				for _, e := range entries {
					if e.Chain == chainFilter {
						filtered = append(filtered, e)
					}
				}
				entries = filtered
			}

			if last > 0 && len(entries) > last {
				entries = entries[len(entries)-last:]
			}

			if len(entries) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no audit entries found")
				return nil
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIMESTAMP\tACTION\tCHAIN\tKEY\tEXTRA")
			for _, e := range entries {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					e.Timestamp.Format(time.RFC3339),
					e.Action, e.Chain, e.Key, e.Extra)
			}
			return w.Flush()
		},
	}

	auditCmd.Flags().StringVarP(&chainFilter, "chain", "c", "", "filter entries by chain name")
	auditCmd.Flags().IntVarP(&last, "last", "n", 0, "show only the last N entries")

	rootCmd.AddCommand(auditCmd)
}
