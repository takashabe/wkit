package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"wkit/internal/config"
	"wkit/internal/worktree"
)

func NewSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync [worktree]",
		Short: "Sync worktree with main branch",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := worktree.NewManager()
			if err != nil {
				return fmt.Errorf("failed to create manager: %w", err)
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			var targetWorktreePath string
			if len(args) > 0 {
				targetWorktreePath, err = manager.FindWorktreePath(args[0])
				if err != nil {
					return fmt.Errorf("failed to find worktree path: %w", err)
				}
			} else {
				currentDir, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
				worktrees, err := manager.ListWorktrees()
				if err != nil {
					return fmt.Errorf("failed to list worktrees: %w", err)
				}
				found := false
				for _, wt := range worktrees {
					if wt.Path == currentDir {
						targetWorktreePath = currentDir
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("current directory is not a worktree")
				}
			}

			rebaseFlag, _ := cmd.Flags().GetBool("rebase")
			useRebase := rebaseFlag || (cfg.DefaultSyncStrategy == "rebase")
			syncStrategy := "merge"
			if useRebase {
				syncStrategy = "rebase"
			}

			fmt.Printf("Syncing worktree '%s' with %s branch using %s...\n",
				targetWorktreePath, cfg.MainBranch, syncStrategy)

			err = manager.SyncWorktreeWithBranch(targetWorktreePath, cfg.MainBranch, useRebase)
			if err != nil {
				return fmt.Errorf("failed to sync worktree: %w", err)
			}

			fmt.Printf("✓ Successfully synced worktree '%s'\n", targetWorktreePath)
			return nil
		},
	}

	cmd.Flags().BoolP("rebase", "r", false, "Use rebase instead of merge")
	return cmd
}