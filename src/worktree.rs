use anyhow::{Context, Result};
use serde::{Deserialize, Serialize};
use std::path::{Path, PathBuf};
use std::process::Command;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Worktree {
    pub path: PathBuf,
    pub branch: String,
    pub head: String,
    pub is_bare: bool,
}

#[derive(Debug, Clone)]
pub struct WorktreeStatus {
    pub is_clean: bool,
    pub modified: usize,
    pub added: usize,
    pub deleted: usize,
    pub untracked: usize,
    pub ahead: usize,
    pub behind: usize,
}

pub struct WorktreeManager;

impl WorktreeManager {
    pub fn new() -> Self {
        Self
    }

    pub fn list_worktrees(&self) -> Result<Vec<Worktree>> {
        let output = Command::new("git")
            .args(["worktree", "list", "--porcelain"])
            .output()
            .context("Failed to execute git worktree list")?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            anyhow::bail!("git worktree list failed: {}", stderr);
        }

        let stdout = String::from_utf8_lossy(&output.stdout);
        self.parse_worktree_list(&stdout)
    }

    fn parse_worktree_list(&self, output: &str) -> Result<Vec<Worktree>> {
        let mut worktrees = Vec::new();
        let mut current_worktree: Option<Worktree> = None;

        for line in output.lines() {
            if line.starts_with("worktree ") {
                if let Some(wt) = current_worktree.take() {
                    worktrees.push(wt);
                }
                let path = line.strip_prefix("worktree ").unwrap();
                current_worktree = Some(Worktree {
                    path: PathBuf::from(path),
                    branch: String::new(),
                    head: String::new(),
                    is_bare: false,
                });
            } else if line.starts_with("HEAD ") {
                if let Some(ref mut wt) = current_worktree {
                    wt.head = line.strip_prefix("HEAD ").unwrap().to_string();
                }
            } else if line.starts_with("branch ") {
                if let Some(ref mut wt) = current_worktree {
                    let branch = line.strip_prefix("branch ").unwrap();
                    wt.branch = branch.strip_prefix("refs/heads/").unwrap_or(branch).to_string();
                }
            } else if line == "bare" {
                if let Some(ref mut wt) = current_worktree {
                    wt.is_bare = true;
                }
            }
        }

        if let Some(wt) = current_worktree {
            worktrees.push(wt);
        }

        Ok(worktrees)
    }

    pub fn add_worktree(&self, branch: &str, path: Option<&str>) -> Result<()> {
        // path should be provided by the caller (using config)
        let target_path = path
            .ok_or_else(|| anyhow::anyhow!("Path must be provided"))?
            .to_string();

        let mut cmd = Command::new("git");
        cmd.args(["worktree", "add"]);

        // Check if branch exists
        let branch_exists = self.branch_exists(branch)?;
        
        if !branch_exists {
            // Use -b flag to create new branch
            cmd.arg("-b").arg(branch).arg(&target_path);
        } else {
            cmd.arg(&target_path).arg(branch);
        }

        let output = cmd.output()
            .context("Failed to execute git worktree add")?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            anyhow::bail!("git worktree add failed: {}", stderr);
        }

        Ok(())
    }

    fn branch_exists(&self, branch: &str) -> Result<bool> {
        let output = Command::new("git")
            .args(["show-ref", "--verify", "--quiet", &format!("refs/heads/{}", branch)])
            .output()
            .context("Failed to check if branch exists")?;

        Ok(output.status.success())
    }

    pub fn remove_worktree(&self, path: &str) -> Result<()> {
        let output = Command::new("git")
            .args(["worktree", "remove", path])
            .output()
            .context("Failed to execute git worktree remove")?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            anyhow::bail!("git worktree remove failed: {}", stderr);
        }

        Ok(())
    }

    pub fn find_worktree_by_name(&self, name: &str) -> Result<Option<Worktree>> {
        let worktrees = self.list_worktrees()?;
        
        // Exact match by branch name
        if let Some(wt) = worktrees.iter().find(|w| w.branch == name) {
            return Ok(Some(wt.clone()));
        }

        // Partial match by path
        if let Some(wt) = worktrees.iter().find(|w| 
            w.path.file_name()
                .and_then(|n| n.to_str())
                .map(|n| n.contains(name))
                .unwrap_or(false)
        ) {
            return Ok(Some(wt.clone()));
        }

        Ok(None)
    }

    pub fn get_worktree_status(&self, worktree_path: &Path) -> Result<WorktreeStatus> {
        let output = Command::new("git")
            .args(["status", "--porcelain", "--ahead-behind"])
            .current_dir(worktree_path)
            .output()
            .context("Failed to execute git status")?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            anyhow::bail!("git status failed: {}", stderr);
        }

        let stdout = String::from_utf8_lossy(&output.stdout);
        self.parse_git_status(&stdout)
    }

    fn parse_git_status(&self, output: &str) -> Result<WorktreeStatus> {
        let mut modified = 0;
        let mut added = 0;
        let mut deleted = 0;
        let mut untracked = 0;

        for line in output.lines() {
            if line.len() < 2 {
                continue;
            }

            let staged = &line[0..1];
            let unstaged = &line[1..2];

            match (staged, unstaged) {
                ("A", _) => added += 1,
                ("M", _) | (_, "M") => modified += 1,
                ("D", _) | (_, "D") => deleted += 1,
                ("?", "?") => untracked += 1,
                _ => {}
            }
        }

        let is_clean = modified == 0 && added == 0 && deleted == 0 && untracked == 0;

        Ok(WorktreeStatus {
            is_clean,
            modified,
            added,
            deleted,
            untracked,
            ahead: 0,
            behind: 0,
        })
    }

    pub fn find_unnecessary_worktrees_with_main(&self, main_branch: &str) -> Result<Vec<(Worktree, String)>> {
        let worktrees = self.list_worktrees()?;
        let mut unnecessary = Vec::new();

        for worktree in worktrees {
            if worktree.is_bare {
                continue;
            }

            // Skip the main branch itself
            if worktree.branch == main_branch {
                continue;
            }

            // Check if branch is merged into main
            if self.is_branch_merged_into(&worktree.branch, main_branch)? {
                unnecessary.push((worktree.clone(), format!("Branch merged into {}", main_branch)));
                continue;
            }

            // Check if worktree path doesn't exist
            if !worktree.path.exists() {
                unnecessary.push((worktree.clone(), "Worktree path does not exist".to_string()));
                continue;
            }

            // Check if branch doesn't exist remotely
            if self.is_branch_deleted_remotely(&worktree.branch)? {
                unnecessary.push((worktree.clone(), "Branch deleted remotely".to_string()));
            }
        }

        Ok(unnecessary)
    }

    fn is_branch_merged_into(&self, branch: &str, main_branch: &str) -> Result<bool> {
        let output = Command::new("git")
            .args(["branch", "--merged", main_branch])
            .output()
            .context("Failed to check merged branches")?;

        if !output.status.success() {
            return Ok(false);
        }

        let stdout = String::from_utf8_lossy(&output.stdout);
        Ok(stdout.lines().any(|line| line.trim() == branch || line.trim() == &format!("* {}", branch)))
    }

    fn is_branch_deleted_remotely(&self, branch: &str) -> Result<bool> {
        let output = Command::new("git")
            .args(["ls-remote", "--heads", "origin", branch])
            .output()
            .context("Failed to check remote branch")?;

        if !output.status.success() {
            return Ok(false);
        }

        let stdout = String::from_utf8_lossy(&output.stdout);
        Ok(stdout.trim().is_empty())
    }

    pub fn sync_worktree_with_branch(&self, worktree: &Worktree, main_branch: &str, use_rebase: bool) -> Result<()> {
        // First, fetch latest changes
        let output = Command::new("git")
            .args(["fetch", "origin"])
            .current_dir(&worktree.path)
            .output()
            .context("Failed to fetch from origin")?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            anyhow::bail!("git fetch failed: {}", stderr);
        }

        // Then sync with specified branch
        let origin_branch = format!("origin/{}", main_branch);
        let sync_cmd = if use_rebase {
            vec!["rebase", &origin_branch]
        } else {
            vec!["merge", &origin_branch]
        };

        let output = Command::new("git")
            .args(&sync_cmd)
            .current_dir(&worktree.path)
            .output()
            .with_context(|| format!("Failed to {} with {}", if use_rebase { "rebase" } else { "merge" }, main_branch))?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            anyhow::bail!("{} with {} failed: {}", if use_rebase { "Rebase" } else { "Merge" }, main_branch, stderr);
        }

        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_worktree_list_empty() {
        let manager = WorktreeManager::new();
        let result = manager.parse_worktree_list("").unwrap();
        assert!(result.is_empty());
    }

    #[test]
    fn test_parse_worktree_list_single() {
        let manager = WorktreeManager::new();
        let output = "worktree /path/to/repo\nHEAD abcd1234\nbranch refs/heads/main\n";
        let result = manager.parse_worktree_list(output).unwrap();
        
        assert_eq!(result.len(), 1);
        assert_eq!(result[0].path, PathBuf::from("/path/to/repo"));
        assert_eq!(result[0].branch, "main");
        assert_eq!(result[0].head, "abcd1234");
        assert!(!result[0].is_bare);
    }

    #[test]
    fn test_parse_worktree_list_multiple() {
        let manager = WorktreeManager::new();
        let output = r#"worktree /path/to/repo
HEAD abcd1234
branch refs/heads/main

worktree /path/to/feature
HEAD efgh5678
branch refs/heads/feature-branch

worktree /path/to/bare
HEAD ijkl9012
bare
"#;
        let result = manager.parse_worktree_list(output).unwrap();
        
        assert_eq!(result.len(), 3);
        
        assert_eq!(result[0].branch, "main");
        assert!(!result[0].is_bare);
        
        assert_eq!(result[1].branch, "feature-branch");
        assert!(!result[1].is_bare);
        
        assert_eq!(result[2].head, "ijkl9012");
        assert!(result[2].is_bare);
    }

    #[test]
    fn test_parse_worktree_list_detached_head() {
        let manager = WorktreeManager::new();
        let output = "worktree /path/to/repo\nHEAD abcd1234\ndetached\n";
        let result = manager.parse_worktree_list(output).unwrap();
        
        assert_eq!(result.len(), 1);
        assert_eq!(result[0].branch, "");
        assert_eq!(result[0].head, "abcd1234");
    }
}