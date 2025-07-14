package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"wkit/internal/worktree"
)

func NewListCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all worktrees",
		Long:  `List all Git worktrees associated with the current repository.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := worktree.NewManager()
			if err != nil {
				return fmt.Errorf("failed to create manager: %w", err)
			}

			worktrees, err := manager.ListWorktrees()
			if err != nil {
				return fmt.Errorf("failed to list worktrees: %w", err)
			}

			repoRoot, err := worktree.GetRepositoryRoot()
			if err != nil {
				return fmt.Errorf("failed to get repository root: %w", err)
			}

			// Convert absolute paths to relative paths for output
			type outputWorktree struct {
				Path   string `json:"path"`
				Branch string `json:"branch"`
				HEAD   string `json:"head"`
			}

			outputWorktrees := make([]outputWorktree, 0, len(worktrees))
			for _, wt := range worktrees {
				var relativePath string
				if wt.Path == repoRoot {
					relativePath = "(root)"
				} else {
					var err error
					relativePath, err = filepath.Rel(repoRoot, wt.Path)
					if err != nil {
						relativePath = wt.Path // Fallback if relative path calculation fails
					}
				}
				outputWorktrees = append(outputWorktrees, outputWorktree{
					Path:   relativePath,
					Branch: wt.Branch,
					HEAD:   wt.HEAD,
				})
			}

			// Output based on format
			if format == "json" {
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				return encoder.Encode(outputWorktrees)
			}

			// Default human-readable format using tabwriter for proper alignment
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 4, 2, ' ', 0)
			defer w.Flush()

			// Header
			fmt.Fprintln(w, "PATH\tHEAD\tBRANCH")
			fmt.Fprintln(w, "----\t----\t------")

			for _, wt := range outputWorktrees {
				// Truncate HEAD to 7 characters for display
				displayHEAD := wt.HEAD
				if len(displayHEAD) > 7 {
					displayHEAD = displayHEAD[:7]
				}
				// Format with tabs for proper alignment
				fmt.Fprintf(w, "%s\t%s\t%s\n", wt.Path, displayHEAD, wt.Branch)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "", "Output format (json)")

	return cmd
}
