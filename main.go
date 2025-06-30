package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"wkit/internal/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "wkit",
	Short: "A Git worktree management toolkit",
	Long:  `wkit is a CLI tool for convenient Git worktree management.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(cmd.NewListCmd())
	rootCmd.AddCommand(cmd.NewAddCmd())
	rootCmd.AddCommand(cmd.NewRemoveCmd())
	rootCmd.AddCommand(cmd.NewSwitchCmd())
	rootCmd.AddCommand(cmd.NewConfigCmd())
	rootCmd.AddCommand(cmd.NewStatusCmd())
	rootCmd.AddCommand(cmd.NewCleanCmd())
	rootCmd.AddCommand(cmd.NewSyncCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
