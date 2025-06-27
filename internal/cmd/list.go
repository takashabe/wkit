package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"wkit/internal/worktree"
)

func NewListCmd() *cobra.Command {
	return &cobra.Command{
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

			fmt.Printf("%-30s %-20s %-12s\n", "PATH", "BRANCH", "HEAD")
			fmt.Println(strings.Repeat("-", 65))

			for _, wt := range worktrees {
				relativePath, err := filepath.Rel(repoRoot, wt.Path)
				if err != nil {
					relativePath = wt.Path // Fallback if relative path calculation fails
				}
				if relativePath == "." {
					relativePath = "(root)"
				}
				fmt.Printf("%-30s %-20s %-12s\n", relativePath, wt.Branch, wt.HEAD)
			}
			return nil
		},
	}
}