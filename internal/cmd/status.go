package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"wkit/internal/worktree"
)

func NewStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show git status of all worktrees",
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

			fmt.Printf("%-30s %-20s %-12s %-15s\n", "PATH", "BRANCH", "HEAD", "STATUS")
			fmt.Println(strings.Repeat("-", 80))

			for _, wt := range worktrees {
				relativePath, err := filepath.Rel(repoRoot, wt.Path)
				if err != nil {
					relativePath = wt.Path // Fallback if relative path calculation fails
				}
				if relativePath == "." {
					relativePath = "(root)"
				}

				status, err := manager.GetWorktreeStatus(wt.Path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting status for %s: %v\n", relativePath, err)
					continue
				}

				statusStr := ""
				if status.IsClean {
					statusStr = "Clean"
				} else {
					statusStr = fmt.Sprintf("%dM %dA %dD", status.Modified, status.Added, status.Deleted)
				}

				fmt.Printf("%-30s %-20s %-12s %-15s\n",
					relativePath,
					wt.Branch,
					wt.HEAD,
					statusStr,
				)

				if !status.IsClean {
					if status.Modified > 0 {
						fmt.Printf("  📝 %d modified files\n", status.Modified)
					}
					if status.Added > 0 {
						fmt.Printf("  ➕ %d added files\n", status.Added)
					}
					if status.Deleted > 0 {
						fmt.Printf("  ❌ %d deleted files\n", status.Deleted)
					}
					if status.Untracked > 0 {
						fmt.Printf("  ❓ %d untracked files\n", status.Untracked)
					}
				}
			}
			return nil
		},
	}
}
