package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"wkit/internal/config"
	"wkit/internal/worktree"
)

func NewAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <branch> [path]",
		Short: "Add a new worktree",
		Args:  cobra.RangeArgs(1, 2), // branch (required), path (optional)
		RunE: func(cmd *cobra.Command, args []string) error {
			branch := args[0]
			var worktreePath string

			manager, err := worktree.NewManager()
			if err != nil {
				return fmt.Errorf("failed to create manager: %w", err)
			}

			if len(args) > 1 {
				worktreePath = args[1]
			} else {
				cfg, err := config.Load()
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}
				repoRoot, err := worktree.GetRepositoryRoot()
				if err != nil {
					return fmt.Errorf("failed to get repository root: %w", err)
				}
				worktreePath = cfg.ResolveWorktreePath(branch, "", repoRoot)
			}

			noSwitch, _ := cmd.Flags().GetBool("no-switch")
			baseBranch, _ := cmd.Flags().GetString("base-branch")

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Use provided base branch or fall back to config main branch
			if baseBranch == "" {
				baseBranch = cfg.MainBranch
			}

			err = manager.AddWorktree(branch, worktreePath, baseBranch)
			if err != nil {
				return fmt.Errorf("failed to add worktree: %w", err)
			}

			fmt.Printf("✓ Created worktree for branch '%s' at '%s'\n", branch, worktreePath)

			// Copy configured files if enabled
			repoRoot, err := worktree.GetRepositoryRoot()
			if err != nil {
				return fmt.Errorf("failed to get repository root: %w", err)
			}

			copiedFiles, err := cfg.CopyFilesToWorktree(repoRoot, worktreePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to copy files: %v\n", err)
			} else if len(copiedFiles) > 0 {
				fmt.Printf("✓ Copied files: %v\n", copiedFiles)
			}

			if !noSwitch {
				relativePath, err := worktree.GetRelativePathFromRoot()
				if err != nil {
					// If we can't get relative path, just output the worktree path
					fmt.Println(worktreePath)
				} else if relativePath != "" {
					fmt.Printf("%s:%s\n", worktreePath, relativePath)
				} else {
					fmt.Println(worktreePath)
				}
			}
			return nil
		},
	}

	cmd.Flags().Bool("no-switch", false, "Skip automatic switching to new worktree")
	cmd.Flags().StringP("base-branch", "b", "", "Base branch to create new branch from (defaults to config main_branch)")
	return cmd
}
