use anyhow::{Context, Result};
use serde::{Deserialize, Serialize};
use std::path::PathBuf;
use std::process::Command;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Worktree {
    pub path: PathBuf,
    pub branch: String,
    pub head: String,
    pub is_bare: bool,
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
        let mut cmd = Command::new("git");
        cmd.args(["worktree", "add"]);

        let target_path = if let Some(p) = path {
            p.to_string()
        } else {
            format!("../{}", branch)
        };

        cmd.arg(&target_path).arg(branch);

        let output = cmd.output()
            .context("Failed to execute git worktree add")?;

        if !output.status.success() {
            let stderr = String::from_utf8_lossy(&output.stderr);
            anyhow::bail!("git worktree add failed: {}", stderr);
        }

        Ok(())
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