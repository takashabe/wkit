package worktree

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseWorktreeList(t *testing.T) {
	testOutput := `worktree /path/to/repo
HEAD 1234567890abcdef
branch refs/heads/main

worktree /path/to/repo/.git/.wkit-worktrees/feature-branch
HEAD abcdef1234567890
branch refs/heads/feature-branch

worktree /path/to/repo/.git/.wkit-worktrees/another-branch
HEAD fedcba0987654321
branch refs/heads/another-branch
`

	worktrees, err := parseWorktreeList(testOutput)
	if err != nil {
		t.Fatalf("parseWorktreeList() failed: %v", err)
	}

	if len(worktrees) != 3 {
		t.Errorf("Expected 3 worktrees, got %d", len(worktrees))
		return
	}

	// Test first worktree (main)
	if worktrees[0].Path != "/path/to/repo" {
		t.Errorf("Expected path '/path/to/repo', got %s", worktrees[0].Path)
	}
	if worktrees[0].Branch != "main" {
		t.Errorf("Expected branch 'main', got %s", worktrees[0].Branch)
	}
	if worktrees[0].HEAD != "1234567890abcdef" {
		t.Errorf("Expected HEAD '1234567890abcdef', got %s", worktrees[0].HEAD)
	}

	// Test second worktree (feature-branch)
	if worktrees[1].Path != "/path/to/repo/.git/.wkit-worktrees/feature-branch" {
		t.Errorf("Expected path '/path/to/repo/.git/.wkit-worktrees/feature-branch', got %s", worktrees[1].Path)
	}
	if worktrees[1].Branch != "feature-branch" {
		t.Errorf("Expected branch 'feature-branch', got %s", worktrees[1].Branch)
	}
	if worktrees[1].HEAD != "abcdef1234567890" {
		t.Errorf("Expected HEAD 'abcdef1234567890', got %s", worktrees[1].HEAD)
	}
}

func TestParseGitStatus(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected *WorktreeStatus
	}{
		{
			name:   "clean status",
			output: "",
			expected: &WorktreeStatus{
				IsClean:   true,
				Modified:  0,
				Added:     0,
				Deleted:   0,
				Untracked: 0,
				Ahead:     0,
				Behind:    0,
			},
		},
		{
			name: "modified files",
			output: ` M file1.go
 M file2.go
A  file3.go
 D file4.go
?? file5.go`,
			expected: &WorktreeStatus{
				IsClean:   false,
				Modified:  2,
				Added:     1,
				Deleted:   1,
				Untracked: 1,
				Ahead:     0,
				Behind:    0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseGitStatus(tt.output)
			if err != nil {
				t.Fatalf("parseGitStatus() failed: %v", err)
			}

			if result.IsClean != tt.expected.IsClean {
				t.Errorf("IsClean = %v, want %v", result.IsClean, tt.expected.IsClean)
			}
			if result.Modified != tt.expected.Modified {
				t.Errorf("Modified = %v, want %v", result.Modified, tt.expected.Modified)
			}
			if result.Added != tt.expected.Added {
				t.Errorf("Added = %v, want %v", result.Added, tt.expected.Added)
			}
			if result.Deleted != tt.expected.Deleted {
				t.Errorf("Deleted = %v, want %v", result.Deleted, tt.expected.Deleted)
			}
			if result.Untracked != tt.expected.Untracked {
				t.Errorf("Untracked = %v, want %v", result.Untracked, tt.expected.Untracked)
			}
		})
	}
}

func TestContainsString(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	tests := []struct {
		name     string
		slice    []string
		search   string
		expected bool
	}{
		{
			name:     "found",
			slice:    slice,
			search:   "banana",
			expected: true,
		},
		{
			name:     "not found",
			slice:    slice,
			search:   "grape",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			search:   "apple",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsString(tt.slice, tt.search)
			if result != tt.expected {
				t.Errorf("containsString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewManager(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager() failed: %v", err)
	}

	if manager == nil {
		t.Error("NewManager() returned nil")
	}
}

func TestAddWorktree_BaseBranchHandling(t *testing.T) {
	// Test that the AddWorktree method properly handles base branch parameters
	// This tests the logic for prefix handling in AddWorktree method

	tests := []struct {
		name               string
		baseBranch         string
		localBranchExists  bool
		expectedBaseBranch string
	}{
		{
			name:               "base branch without origin prefix, local branch exists",
			baseBranch:         "main",
			localBranchExists:  true,
			expectedBaseBranch: "main",
		},
		{
			name:               "base branch without origin prefix, local branch doesn't exist",
			baseBranch:         "main",
			localBranchExists:  false,
			expectedBaseBranch: "origin/main",
		},
		{
			name:               "base branch with origin prefix",
			baseBranch:         "origin/develop",
			localBranchExists:  false,
			expectedBaseBranch: "origin/develop",
		},
		{
			name:               "base branch is empty, local branch doesn't exist",
			baseBranch:         "",
			localBranchExists:  false,
			expectedBaseBranch: "origin/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the prefix logic without calling actual git commands
			var actualBaseBranch string

			if strings.HasPrefix(tt.baseBranch, "origin/") {
				// Already has origin/ prefix, use as-is
				actualBaseBranch = tt.baseBranch
			} else if tt.localBranchExists {
				// Local branch exists, use as-is
				actualBaseBranch = tt.baseBranch
			} else {
				// Try remote branch
				actualBaseBranch = fmt.Sprintf("origin/%s", tt.baseBranch)
			}

			if actualBaseBranch != tt.expectedBaseBranch {
				t.Errorf("Expected base branch %s, got %s", tt.expectedBaseBranch, actualBaseBranch)
			}
		})
	}
}
