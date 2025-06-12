use anyhow::{Context, Result};
use std::path::{Path, PathBuf};
use std::fs::OpenOptions;
use std::io::{BufRead, BufReader, Write, BufWriter};
use std::time::{SystemTime, UNIX_EPOCH};

/// Represents a z database entry
#[derive(Debug, Clone)]
pub struct ZEntry {
    pub path: PathBuf,
    pub rank: f64,
    pub timestamp: u64,
}

impl ZEntry {
    pub fn from_line(line: &str) -> Option<Self> {
        let parts: Vec<&str> = line.split('|').collect();
        if parts.len() != 3 {
            return None;
        }

        Some(ZEntry {
            path: PathBuf::from(parts[0]),
            rank: parts[1].parse().ok()?,
            timestamp: parts[2].parse().ok()?,
        })
    }

    pub fn to_line(&self) -> String {
        format!("{}|{}|{}", self.path.display(), self.rank, self.timestamp)
    }

    /// Calculate frecency score for sorting
    pub fn frecency_score(&self) -> f64 {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs();
        
        let age_hours = (now.saturating_sub(self.timestamp)) as f64 / 3600.0;
        
        // Recency multiplier similar to z algorithm
        let recency_weight = if age_hours < 1.0 {
            4.0
        } else if age_hours < 24.0 {
            2.0
        } else if age_hours < 168.0 { // 1 week
            1.0
        } else {
            0.5
        };
        
        self.rank * recency_weight
    }
}

/// Manager for z database integration
pub struct ZIntegration {
    data_file: PathBuf,
}

impl ZIntegration {
    pub fn new() -> Self {
        let data_file = dirs::home_dir()
            .unwrap_or_else(|| PathBuf::from("."))
            .join(".z");
        
        Self { data_file }
    }

    /// Read all entries from z database
    pub fn read_entries(&self) -> Result<Vec<ZEntry>> {
        if !self.data_file.exists() {
            return Ok(Vec::new());
        }

        let file = std::fs::File::open(&self.data_file)
            .with_context(|| format!("Failed to open z data file: {}", self.data_file.display()))?;
        
        let reader = BufReader::new(file);
        let mut entries = Vec::new();

        for line in reader.lines() {
            let line = line?;
            if let Some(entry) = ZEntry::from_line(&line) {
                entries.push(entry);
            }
        }

        Ok(entries)
    }

    /// Write entries back to z database
    pub fn write_entries(&self, entries: &[ZEntry]) -> Result<()> {
        // Create parent directory if it doesn't exist
        if let Some(parent) = self.data_file.parent() {
            std::fs::create_dir_all(parent)?;
        }

        let file = OpenOptions::new()
            .write(true)
            .truncate(true)
            .create(true)
            .open(&self.data_file)
            .with_context(|| format!("Failed to open z data file for writing: {}", self.data_file.display()))?;
        
        let mut writer = BufWriter::new(file);
        
        for entry in entries {
            writeln!(writer, "{}", entry.to_line())?;
        }
        
        writer.flush()?;
        Ok(())
    }

    /// Add a worktree path to z database
    pub fn add_worktree<P: AsRef<Path>>(&self, path: P, alias: Option<String>) -> Result<()> {
        let canonical_path = path.as_ref().canonicalize()
            .unwrap_or_else(|_| path.as_ref().to_path_buf());
        
        let mut entries = self.read_entries()?;
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs();

        // Check if path already exists and update it
        let mut found = false;
        for entry in &mut entries {
            if entry.path == canonical_path {
                entry.rank += 1.0;
                entry.timestamp = now;
                found = true;
                break;
            }
        }

        // Add new entry if not found
        if !found {
            entries.push(ZEntry {
                path: canonical_path.clone(),
                rank: 1.0,
                timestamp: now,
            });
        }

        // Add alias entry if provided
        if let Some(alias_name) = alias {
            let alias_path = canonical_path.parent()
                .unwrap_or(&canonical_path)
                .join(&alias_name);
            
            // Add alias entry pointing to the worktree
            let mut alias_found = false;
            for entry in &mut entries {
                if entry.path == alias_path {
                    entry.rank += 0.5; // Lower rank for alias
                    entry.timestamp = now;
                    alias_found = true;
                    break;
                }
            }

            if !alias_found {
                entries.push(ZEntry {
                    path: alias_path,
                    rank: 0.5,
                    timestamp: now,
                });
            }
        }

        self.write_entries(&entries)?;
        Ok(())
    }

    /// Remove worktree paths from z database
    pub fn remove_worktree<P: AsRef<Path>>(&self, path: P) -> Result<()> {
        let canonical_path = path.as_ref().canonicalize()
            .unwrap_or_else(|_| path.as_ref().to_path_buf());
        
        let entries = self.read_entries()?;
        let filtered_entries: Vec<ZEntry> = entries
            .into_iter()
            .filter(|entry| entry.path != canonical_path)
            .collect();

        self.write_entries(&filtered_entries)?;
        Ok(())
    }

