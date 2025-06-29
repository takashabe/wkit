package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Executor handles Git command execution
type Executor struct {
	workDir string
}

// NewExecutor creates a new Git command executor
func NewExecutor(workDir string) *Executor {
	return &Executor{workDir: workDir}
}

// Execute runs a Git command with the given arguments
func (e *Executor) Execute(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if e.workDir != "" {
		cmd.Dir = e.workDir
	}

	output, err := cmd.Output()
	if err != nil {
		// Try to get stderr for better error messages
		if exitError, ok := err.(*exec.ExitError); ok {
			stderr := string(exitError.Stderr)
			return "", fmt.Errorf("git %s failed: %w\nStderr: %s", strings.Join(args, " "), err, stderr)
		}
		return "", fmt.Errorf("git %s failed: %w", strings.Join(args, " "), err)
	}

	return strings.TrimSpace(string(output)), nil
}

// ExecuteWithStderr runs a Git command and returns both stdout and stderr
func (e *Executor) ExecuteWithStderr(args ...string) (string, string, error) {
	cmd := exec.Command("git", args...)
	if e.workDir != "" {
		cmd.Dir = e.workDir
	}

	stdout, err := cmd.Output()
	stderr := ""
	
	if exitError, ok := err.(*exec.ExitError); ok {
		stderr = string(exitError.Stderr)
	}

	if err != nil {
		return string(stdout), stderr, fmt.Errorf("git %s failed: %w", strings.Join(args, " "), err)
	}

	return strings.TrimSpace(string(stdout)), stderr, nil
}

// GetRepositoryRoot returns the absolute path to the repository root
func GetRepositoryRoot() (string, error) {
	executor := NewExecutor("")
	return executor.Execute("rev-parse", "--show-toplevel")
}

// WorktreeList returns the output of 'git worktree list --porcelain'
func (e *Executor) WorktreeList() (string, error) {
	return e.Execute("worktree", "list", "--porcelain")
}

// WorktreeAdd adds a new worktree
func (e *Executor) WorktreeAdd(args ...string) error {
	cmdArgs := append([]string{"worktree", "add"}, args...)
	_, err := e.Execute(cmdArgs...)
	return err
}

// WorktreeRemove removes a worktree
func (e *Executor) WorktreeRemove(path string, force bool) error {
	args := []string{"worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, path)
	
	_, err := e.Execute(args...)
	return err
}

// Status returns git status output in porcelain format
func (e *Executor) Status() (string, error) {
	return e.Execute("status", "--porcelain", "--ahead-behind")
}

// BranchExists checks if a local branch exists
func (e *Executor) BranchExists(branch string) bool {
	_, err := e.Execute("show-ref", "--verify", "--quiet", fmt.Sprintf("refs/heads/%s", branch))
	return err == nil
}

// BranchMerged returns branches merged into the specified branch
func (e *Executor) BranchMerged(branch string) ([]string, error) {
	output, err := e.Execute("branch", "--merged", branch)
	if err != nil {
		return nil, err
	}

	var branches []string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Remove the "* " prefix if present
		branch := strings.TrimPrefix(line, "* ")
		branches = append(branches, branch)
	}
	
	return branches, nil
}

// RemoteBranches returns all remote branches
func (e *Executor) RemoteBranches(remote string) ([]string, error) {
	output, err := e.Execute("ls-remote", "--heads", remote)
	if err != nil {
		return nil, err
	}

	var branches []string
	lines := strings.Split(output, "\n")
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

// Fetch fetches from origin
func (e *Executor) Fetch(remote string) error {
	_, err := e.Execute("fetch", remote)
	return err
}

// Merge merges the specified branch
func (e *Executor) Merge(branch string) error {
	_, err := e.Execute("merge", branch)
	return err
}

// Rebase rebases onto the specified branch
func (e *Executor) Rebase(branch string) error {
	_, err := e.Execute("rebase", branch)
	return err
}