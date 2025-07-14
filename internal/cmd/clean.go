package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"wkit/internal/config"
	"wkit/internal/worktree"
)

func NewCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean up unnecessary worktrees",
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")

			manager, err := worktree.NewManager()
			if err != nil {
				return fmt.Errorf("failed to create manager: %w", err)
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			unnecessaryWorktrees, err := manager.FindUnnecessaryWorktrees(cfg.MainBranch)
			if err != nil {
				return fmt.Errorf("failed to find unnecessary worktrees: %w", err)
			}

			if len(unnecessaryWorktrees) == 0 {
				fmt.Println("No unnecessary worktrees found.")
				return nil
			}

			fmt.Printf("Found %d unnecessary worktree(s):\n", len(unnecessaryWorktrees))
			for _, uw := range unnecessaryWorktrees {
				fmt.Printf("  %s - %s\n", uw.Worktree.Path, uw.Reason)
			}

			if !force {
				fmt.Print("\nRemove these worktrees? (y/N): ")
				var confirm string
				fmt.Scanln(&confirm)
				if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
					fmt.Println("Cancelled.")
					return nil
				}
			}

			for _, uw := range unnecessaryWorktrees {
				err := manager.RemoveWorktree(uw.Worktree.Path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error removing worktree %s: %v\n", uw.Worktree.Path, err)
					continue
				}
				fmt.Printf("✓ Removed worktree at '%s'\n", uw.Worktree.Path)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	return cmd
}
