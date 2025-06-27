package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	DefaultWorktreePath string   `mapstructure:"default_worktree_path"`
	AutoCleanup         bool     `mapstructure:"auto_cleanup"`
	ZIntegration        bool     `mapstructure:"z_integration"` // 削除予定だが、互換性のため残す
	DefaultSyncStrategy string   `mapstructure:"default_sync_strategy"`
	MainBranch          string   `mapstructure:"main_branch"`
	CopyFiles           CopyFiles `mapstructure:"copy_files"`
}

// CopyFiles represents the configuration for copying files
type CopyFiles struct {
	Enabled bool     `mapstructure:"enabled"`
	Files   []string `mapstructure:"files"`
}

// Load loads the configuration from local or global config files
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config") // global config file name
	v.SetConfigType("toml")
	if home, err := os.UserHomeDir(); err == nil {
		v.AddConfigPath(filepath.Join(home, ".config", "wkit")) // global config path
	}

	// Set default values
	v.SetDefault("default_worktree_path", ".git/.wkit-worktrees")
	v.SetDefault("auto_cleanup", false)
	v.SetDefault("z_integration", false)
	v.SetDefault("default_sync_strategy", "merge")
	v.SetDefault("main_branch", "main")
	v.SetDefault("copy_files.enabled", false)
	v.SetDefault("copy_files.files", []string{".envrc", "compose.override.yaml", ".env.local", "config/local.yaml"})

	// Read global config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read global config file: %w", err)
		}
	}

	// Read local config (if exists) and merge
	localV := viper.New()
	localV.SetConfigName(".wkit") // local config file name
	localV.SetConfigType("toml")
	localV.AddConfigPath(".")

	if err := localV.ReadInConfig(); err == nil {
		// Merge local config into global config
		for _, key := range localV.AllKeys() {
			v.Set(key, localV.Get(key))
		}
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read local config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// SaveGlobal saves the configuration to the global config file
func SaveGlobal(cfg *Config) error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	if home, err := os.UserHomeDir(); err == nil {
		v.AddConfigPath(filepath.Join(home, ".config", "wkit"))
	} else {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Ensure the config directory exists
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "wkit")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set values from the provided config struct
	v.Set("default_worktree_path", cfg.DefaultWorktreePath)
	v.Set("auto_cleanup", cfg.AutoCleanup)
	v.Set("z_integration", cfg.ZIntegration)
	v.Set("default_sync_strategy", cfg.DefaultSyncStrategy)
	v.Set("main_branch", cfg.MainBranch)
	v.Set("copy_files.enabled", cfg.CopyFiles.Enabled)
	v.Set("copy_files.files", cfg.CopyFiles.Files)

	configPath := filepath.Join(configDir, "config.toml")
	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write global config file: %w", err)
	}

	return nil
}

// InitLocal creates a local .wkit.toml file with default values
func InitLocal() error {
	v := viper.New()
	v.SetConfigName(".wkit")
	v.SetConfigType("toml")
	v.AddConfigPath(".")

	// Set default values
	v.SetDefault("default_worktree_path", ".git/.wkit-worktrees")
	v.SetDefault("auto_cleanup", false)
	v.SetDefault("z_integration", false)
	v.SetDefault("default_sync_strategy", "merge")
	v.SetDefault("main_branch", "main")
	v.SetDefault("copy_files.enabled", false)
	v.SetDefault("copy_files.files", []string{".envrc", "compose.override.yaml", ".env.local", "config/local.yaml"})

	if err := v.SafeWriteConfigAs(".wkit.toml"); err != nil {
		return fmt.Errorf("failed to create local config file: %w", err)
	}

	return nil
}

// ResolveWorktreePath resolves the worktree path based on config and provided path
func (c *Config) ResolveWorktreePath(branch string, providedPath string, repositoryRoot string) string {
	if providedPath != "" {
		return providedPath
	}

	if filepath.IsAbs(c.DefaultWorktreePath) {
		return filepath.Join(c.DefaultWorktreePath, branch)
	} else {
		return filepath.Join(repositoryRoot, c.DefaultWorktreePath, branch)
	}
}

// CopyFilesToWorktree copies configured files to the new worktree
func (c *Config) CopyFilesToWorktree(sourceDir string, targetDir string) ([]string, error) {
	if !c.CopyFiles.Enabled {
		return []string{}, nil
	}

	var copiedFiles []string

	for _, filePattern := range c.CopyFiles.Files {
		// Check if it's a relative path or just a filename
		if strings.Contains(filePattern, "/") || strings.Contains(filePattern, "\\") {
			sourceFile := filepath.Join(sourceDir, filePattern)
			targetFile := filepath.Join(targetDir, filePattern)

			if _, err := os.Stat(sourceFile); err == nil { // Check if source file exists
				if err := c.copySingleFile(sourceFile, targetFile, filePattern, &copiedFiles); err != nil {
					fmt.Fprintf(os.Stderr, "  Warning: Failed to copy %s: %v\n", filePattern, err)
				}
			}
		} else {
			// It's just a filename, search for all matching files in the repository
			foundFiles, err := findFilesByName(sourceDir, filePattern)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  Warning: Failed to find files for pattern %s: %v\n", filePattern, err)
				continue
			}
			for _, relativePath := range foundFiles {
				sourceFile := filepath.Join(sourceDir, relativePath)
				targetFile := filepath.Join(targetDir, relativePath)
				if err := c.copySingleFile(sourceFile, targetFile, relativePath, &copiedFiles); err != nil {
					fmt.Fprintf(os.Stderr, "  Warning: Failed to copy %s: %v\n", relativePath, err)
				}
			}
		}
	}

	return copiedFiles, nil
}

func (c *Config) copySingleFile(sourceFile string, targetFile string, relativePath string, copiedFiles *[]string) error {
	// Create parent directories if needed
	if err := os.MkdirAll(filepath.Dir(targetFile), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Skip if target file already exists
	if _, err := os.Stat(targetFile); err == nil {
		return nil
	}

	input, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file %s: %w", sourceFile, err)
	}

	err = os.WriteFile(targetFile, input, 0644)
	if err != nil {
		return fmt.Errorf("failed to write target file %s: %w", targetFile, err)
	}

	*copiedFiles = append(*copiedFiles, relativePath)
	return nil
}

func findFilesByName(baseDir string, filename string) ([]string, error) {
	var foundFiles []string
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir // Skip .git directory
		}
		if !info.IsDir() && info.Name() == filename {
			relativePath, err := filepath.Rel(baseDir, path)
			if err != nil {
				return err
			}
			foundFiles = append(foundFiles, relativePath)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", baseDir, err)
	}
	return foundFiles, nil
}
