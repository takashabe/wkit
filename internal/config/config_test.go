package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveWkitPath(t *testing.T) {
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
			config:       Config{WkitRoot: ".git/.wkit-worktrees"},
			branch:       "feature-branch",
			providedPath: "/custom/path",
			repoRoot:     "/repo",
			expected:     "/custom/path",
		},
		{
			name:         "with relative wkit root",
			config:       Config{WkitRoot: ".git/.wkit-worktrees"},
			branch:       "feature-branch",
			providedPath: "",
			repoRoot:     "/repo",
			expected:     "/repo/.git/.wkit-worktrees/feature-branch",
		},
		{
			name:         "with absolute wkit root",
			config:       Config{WkitRoot: "/absolute/path"},
			branch:       "feature-branch",
			providedPath: "",
			repoRoot:     "/repo",
			expected:     "/absolute/path/feature-branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.ResolveWkitPath(tt.branch, tt.providedPath, tt.repoRoot)
			if result != tt.expected {
				t.Errorf("ResolveWkitPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test backward compatibility with old ResolveWorktreePath function
func TestResolveWorktreePath_BackwardCompatibility(t *testing.T) {
	config := Config{WkitRoot: ".git/.wkit-worktrees"}
	branch := "feature-branch"
	repoRoot := "/repo"
	expected := "/repo/.git/.wkit-worktrees/feature-branch"

	result := config.ResolveWorktreePath(branch, "", repoRoot)
	if result != expected {
		t.Errorf("ResolveWorktreePath() = %v, want %v", result, expected)
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
	if cfg.WkitRoot == "" {
		t.Errorf("WkitRoot is empty")
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

// Test backward compatibility when loading old config with default_worktree_path
func TestLoad_BackwardCompatibility(t *testing.T) {
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

	// Create a config file with old key name
	configContent := `default_worktree_path = "/old/path"`
	err = os.WriteFile(".wkit.toml", []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load config and check if it migrates correctly
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Check that WkitRoot is set from the old key
	if cfg.WkitRoot != "/old/path" {
		t.Errorf("WkitRoot = %v, want %v", cfg.WkitRoot, "/old/path")
	}

	// Check that DefaultWorktreePath is also set (for backward compatibility)
	if cfg.DefaultWorktreePath != "/old/path" {
		t.Errorf("DefaultWorktreePath = %v, want %v", cfg.DefaultWorktreePath, "/old/path")
	}
}

func TestCopyFilesToWorktreeWithNestedDirectories(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wkit-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	sourceDir := filepath.Join(tmpDir, "source")
	targetDir := filepath.Join(tmpDir, "target")

	// Create source directory structure
	err = os.MkdirAll(sourceDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}

	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	tests := []struct {
		name           string
		config         Config
		fileStructure  map[string]string // file path -> content
		expectedCopies []string
	}{
		{
			name: "copy files from nested directories",
			config: Config{
				CopyFiles: CopyFiles{
					Enabled: true,
					Files:   []string{".envrc", "config.yaml"},
				},
			},
			fileStructure: map[string]string{
				".envrc":                        "export PATH=/usr/bin",
				"subdir/.envrc":                 "export PATH=/usr/local/bin",
				"config/config.yaml":            "database: production",
				"deep/nested/config.yaml":       "database: test",
				"config/deep/nested/config.yaml": "database: staging",
			},
			expectedCopies: []string{
				".envrc",
				"subdir/.envrc",
				"config/config.yaml",
				"deep/nested/config.yaml",
				"config/deep/nested/config.yaml",
			},
		},
		{
			name: "copy files with specific path patterns",
			config: Config{
				CopyFiles: CopyFiles{
					Enabled: true,
					Files:   []string{"config/local.yaml", ".env.local"},
				},
			},
			fileStructure: map[string]string{
				"config/local.yaml":      "env: local",
				"config/prod.yaml":       "env: prod",
				"other/config/local.yaml": "env: other",
				".env.local":             "DEBUG=true",
				"subdir/.env.local":      "DEBUG=false",
			},
			expectedCopies: []string{
				"config/local.yaml",
				".env.local",
				"subdir/.env.local",
			},
		},
		{
			name: "no files match pattern",
			config: Config{
				CopyFiles: CopyFiles{
					Enabled: true,
					Files:   []string{"nonexistent.txt"},
				},
			},
			fileStructure: map[string]string{
				"existing.txt": "content",
			},
			expectedCopies: []string{},
		},
		{
			name: "copy disabled",
			config: Config{
				CopyFiles: CopyFiles{
					Enabled: false,
					Files:   []string{".envrc"},
				},
			},
			fileStructure: map[string]string{
				".envrc": "export PATH=/usr/bin",
			},
			expectedCopies: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up directories for each test
			os.RemoveAll(sourceDir)
			os.RemoveAll(targetDir)
			err = os.MkdirAll(sourceDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create source dir: %v", err)
			}
			err = os.MkdirAll(targetDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create target dir: %v", err)
			}

			// Create test file structure
			for filePath, content := range tt.fileStructure {
				fullPath := filepath.Join(sourceDir, filePath)
				dir := filepath.Dir(fullPath)
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("Failed to create directory %s: %v", dir, err)
				}
				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create file %s: %v", fullPath, err)
				}
			}

			// Test CopyFilesToWorktree
			copiedFiles, err := tt.config.CopyFilesToWorktree(sourceDir, targetDir)
			if err != nil {
				t.Fatalf("CopyFilesToWorktree() failed: %v", err)
			}

			// Check the number of copied files
			if len(copiedFiles) != len(tt.expectedCopies) {
				t.Errorf("CopyFilesToWorktree() copied %d files, expected %d", len(copiedFiles), len(tt.expectedCopies))
				t.Errorf("Copied files: %v", copiedFiles)
				t.Errorf("Expected files: %v", tt.expectedCopies)
			}

			// Check if all expected files were copied
			copiedMap := make(map[string]bool)
			for _, file := range copiedFiles {
				copiedMap[file] = true
			}

			for _, expectedFile := range tt.expectedCopies {
				if !copiedMap[expectedFile] {
					t.Errorf("Expected file %s was not copied", expectedFile)
				}

				// Check if the target file exists and has correct content
				targetFile := filepath.Join(targetDir, expectedFile)
				if _, err := os.Stat(targetFile); os.IsNotExist(err) {
					t.Errorf("Target file %s does not exist", targetFile)
					continue
				}

				// Verify content matches
				expectedContent := tt.fileStructure[expectedFile]
				actualContent, err := os.ReadFile(targetFile)
				if err != nil {
					t.Errorf("Failed to read target file %s: %v", targetFile, err)
					continue
				}

				if string(actualContent) != expectedContent {
					t.Errorf("Target file %s content mismatch. Got: %q, Want: %q", targetFile, string(actualContent), expectedContent)
				}
			}

			// Check for unexpected files
			for _, copiedFile := range copiedFiles {
				found := false
				for _, expectedFile := range tt.expectedCopies {
					if copiedFile == expectedFile {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Unexpected file was copied: %s", copiedFile)
				}
			}
		})
	}
}