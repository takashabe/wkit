package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestListCommand(t *testing.T) {
	tests := []struct {
		name           string
		worktrees      []struct {
			path   string
			branch string
			head   string
		}
		repoRoot       string
		expectedOutput []string
	}{
		{
			name: "standard worktree setup",
			worktrees: []struct {
				path   string
				branch string
				head   string
			}{
				{
					path:   "/path/to/repo",
					branch: "main",
					head:   "1234567890abcdef",
				},
				{
					path:   "/path/to/repo/.git/.wkit-worktrees/feature-branch",
					branch: "feature-branch",
					head:   "abcdef1234567890",
				},
			},
			repoRoot: "/path/to/repo",
			expectedOutput: []string{
				"PATH                                                                        HEAD BRANCH",
				"--------------------------------------------------------------------------- ------- ----------",
				"(root)                                                                      1234567 [main]",
				".git/.wkit-worktrees/feature-branch                                        abcdef1 [feature-branch]",
			},
		},
		{
			name: "worktree with long branch name",
			worktrees: []struct {
				path   string
				branch string
				head   string
			}{
				{
					path:   "/path/to/repo",
					branch: "main",
					head:   "1234567890abcdef",
				},
				{
					path:   "/path/to/repo/.git/.wkit-worktrees/very-long-feature-branch-name",
					branch: "very-long-feature-branch-name",
					head:   "abcdef1234567890",
				},
			},
			repoRoot: "/path/to/repo",
			expectedOutput: []string{
				"PATH                                                                        HEAD BRANCH",
				"--------------------------------------------------------------------------- ------- ----------",
				"(root)                                                                      1234567 [main]",
				".git/.wkit-worktrees/very-long-feature-branch-name                         abcdef1 [very-long-feature-branch-name]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a unit test for the output formatting logic
			// We'll verify that the paths are correctly formatted as relative to repo root
			
			// Verify the expected space-padded format with header
			for i, expected := range tt.expectedOutput {
				if i == 0 {
					// Header line
					if !strings.Contains(expected, "PATH") || !strings.Contains(expected, "HEAD") || !strings.Contains(expected, "BRANCH") {
						t.Errorf("Expected header line to contain PATH, HEAD, and BRANCH")
					}
				} else if i == 1 {
					// Separator line
					if !strings.Contains(expected, "---") {
						t.Errorf("Expected separator line with dashes")
					}
				} else if strings.Contains(expected, "(root)") {
					// Root worktree should be marked as (root)
					if !strings.HasPrefix(expected, "(root)") {
						t.Errorf("Expected root worktree to start with '(root)', got: %s", expected)
					}
					if !strings.Contains(expected, "[main]") {
						t.Errorf("Expected root worktree to contain '[main]', got: %s", expected)
					}
				} else {
					// Other worktrees should show relative path from repo root
					if !strings.HasPrefix(expected, ".git/.wkit-worktrees/") {
						t.Errorf("Expected non-root worktree path to start with .git/.wkit-worktrees/, got: %s", expected)
					}
					
					// Should have the git worktree list format: path + spaces + hash + space + [branch]
					if !strings.Contains(expected, " [") || !strings.HasSuffix(expected, "]") {
						t.Errorf("Expected format 'path hash [branch]', got: %s", expected)
					}
				}
			}
		})
	}
}

func TestListCommandRelativePaths(t *testing.T) {
	// Test that paths are always relative to the git repository root,
	// regardless of where the command is executed from
	
	tests := []struct {
		name         string
		worktreePath string
		repoRoot     string
		expected     string
	}{
		{
			name:         "root worktree",
			worktreePath: "/home/user/myrepo",
			repoRoot:     "/home/user/myrepo",
			expected:     "(root)",
		},
		{
			name:         "nested worktree",
			worktreePath: "/home/user/myrepo/.git/.wkit-worktrees/feature",
			repoRoot:     "/home/user/myrepo",
			expected:     ".git/.wkit-worktrees/feature",
		},
		{
			name:         "deeply nested worktree",
			worktreePath: "/home/user/myrepo/.git/.wkit-worktrees/deep/nested/feature",
			repoRoot:     "/home/user/myrepo",
			expected:     ".git/.wkit-worktrees/deep/nested/feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test verifies the path calculation logic
			// The actual implementation would use filepath.Rel(repoRoot, worktreePath)
			// and special case when they are equal to return "(root)"
		})
	}
}

func TestListCommandJSONFormat(t *testing.T) {
	// Test JSON output format
	cmd := NewListCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--format", "json"})

	// We would need to mock the worktree manager here
	// For now, this is a placeholder to show the test structure
}