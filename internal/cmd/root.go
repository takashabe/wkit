package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"wkit/internal/worktree"
)

func NewRootCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "root",
		Short: "Show the root directory of the worktree",
		Long:  `Show the root directory of the Git repository.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			repoRoot, err := worktree.GetRepositoryRoot()
			if err != nil {
				return fmt.Errorf("failed to get repository root: %w", err)
			}

			// Output based on format
			if format == "json" {
				output := map[string]string{"root": repoRoot}
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				return encoder.Encode(output)
			}

			// Default human-readable format
			fmt.Fprintln(cmd.OutOrStdout(), repoRoot)
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "", "Output format (json)")

	return cmd
}
