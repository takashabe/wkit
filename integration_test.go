package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIIntegration(t *testing.T) {
	// Skip integration tests if not in a git repository
	if !isGitRepository() {
		t.Skip("Skipping integration tests: not in a git repository")
	}

	// Build the binary for testing
	binaryPath := buildBinary(t)
	defer os.Remove(binaryPath)

	t.Run("help command", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "--help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Help command failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "wkit is a CLI tool for convenient Git worktree management") {
			t.Errorf("Help output doesn't contain expected description")
		}
	})

	t.Run("list command", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "list")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("List command failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		// New format has header and matches git worktree list format: space-padded with [branch] format
		if !strings.Contains(outputStr, "PATH") || !strings.Contains(outputStr, "HEAD") || !strings.Contains(outputStr, "BRANCH") {
			t.Errorf("List output doesn't contain expected headers")
		}
		if !strings.Contains(outputStr, "(root)") || !strings.Contains(outputStr, " [") {
			t.Errorf("List output doesn't contain expected git worktree list format with (root) and [branch]")
		}
	})

	t.Run("config show command", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "config", "show")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Config show command failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Current configuration:") {
			t.Errorf("Config show output doesn't contain expected header")
		}
		if !strings.Contains(outputStr, "wkit_root:") {
			t.Errorf("Config show output doesn't contain wkit_root")
		}
	})

	t.Run("status command", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "status")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Status command failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "PATH") || !strings.Contains(outputStr, "BRANCH") || !strings.Contains(outputStr, "STATUS") {
			t.Errorf("Status output doesn't contain expected headers")
		}
	})
}

func TestInvalidCommand(t *testing.T) {
	// Build the binary for testing
	binaryPath := buildBinary(t)
	defer os.Remove(binaryPath)

	cmd := exec.Command(binaryPath, "nonexistent-command")
	output, err := cmd.CombinedOutput()
	
	// Command should fail
	if err == nil {
		t.Fatalf("Expected invalid command to fail, but it succeeded")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Error:") || !strings.Contains(outputStr, "unknown command") {
		t.Errorf("Invalid command output doesn't contain expected error message: %s", outputStr)
	}
}

func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func buildBinary(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "wkit-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	binaryPath := filepath.Join(tmpDir, "wkit")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	return binaryPath
}