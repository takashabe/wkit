package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"wkit/internal/config"
)

func TestCopyFilesIntegrationInAddCommand(t *testing.T) {
	// Test copy_files functionality without actual git worktree creation
	// This tests the logic in add.go that calls CopyFilesToWorktree

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "wkit-copy-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	sourceDir := filepath.Join(tmpDir, "source")
	targetDir := filepath.Join(tmpDir, "target")

	// Create source and target directories
	err = os.MkdirAll(sourceDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	err = os.MkdirAll(targetDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	// Create test files in source directory
	err = os.WriteFile(filepath.Join(sourceDir, ".envrc"), []byte("export TEST_VAR=value"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create .envrc: %v", err)
	}

	err = os.MkdirAll(filepath.Join(sourceDir, "config"), 0o755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	err = os.WriteFile(filepath.Join(sourceDir, "config", "local.yaml"), []byte("env: test"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create config/local.yaml: %v", err)
	}

	tests := []struct {
		name          string
		copyEnabled   bool
		expectedFiles []string
	}{
		{
			name:          "copy_files enabled",
			copyEnabled:   true,
			expectedFiles: []string{".envrc", "config/local.yaml"},
		},
		{
			name:          "copy_files disabled",
			copyEnabled:   false,
			expectedFiles: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean target directory for each test
			os.RemoveAll(targetDir)
			err = os.MkdirAll(targetDir, 0o755)
			if err != nil {
				t.Fatalf("Failed to recreate target dir: %v", err)
			}

			// Create config
			cfg := &config.Config{
				CopyFiles: config.CopyFiles{
					Enabled: tt.copyEnabled,
					Files:   []string{".envrc", "config/local.yaml"},
				},
			}

			// Test the copy functionality directly (this is what add.go calls)
			copiedFiles, err := cfg.CopyFilesToWorktree(sourceDir, targetDir)
			if err != nil {
				t.Fatalf("CopyFilesToWorktree() failed: %v", err)
			}

			// Check the number of copied files
			if len(copiedFiles) != len(tt.expectedFiles) {
				t.Errorf("Expected %d copied files, got %d", len(tt.expectedFiles), len(copiedFiles))
				t.Errorf("Copied files: %v", copiedFiles)
				t.Errorf("Expected files: %v", tt.expectedFiles)
			}

			// Check if expected files were copied
			copiedMap := make(map[string]bool)
			for _, file := range copiedFiles {
				copiedMap[file] = true
			}

			for _, expectedFile := range tt.expectedFiles {
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
				sourceFile := filepath.Join(sourceDir, expectedFile)
				sourceContent, err := os.ReadFile(sourceFile)
				if err != nil {
					t.Fatalf("Failed to read source file %s: %v", sourceFile, err)
				}

				targetContent, err := os.ReadFile(targetFile)
				if err != nil {
					t.Fatalf("Failed to read target file %s: %v", targetFile, err)
				}

				if string(sourceContent) != string(targetContent) {
					t.Errorf("Content mismatch for file %s. Source: %q, Target: %q",
						expectedFile, string(sourceContent), string(targetContent))
				}
			}

			// When copy is disabled, ensure no files are copied
			if !tt.copyEnabled {
				for _, file := range []string{".envrc", "config/local.yaml"} {
					targetFile := filepath.Join(targetDir, file)
					if _, err := os.Stat(targetFile); err == nil {
						t.Errorf("File %s was copied when copy_files is disabled", file)
					}
				}
			}
		})
	}
}

func TestCopyDirectoryIntegration(t *testing.T) {
	// Test directory copying functionality
	tmpDir, err := os.MkdirTemp("", "wkit-copy-dir-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	sourceDir := filepath.Join(tmpDir, "source")
	targetDir := filepath.Join(tmpDir, "target")

	// Create source and target directories
	err = os.MkdirAll(sourceDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	err = os.MkdirAll(targetDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	// Create test directory structure
	testDir := filepath.Join(sourceDir, "scripts")
	err = os.MkdirAll(testDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create scripts directory: %v", err)
	}

	// Create files in test directory
	err = os.WriteFile(filepath.Join(testDir, "build.sh"), []byte("#!/bin/bash\necho 'Building...'"), 0o755)
	if err != nil {
		t.Fatalf("Failed to create build.sh: %v", err)
	}

	// Create subdirectory with files
	subDir := filepath.Join(testDir, "utils")
	err = os.MkdirAll(subDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create utils subdirectory: %v", err)
	}

	err = os.WriteFile(filepath.Join(subDir, "helper.py"), []byte("# Helper script"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create helper.py: %v", err)
	}

	tests := []struct {
		name          string
		copyEnabled   bool
		files         []string
		expectedFiles []string
	}{
		{
			name:        "copy directory enabled",
			copyEnabled: true,
			files:       []string{"scripts/"},
			expectedFiles: []string{
				"scripts/build.sh",
				"scripts/utils/helper.py",
			},
		},
		{
			name:          "copy directory disabled",
			copyEnabled:   false,
			files:         []string{"scripts/"},
			expectedFiles: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean target directory for each test
			os.RemoveAll(targetDir)
			err = os.MkdirAll(targetDir, 0o755)
			if err != nil {
				t.Fatalf("Failed to recreate target dir: %v", err)
			}

			// Create config
			cfg := &config.Config{
				CopyFiles: config.CopyFiles{
					Enabled: tt.copyEnabled,
					Files:   tt.files,
				},
			}

			// Test the copy functionality
			copiedFiles, err := cfg.CopyFilesToWorktree(sourceDir, targetDir)
			if err != nil {
				t.Fatalf("CopyFilesToWorktree() failed: %v", err)
			}

			// Check the number of copied files
			if len(copiedFiles) != len(tt.expectedFiles) {
				t.Errorf("Expected %d copied files, got %d", len(tt.expectedFiles), len(copiedFiles))
				t.Errorf("Copied files: %v", copiedFiles)
				t.Errorf("Expected files: %v", tt.expectedFiles)
			}

			// Check if expected files were copied
			copiedMap := make(map[string]bool)
			for _, file := range copiedFiles {
				copiedMap[file] = true
			}

			for _, expectedFile := range tt.expectedFiles {
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
				sourceFile := filepath.Join(sourceDir, expectedFile)
				sourceContent, err := os.ReadFile(sourceFile)
				if err != nil {
					t.Fatalf("Failed to read source file %s: %v", sourceFile, err)
				}

				targetContent, err := os.ReadFile(targetFile)
				if err != nil {
					t.Fatalf("Failed to read target file %s: %v", targetFile, err)
				}

				if string(sourceContent) != string(targetContent) {
					t.Errorf("Content mismatch for file %s. Source: %q, Target: %q",
						expectedFile, string(sourceContent), string(targetContent))
				}
			}

			// When copy is disabled, ensure no files are copied
			if !tt.copyEnabled {
				for _, expectedFile := range tt.expectedFiles {
					targetFile := filepath.Join(targetDir, expectedFile)
					if _, err := os.Stat(targetFile); err == nil {
						t.Errorf("File %s was copied when copy_files is disabled", expectedFile)
					}
				}
			}
		})
	}
}
