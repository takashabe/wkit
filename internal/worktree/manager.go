package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Worktree represents a Git worktree
type Worktree struct {
	Path   string
	Branch string
	HEAD   string
}

// Manager handles Git worktree operations
type Manager struct {
	// repo *git.Repository // go-git の Repository オブジェクトは直接使わない
}

// NewManager creates a new WorktreeManager
func NewManager() (*Manager, error) {
	return &Manager{}, nil
}

// ListWorktrees lists all worktrees associated with the repository
func (m *Manager) ListWorktrees() ([]Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute git worktree list: %w", err)
	}

	return parseWorktreeList(string(output))
}

func parseWorktreeList(output string) ([]Worktree, error) {
	var worktrees []Worktree
	var currentWorktree *Worktree

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if currentWorktree != nil {
				worktrees = append(worktrees, *currentWorktree)
				currentWorktree = nil
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			if currentWorktree != nil {
				worktrees = append(worktrees, *currentWorktree)
			}
			path := strings.TrimPrefix(line, "worktree ")
			currentWorktree = &Worktree{
				Path: path,
			}
		} else if strings.HasPrefix(line, "HEAD ") {
			if currentWorktree != nil {
				currentWorktree.HEAD = strings.TrimPrefix(line, "HEAD ")
			}
		} else if strings.HasPrefix(line, "branch ") {
			if currentWorktree != nil {
				branch := strings.TrimPrefix(line, "branch ")
				currentWorktree.Branch = strings.TrimPrefix(branch, "refs/heads/")
			}
		}
	}

	if currentWorktree != nil {
		worktrees = append(worktrees, *currentWorktree)
	}

	return worktrees, nil
}

// GetRepositoryRoot returns the absolute path to the repository root.
func GetRepositoryRoot() (string, error) {
	// Get the common git directory first
	cmd := exec.Command("git", "rev-parse", "--git-common-dir")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute git rev-parse --git-common-dir: %w", err)
	}
	gitDir := strings.TrimSpace(string(output))
	
	// The repository root is the parent of .git directory
	if strings.HasSuffix(gitDir, "/.git") {
		return strings.TrimSuffix(gitDir, "/.git"), nil
	}
	
	// Fallback to show-toplevel if not a standard .git directory
	cmd = exec.Command("git", "rev-parse", "--show-toplevel")
	output, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute git rev-parse --show-toplevel: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// AddWorktree adds a new worktree
func (m *Manager) AddWorktree(branch string, path string, mainBranch string) error {
	// Check if branch exists
	branchExists := m.branchExists(branch)

	var cmdArgs []string
	cmdArgs = append(cmdArgs, "worktree", "add")

	if !branchExists {
		// Use -b flag to create new branch from origin/<main_branch>
		baseBranch := fmt.Sprintf("origin/%s", mainBranch)
		cmdArgs = append(cmdArgs, "-b", branch, path, baseBranch)
	} else {
		cmdArgs = append(cmdArgs, path, branch)
	}

	cmd := exec.Command("git", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute git worktree add: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}

// branchExists checks if a local branch exists
func (m *Manager) branchExists(branch string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", fmt.Sprintf("refs/heads/%s", branch))
	err := cmd.Run()
	return err == nil
}

// FindWorktreePath finds a worktree path by name or partial path
func (m *Manager) FindWorktreePath(name string) (string, error) {
	worktrees, err := m.ListWorktrees()
	if err != nil {
		return "", err
	}

	// Exact match by branch name
	for _, wt := range worktrees {
		if wt.Branch == name {
			return wt.Path, nil
		}
	}

	// Partial match by path
	for _, wt := range worktrees {
		if strings.Contains(wt.Path, name) {
			return wt.Path, nil
		}
	}

	return "", fmt.Errorf("worktree '%s' not found", name)
}

// WorktreeStatus represents the status of a worktree
type WorktreeStatus struct {
	IsClean   bool
	Modified  int
	Added     int
	Deleted   int
	Untracked int
	Ahead     int
	Behind    int
}

// GetWorktreeStatus gets the status of a specific worktree
func (m *Manager) GetWorktreeStatus(worktreePath string) (*WorktreeStatus, error) {
	cmd := exec.Command("git", "status", "--porcelain", "--ahead-behind")
	cmd.Dir = worktreePath
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute git status for %s: %w", worktreePath, err)
	}

	return parseGitStatus(string(output))
}

