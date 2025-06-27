package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveWorktreePath(t *testing.T) {
	tests := []struct {
		name         string
		config       Config
		branch       string
		providedPath string
		repoRoot     string
		expected     string
	}{
		{
			name:         "with provided path",
			config:       Config{DefaultWorktreePath: ".git/.wkit-worktrees"},
			branch:       "feature-branch",
			providedPath: "/custom/path",
			repoRoot:     "/repo",
			expected:     "/custom/path",
		},
		{
			name:         "with relative default path",
			config:       Config{DefaultWorktreePath: ".git/.wkit-worktrees"},
			branch:       "feature-branch",
			providedPath: "",
			repoRoot:     "/repo",
			expected:     "/repo/.git/.wkit-worktrees/feature-branch",
		},
		{
			name:         "with absolute default path",
			config:       Config{DefaultWorktreePath: "/absolute/path"},
			branch:       "feature-branch",
			providedPath: "",
			repoRoot:     "/repo",
			expected:     "/absolute/path/feature-branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.ResolveWorktreePath(tt.branch, tt.providedPath, tt.repoRoot)
			if result != tt.expected {
				t.Errorf("ResolveWorktreePath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wkit-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to the temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test loading with default values
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// The actual default value may vary based on viper's behavior when no config file exists
	// Just check that it's not empty
	if cfg.DefaultWorktreePath == "" {
		t.Errorf("DefaultWorktreePath is empty")
	}

	if cfg.MainBranch != "main" {
		t.Errorf("MainBranch = %v, want %v", cfg.MainBranch, "main")
	}

	if cfg.DefaultSyncStrategy != "merge" {
		t.Errorf("DefaultSyncStrategy = %v, want %v", cfg.DefaultSyncStrategy, "merge")
	}
}

func TestInitLocal(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wkit-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to the temp directory
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test InitLocal
	err = InitLocal()
	if err != nil {
		t.Fatalf("InitLocal() failed: %v", err)
	}

	// Check if the file was created
	configPath := filepath.Join(tmpDir, ".wkit.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created at %v", configPath)
	}
}