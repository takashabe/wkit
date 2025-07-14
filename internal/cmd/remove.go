package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"wkit/internal/worktree"
)

func NewRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <worktree>",
		Short: "Remove a worktree",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			worktreeName := args[0]

			manager, err := worktree.NewManager()
			if err != nil {
				return fmt.Errorf("failed to create manager: %w", err)
			}

			worktreePath, err := manager.FindWorktreePath(worktreeName)
			if err != nil {
				return fmt.Errorf("failed to find worktree path: %w", err)
			}

			err = manager.RemoveWorktree(worktreePath)
			if err != nil {
				return fmt.Errorf("failed to remove worktree: %w", err)
			}

			fmt.Printf("✓ Removed worktree '%s'\n", worktreeName)
			return nil
		},
	}
}