func parseGitStatus(output string) (*WorktreeStatus, error) {
	status := &WorktreeStatus{}
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if len(line) < 2 {
			continue
		}

		staged := string(line[0])
		unstaged := string(line[1])

		switch {
		case staged == "A":
			status.Added++
		case staged == "M" || unstaged == "M":
			status.Modified++
		case staged == "D" || unstaged == "D":
			status.Deleted++
		case staged == "?" && unstaged == "?":
			status.Untracked++
		}

		// Parse ahead/behind information
		if strings.HasPrefix(line, "##") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasPrefix(part, "ahead") {
					if n, err := fmt.Sscanf(part, "ahead %d", &status.Ahead); err == nil && n == 1 {
						// Successfully parsed
					}
				} else if strings.HasPrefix(part, "behind") {
					if n, err := fmt.Sscanf(part, "behind %d", &status.Behind); err == nil && n == 1 {
						// Successfully parsed
					}
				}
			}
		}
	}

	status.IsClean = (status.Modified == 0 && status.Added == 0 && status.Deleted == 0 && status.Untracked == 0)

	return status, nil
}

// UnnecessaryWorktree represents an unnecessary worktree with a reason
type UnnecessaryWorktree struct {
	Worktree Worktree
	Reason   string
}

// FindUnnecessaryWorktrees finds worktrees that are no longer needed
func (m *Manager) FindUnnecessaryWorktrees(mainBranch string) ([]UnnecessaryWorktree, error) {
	var unnecessary []UnnecessaryWorktree
	worktrees, err := m.ListWorktrees()
	if err != nil {
		return nil, err
	}

	mergedBranches, err := m.getAllMergedBranches(mainBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to get merged branches: %w", err)
	}

	remoteBranches, err := m.getAllRemoteBranches()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote branches: %w", err)
	}

	for _, wt := range worktrees {
		// Skip the main branch itself
		if wt.Branch == mainBranch {
			continue
		}

		// Check if branch is merged into main
		if containsString(mergedBranches, wt.Branch) {
			unnecessary = append(unnecessary, UnnecessaryWorktree{Worktree: wt, Reason: fmt.Sprintf("Branch merged into %s", mainBranch)})
			continue
		}

		// Check if worktree path doesn't exist
		if _, err := os.Stat(wt.Path); os.IsNotExist(err) {
			unnecessary = append(unnecessary, UnnecessaryWorktree{Worktree: wt, Reason: "Worktree path does not exist"})
			continue
		}

		// Check if branch doesn't exist remotely
		if !containsString(remoteBranches, wt.Branch) {
			unnecessary = append(unnecessary, UnnecessaryWorktree{Worktree: wt, Reason: "Branch deleted remotely"})
		}
	}

	return unnecessary, nil
}

func (m *Manager) getAllMergedBranches(mainBranch string) ([]string, error) {
	cmd := exec.Command("git", "branch", "--merged", mainBranch)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute git branch --merged: %w", err)
	}

	var branches []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		branches = append(branches, strings.TrimPrefix(line, "* "))
	}
	return branches, nil
}

func (m *Manager) getAllRemoteBranches() ([]string, error) {
	cmd := exec.Command("git", "ls-remote", "--heads", "origin")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute git ls-remote: %w", err)
	}

	var branches []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 1 {
			branch := strings.TrimPrefix(parts[1], "refs/heads/")
			branches = append(branches, branch)
		}
	}
	return branches, nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// SyncWorktreeWithBranch syncs a worktree with the main branch
func (m *Manager) SyncWorktreeWithBranch(worktreePath string, mainBranch string, useRebase bool) error {
	// First, fetch latest changes
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Dir = worktreePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute git fetch: %w: %s", err, strings.TrimSpace(string(output)))
	}

	// Then sync with specified branch
	originBranch := fmt.Sprintf("origin/%s", mainBranch)
	var syncCmdArgs []string
	if useRebase {
		syncCmdArgs = []string{"rebase", originBranch}
	} else {
		syncCmdArgs = []string{"merge", originBranch}
	}

	cmd = exec.Command("git", syncCmdArgs...)
	cmd.Dir = worktreePath
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute git %s: %w: %s", syncCmdArgs[0], err, strings.TrimSpace(string(output)))
	}

	return nil
}

// RemoveWorktree removes a worktree
func (m *Manager) RemoveWorktree(worktreePath string) error {
	cmd := exec.Command("git", "worktree", "remove", "--force", worktreePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute git worktree remove: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}