package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"wkit/internal/config"
)

func NewConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management",
		Long:  `Manage wkit configuration.`,
	}

	configCmd.AddCommand(NewConfigShowCmd())
	configCmd.AddCommand(NewConfigSetCmd())
	configCmd.AddCommand(NewConfigInitCmd())

	return configCmd
}

func NewConfigShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			fmt.Println("Current configuration:")
			fmt.Printf("  wkit_root: %s\n", cfg.WkitRoot)
			fmt.Printf("  auto_cleanup: %t\n", cfg.AutoCleanup)
			fmt.Printf("  default_sync_strategy: %s\n", cfg.DefaultSyncStrategy)
			fmt.Printf("  main_branch: %s\n", cfg.MainBranch)
			fmt.Printf("  copy_files.enabled: %t\n", cfg.CopyFiles.Enabled)
			fmt.Printf("  copy_files.files: %v\n", cfg.CopyFiles.Files)
			return nil
		},
	}
}

func NewConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			switch key {
			case "wkit_root":
				cfg.WkitRoot = value
			case "auto_cleanup":
				b, err := parseBool(value)
				if err != nil {
					return fmt.Errorf("invalid boolean value for auto_cleanup: %w", err)
				}
				cfg.AutoCleanup = b
			case "default_sync_strategy":
				if value != "merge" && value != "rebase" {
					return fmt.Errorf("invalid sync strategy: %s. Valid values: merge, rebase", value)
				}
				cfg.DefaultSyncStrategy = value
			case "main_branch":
				cfg.MainBranch = value
			case "copy_files.enabled":
				b, err := parseBool(value)
				if err != nil {
					return fmt.Errorf("invalid boolean value for copy_files.enabled: %w", err)
				}
				cfg.CopyFiles.Enabled = b
			case "copy_files.files":
				cfg.CopyFiles.Files = strings.Split(value, ",")
			default:
				return fmt.Errorf("unknown configuration key: %s", key)
			}

			err = config.SaveGlobal(cfg)
			if err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Printf("✓ Configuration updated: %s = %s\n", key, value)
			return nil
		},
	}
}

func NewConfigInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a local configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.InitLocal()
			if err != nil {
				return fmt.Errorf("failed to create local config file: %w", err)
			}
			fmt.Println("✓ Created local configuration file: .wkit.toml")
			return nil
		},
	}
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