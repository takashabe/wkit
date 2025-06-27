package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"wkit/internal/config"
	"wkit/internal/worktree"
)

var rootCmd = &cobra.Command{
	Use:   "wkit",
	Short: "A Git worktree management toolkit",
	Long:  `wkit is a CLI tool for convenient Git worktree management.`,
	Run: func(cmd *cobra.Command, args []string) {
		// デフォルトの動作、ヘルプを表示するなど
		cmd.Help()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all worktrees",
	Long:  `List all Git worktrees associated with the current repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := worktree.NewManager() // manager の初期化をここへ移動
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		worktrees, err := manager.ListWorktrees()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		repoRoot, err := worktree.GetRepositoryRoot() // ここを修正
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
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
	},
}

var addCmd = &cobra.Command{
	Use:   "add <branch> [path]",
	Short: "Add a new worktree",
	Args:  cobra.RangeArgs(1, 2), // branch (required), path (optional)
	Run: func(cmd *cobra.Command, args []string) {
		branch := args[0]
		var worktreePath string

		manager, err := worktree.NewManager() // manager の初期化をここへ移動
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(args) > 1 {
			worktreePath = args[1]
		} else {
			// デフォルトのパスを生成 (Rust版の default_worktree_path を考慮)
			// ここでは仮に .git/.wkit-worktrees/<branch> とする
			// 実際には設定ファイルから読み込む必要がある
			repoRoot, err := worktree.GetRepositoryRoot() // ここを修正
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			worktreePath = filepath.Join(repoRoot, ".git", ".wkit-worktrees", branch)
		}

		noSwitch, _ := cmd.Flags().GetBool("no-switch")

		// mainBranch は設定から取得する必要があるが、ここでは仮に "main" とする
		err = manager.AddWorktree(branch, worktreePath, "main")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Created worktree for branch '%s' at '%s'\n", branch, worktreePath)

		if !noSwitch {
			fmt.Println(worktreePath)
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove <worktree>",
	Short: "Remove a worktree",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		worktreeName := args[0]

		manager, err := worktree.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// ワークツリー名からパスを特定する必要がある
		// Rust版では find_worktree_by_name を使っていた
		// ここでは一旦、worktreeName がそのままパスとして使えると仮定する
		// TODO: ワークツリー名からパスを解決するロジックを追加
		err = manager.RemoveWorktree(worktreeName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Removed worktree '%s'\n", worktreeName)
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch <worktree>",
	Short: "Switch to a worktree",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		worktreeName := args[0]

		manager, err := worktree.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		worktreePath, err := manager.FindWorktreePath(worktreeName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(worktreePath)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
	Long:  `Manage wkit configuration.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Current configuration:")
		fmt.Printf("  default_worktree_path: %s\n", cfg.DefaultWorktreePath)
		fmt.Printf("  auto_cleanup: %t\n", cfg.AutoCleanup)
		fmt.Printf("  default_sync_strategy: %s\n", cfg.DefaultSyncStrategy)
		fmt.Printf("  main_branch: %s\n", cfg.MainBranch)
		fmt.Printf("  copy_files.enabled: %t\n", cfg.CopyFiles.Enabled)
		fmt.Printf("  copy_files.files: %v\n", cfg.CopyFiles.Files)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		switch key {
		case "default_worktree_path":
			cfg.DefaultWorktreePath = value
		case "auto_cleanup":
			b, err := parseBool(value)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid boolean value for auto_cleanup: %v\n", err)
				os.Exit(1)
			}
			cfg.AutoCleanup = b
		case "default_sync_strategy":
			if value != "merge" && value != "rebase" {
				fmt.Fprintf(os.Stderr, "Error: Invalid sync strategy: %s. Valid values: merge, rebase\n", value)
				os.Exit(1)
			}
			cfg.DefaultSyncStrategy = value
		case "main_branch":
			cfg.MainBranch = value
		case "copy_files.enabled":
			b, err := parseBool(value)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid boolean value for copy_files.enabled: %v\n", err)
				os.Exit(1)
			}
			cfg.CopyFiles.Enabled = b
		case "copy_files.files":
			cfg.CopyFiles.Files = strings.Split(value, ",")
		default:
			fmt.Fprintf(os.Stderr, "Error: Unknown configuration key: %s\n", key)
			os.Exit(1)
		}

		err = config.SaveGlobal(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Configuration updated: %s = %s\n", key, value)
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a local configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.InitLocal()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating local config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Created local configuration file: .wkit.toml")
	},
}

func parseBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "t", "1":
		return true, nil
	case "false", "f", "0":
		return false, nil
	}
	return false, fmt.Errorf("invalid boolean value: %s", s)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show git status of all worktrees",
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := worktree.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		worktrees, err := manager.ListWorktrees()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		repoRoot, err := worktree.GetRepositoryRoot()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
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
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up unnecessary worktrees",
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

		manager, err := worktree.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		unnecessaryWorktrees, err := manager.FindUnnecessaryWorktrees(cfg.MainBranch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding unnecessary worktrees: %v\n", err)
			os.Exit(1)
		}

		if len(unnecessaryWorktrees) == 0 {
			fmt.Println("No unnecessary worktrees found.")
			return
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
				return
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
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync [worktree]",
	Short: "Sync worktree with main branch",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := worktree.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		var targetWorktreePath string
		if len(args) > 0 {
			// worktreeName が指定された場合、そのパスを解決
			targetWorktreePath, err = manager.FindWorktreePath(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		} else {
			// 指定がない場合、現在のディレクトリをワークツリーとして使用
			currentDir, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
				os.Exit(1)
			}
			// 現在のディレクトリがワークツリーであることを確認
			worktrees, err := manager.ListWorktrees()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing worktrees: %v\n", err)
				os.Exit(1)
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
				fmt.Fprintf(os.Stderr, "Error: Current directory is not a worktree\n")
				os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "Error syncing worktree: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Successfully synced worktree '%s'\n", targetWorktreePath)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(syncCmd) // syncCmd を追加

	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configInitCmd)

	addCmd.Flags().Bool("no-switch", false, "Skip automatic switching to new worktree")
	cleanCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	syncCmd.Flags().BoolP("rebase", "r", false, "Use rebase instead of merge")
	// 他のコマンドもここに追加
}


func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
