package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"wkit/internal/worktree"
)

func NewSwitchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "switch <worktree>",
		Short: "Switch to a worktree",
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

			fmt.Println(worktreePath)
			return nil
		},
	}
}