    /// Clean up non-existent worktree paths
    pub fn cleanup_worktrees(&self) -> Result<Vec<PathBuf>> {
        let entries = self.read_entries()?;
        let mut removed_paths = Vec::new();
        
        let existing_entries: Vec<ZEntry> = entries
            .into_iter()
            .filter(|entry| {
                if entry.path.exists() {
                    true
                } else {
                    removed_paths.push(entry.path.clone());
                    false
                }
            })
            .collect();

        if !removed_paths.is_empty() {
            self.write_entries(&existing_entries)?;
        }

        Ok(removed_paths)
    }

    /// Search for worktree paths in z database
    pub fn search_worktrees(&self, query: &str) -> Result<Vec<ZEntry>> {
        let entries = self.read_entries()?;
        let query_lower = query.to_lowercase();
        
        let mut matching_entries: Vec<ZEntry> = entries
            .into_iter()
            .filter(|entry| {
                entry.path.exists() && 
                (entry.path.to_string_lossy().to_lowercase().contains(&query_lower) ||
                 entry.path.file_name()
                     .and_then(|n| n.to_str())
                     .map(|n| n.to_lowercase().contains(&query_lower))
                     .unwrap_or(false))
            })
            .collect();

        // Sort by frecency score
        matching_entries.sort_by(|a, b| b.frecency_score().partial_cmp(&a.frecency_score()).unwrap());
        
        Ok(matching_entries)
    }

    /// Get the best match for a query
    #[allow(dead_code)]
    pub fn find_best_match(&self, query: &str) -> Result<Option<PathBuf>> {
        let matches = self.search_worktrees(query)?;
        Ok(matches.first().map(|entry| entry.path.clone()))
    }

    /// Create smart aliases for worktrees (project-branch format)
    pub fn create_smart_alias<P: AsRef<Path>>(&self, worktree_path: P, branch_name: &str) -> Result<()> {
        let path = worktree_path.as_ref();
        
        // Try to determine project name from git remote or directory name
        let project_name = self.get_project_name(path)?;
        let smart_alias = format!("{}-{}", project_name, branch_name);
        
        self.add_worktree(path, Some(smart_alias))?;
        Ok(())
    }

    /// Extract project name from git remote or directory structure
    fn get_project_name<P: AsRef<Path>>(&self, path: P) -> Result<String> {
        let path = path.as_ref();
        
        // Try to get project name from git remote
        if let Ok(output) = std::process::Command::new("git")
            .args(&["remote", "get-url", "origin"])
            .current_dir(path)
            .output()
        {
            if output.status.success() {
                let remote_url = String::from_utf8_lossy(&output.stdout);
                if let Some(project_name) = extract_project_from_remote(&remote_url) {
                    return Ok(project_name);
                }
            }
        }
        
        // Fallback to directory name
        Ok(path.file_name()
            .and_then(|n| n.to_str())
            .unwrap_or("unknown")
            .to_string())
    }
}

/// Extract project name from git remote URL
fn extract_project_from_remote(remote_url: &str) -> Option<String> {
    let url = remote_url.trim();
    
    // Handle GitHub SSH URLs: git@github.com:user/repo.git
    if url.starts_with("git@github.com:") {
        return url.strip_prefix("git@github.com:")?
            .strip_suffix(".git")
            .or_else(|| Some(url.strip_prefix("git@github.com:")?))
            .map(|s| s.split('/').last().unwrap_or("unknown").to_string());
    }
    
    // Handle HTTPS URLs: https://github.com/user/repo.git
    if url.starts_with("https://github.com/") {
        return url.strip_prefix("https://github.com/")?
            .strip_suffix(".git")
            .or_else(|| Some(url.strip_prefix("https://github.com/")?))
            .map(|s| s.split('/').last().unwrap_or("unknown").to_string());
    }
    
    // Generic fallback for other remotes
    url.split('/').last()
        .and_then(|s| s.strip_suffix(".git").or(Some(s)))
        .map(|s| s.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_z_entry_parsing() {
        let line = "/Users/test/project|5.5|1640995200";
        let entry = ZEntry::from_line(line).unwrap();
        
        assert_eq!(entry.path, PathBuf::from("/Users/test/project"));
        assert_eq!(entry.rank, 5.5);
        assert_eq!(entry.timestamp, 1640995200);
    }

    #[test]
    fn test_z_entry_serialization() {
        let entry = ZEntry {
            path: PathBuf::from("/Users/test/project"),
            rank: 5.5,
            timestamp: 1640995200,
        };
        
        let line = entry.to_line();
        assert_eq!(line, "/Users/test/project|5.5|1640995200");
    }

    #[test]
    fn test_extract_project_from_remote() {
        assert_eq!(
            extract_project_from_remote("git@github.com:user/project.git"),
            Some("project".to_string())
        );
        
        assert_eq!(
            extract_project_from_remote("https://github.com/user/project.git"),
            Some("project".to_string())
        );
        
        assert_eq!(
            extract_project_from_remote("https://github.com/user/project"),
            Some("project".to_string())
        );
    }
}