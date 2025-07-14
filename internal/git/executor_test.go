package git

import (
	"testing"
)

func TestNewExecutor(t *testing.T) {
	executor := NewExecutor("/test/path")
	if executor == nil {
		t.Error("NewExecutor() returned nil")
	}
	if executor.workDir != "/test/path" {
		t.Errorf("Expected workDir '/test/path', got %s", executor.workDir)
	}
}

func TestNewExecutorEmptyWorkDir(t *testing.T) {
	executor := NewExecutor("")
	if executor == nil {
		t.Error("NewExecutor() returned nil")
	}
	if executor.workDir != "" {
		t.Errorf("Expected empty workDir, got %s", executor.workDir)
	}
}

// Note: The following tests would require a real git repository to work properly
// In a real testing environment, you might want to set up a temporary git repo
// or mock the exec.Command calls

func TestGetRepositoryRoot(t *testing.T) {
	// Skip this test if not in a git repository
	if !isInGitRepo() {
		t.Skip("Skipping TestGetRepositoryRoot: not in a git repository")
	}

	root, err := GetRepositoryRoot()
	if err != nil {
		t.Fatalf("GetRepositoryRoot() failed: %v", err)
	}

	if root == "" {
		t.Error("GetRepositoryRoot() returned empty string")
	}
}

// Helper function to check if we're in a git repository
func isInGitRepo() bool {
	executor := NewExecutor("")
	_, err := executor.Execute("rev-parse", "--git-dir")
	return err == nil
}

func TestBranchExists(t *testing.T) {
	// Skip this test if not in a git repository
	if !isInGitRepo() {
		t.Skip("Skipping TestBranchExists: not in a git repository")
	}

	executor := NewExecutor("")

	// Test with a branch that should exist (current branch)
	// Get current branch first
	currentBranch, err := executor.Execute("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		t.Fatalf("Failed to get current branch: %v", err)
	}

	// Skip test if we're in detached HEAD state (common in CI)
	if currentBranch != "HEAD" {
		exists := executor.BranchExists(currentBranch)
		if !exists {
			t.Errorf("BranchExists(%s) = false, expected true for current branch", currentBranch)
		}
	}

	// Test with a branch that should not exist
	nonExistentBranch := "this-branch-should-not-exist-12345"
	exists := executor.BranchExists(nonExistentBranch)
	if exists {
		t.Errorf("BranchExists(%s) = true, expected false for non-existent branch", nonExistentBranch)
	}
}
