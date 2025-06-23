use anyhow::{Context, Result};
use serde::{Deserialize, Serialize};
use std::path::{Path, PathBuf};
use std::fs;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Config {
    #[serde(default = "default_worktree_path")]
    pub default_worktree_path: String,
    #[serde(default)]
    pub auto_cleanup: bool,
    #[serde(default)]
    pub z_integration: bool,
    #[serde(default = "default_sync_strategy")]
    pub default_sync_strategy: String,
    #[serde(default = "default_main_branch")]
    pub main_branch: String,
    #[serde(default)]
    pub copy_files: CopyFiles,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CopyFiles {
    #[serde(default)]
    pub enabled: bool,
    #[serde(default = "default_copy_files")]
    pub files: Vec<String>,
}

impl Default for Config {
    fn default() -> Self {
        Self {
            default_worktree_path: default_worktree_path(),
            auto_cleanup: false,
            z_integration: true,
            default_sync_strategy: default_sync_strategy(),
            main_branch: default_main_branch(),
            copy_files: CopyFiles::default(),
        }
    }
}

impl Default for CopyFiles {
    fn default() -> Self {
        Self {
            enabled: false,
            files: default_copy_files(),
        }
    }
}

fn default_worktree_path() -> String {
    ".git/.wkit-worktrees".to_string()
}

fn default_sync_strategy() -> String {
    "merge".to_string()
}

fn default_main_branch() -> String {
    "main".to_string()
}

fn default_copy_files() -> Vec<String> {
    vec![
        ".envrc".to_string(),
        "compose.override.yaml".to_string(),
        ".env.local".to_string(),
        "config/local.yaml".to_string(),
    ]
}

impl Config {
    pub fn load() -> Result<Self> {
        let (config, _) = Self::load_with_source()?;
        Ok(config)
    }

    pub fn load_with_source() -> Result<(Self, Option<PathBuf>)> {
        // レガシー設定のマイグレーションを実行
        let _ = Self::migrate_legacy_config();

        // Try to load from local config first, then global config
        let local_path = Path::new(".wkit.toml");
        if local_path.exists() {
            let config = Self::load_from_path(local_path)?;
            let absolute_path = std::env::current_dir()?.join(local_path);
            return Ok((config, Some(absolute_path)));
        }

        if let Some(global_path) = Self::global_config_path() {
            if global_path.exists() {
                let config = Self::load_from_path(&global_path)?;
                return Ok((config, Some(global_path)));
            }
        }

        // Return default config if no config files exist
        Ok((Self::default(), None))
    }

    fn load_from_path<P: AsRef<Path>>(path: P) -> Result<Self> {
        let content = fs::read_to_string(path.as_ref())
            .with_context(|| format!("Failed to read config file: {}", path.as_ref().display()))?;

        let config: Self = toml::from_str(&content)
            .context("Failed to parse config file")?;

        Ok(config)
    }

    pub fn save_global(&self) -> Result<()> {
        let global_path = Self::global_config_path()
            .ok_or_else(|| anyhow::anyhow!("Cannot determine global config directory"))?;

        if let Some(parent) = global_path.parent() {
            fs::create_dir_all(parent)
                .context("Failed to create config directory")?;
        }

        let content = toml::to_string_pretty(self)
            .context("Failed to serialize config")?;

        fs::write(&global_path, content)
            .with_context(|| format!("Failed to write config file: {}", global_path.display()))?;

        Ok(())
    }

    pub fn save_local(&self) -> Result<()> {
        let content = toml::to_string_pretty(self)
            .context("Failed to serialize config")?;

        fs::write(".wkit.toml", content)
            .context("Failed to write local config file")?;

        Ok(())
    }

    fn global_config_path() -> Option<PathBuf> {
        // XDG_CONFIG_HOME環境変数があれば優先、なければ~/.config
        if let Ok(xdg_config_home) = std::env::var("XDG_CONFIG_HOME") {
            Some(PathBuf::from(xdg_config_home).join("wkit").join("config.toml"))
        } else {
            dirs::config_dir().map(|dir| dir.join("wkit").join("config.toml"))
        }
    }

    fn legacy_global_config_path() -> Option<PathBuf> {
        // 以前のApplication Support形式のパス
        dirs::data_dir().map(|dir| dir.join("wkit").join("config.toml"))
    }

    fn migrate_legacy_config() -> Result<()> {
        let new_path = Self::global_config_path();
        let legacy_path = Self::legacy_global_config_path();

        if let (Some(new_path), Some(legacy_path)) = (new_path, legacy_path) {
            // 新しいパスが存在せず、古いパスが存在する場合のみマイグレーション
            if !new_path.exists() && legacy_path.exists() {
                
                // 新しいディレクトリを作成
                if let Some(parent) = new_path.parent() {
                    fs::create_dir_all(parent)
                        .context("Failed to create new config directory")?;
                }

                // ファイルをコピー
                fs::copy(&legacy_path, &new_path)
                    .context("Failed to migrate config file")?;

                // 古いファイルを削除
                fs::remove_file(&legacy_path)
                    .context("Failed to remove legacy config file")?;

                // 古いディレクトリが空なら削除
                if let Some(parent) = legacy_path.parent() {
                    if parent.read_dir().map(|mut d| d.next().is_none()).unwrap_or(false) {
                        let _ = fs::remove_dir(parent);
                    }
                }
            }
        }
        Ok(())
    }

    pub fn resolve_worktree_path(&self, branch: &str, provided_path: Option<&str>) -> String {
        match provided_path {
            Some(path) => path.to_string(),
            None => {
                let base_path = if self.default_worktree_path.starts_with('/') {
                    self.default_worktree_path.clone()
                } else {
                    self.default_worktree_path.clone()
                };
                format!("{}/{}", base_path, branch)
            }
        }
    }

    pub fn copy_files_to_worktree(&self, source_dir: &Path, target_dir: &Path) -> Result<Vec<String>> {
        if !self.copy_files.enabled {
            return Ok(vec![]);
        }

        let mut copied_files = Vec::new();

        for file_pattern in &self.copy_files.files {
            // Check if it's a relative path or just a filename
            if file_pattern.contains('/') || file_pattern.contains('\\') {
                // It's a path, use the existing logic
                let source_file = source_dir.join(file_pattern);
                let target_file = target_dir.join(file_pattern);

                if source_file.exists() {
                    self.copy_single_file(&source_file, &target_file, file_pattern, &mut copied_files)?;
                }
            } else {
                // It's just a filename, search for all matching files in the repository
                let found_files = self.find_files_by_name(source_dir, file_pattern)?;
                for relative_path in found_files {
                    let source_file = source_dir.join(&relative_path);
                    let target_file = target_dir.join(&relative_path);
                    
                    self.copy_single_file(&source_file, &target_file, &relative_path, &mut copied_files)?;
                }
            }
        }

        Ok(copied_files)
    }

    fn copy_single_file(&self, source_file: &Path, target_file: &Path, relative_path: &str, copied_files: &mut Vec<String>) -> Result<()> {
        // Create parent directories if needed
        if let Some(parent) = target_file.parent() {
            fs::create_dir_all(parent)
                .with_context(|| format!("Failed to create directory: {}", parent.display()))?;
        }

        // Skip if target file already exists
        if target_file.exists() {
            return Ok(());
        }

        fs::copy(source_file, target_file)
            .with_context(|| format!("Failed to copy {} to {}",
                source_file.display(), target_file.display()))?;

        copied_files.push(relative_path.to_string());
        Ok(())
    }

    fn find_files_by_name(&self, source_dir: &Path, filename: &str) -> Result<Vec<String>> {
        let mut found_files = Vec::new();
        self.walk_directory(source_dir, source_dir, filename, &mut found_files)?;
        Ok(found_files)
    }

    fn walk_directory(&self, base_dir: &Path, current_dir: &Path, filename: &str, found_files: &mut Vec<String>) -> Result<()> {
        let entries = fs::read_dir(current_dir)
            .with_context(|| format!("Failed to read directory: {}", current_dir.display()))?;

        for entry in entries {
            let entry = entry?;
            let path = entry.path();
            
            // Skip .git directory to avoid copying from other worktrees
            if path.file_name().map_or(false, |name| name == ".git") {
                continue;
            }

            if path.is_file() {
                if let Some(file_name) = path.file_name() {
                    if file_name == filename {
                        // Calculate relative path from base directory
                        let relative_path = path.strip_prefix(base_dir)
                            .with_context(|| format!("Failed to get relative path for: {}", path.display()))?;
                        found_files.push(relative_path.to_string_lossy().to_string());
                    }
                }
            } else if path.is_dir() {
                // Recursively search subdirectories
                self.walk_directory(base_dir, &path, filename, found_files)?;
            }
        }
        
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_default_config() {
        let config = Config::default();
        assert_eq!(config.default_worktree_path, ".git/.wkit-worktrees");
        assert!(!config.auto_cleanup);
        assert!(config.z_integration);
    }

    #[test]
    fn test_resolve_worktree_path_with_provided_path() {
        let config = Config::default();
        let result = config.resolve_worktree_path("feature", Some("/custom/path"));
        assert_eq!(result, "/custom/path");
    }

    #[test]
    fn test_resolve_worktree_path_with_default_relative() {
        let config = Config::default();
        let result = config.resolve_worktree_path("feature", None);
        assert_eq!(result, ".git/.wkit-worktrees/feature");
    }

    #[test]
    fn test_resolve_worktree_path_with_absolute_default() {
        let mut config = Config::default();
        config.default_worktree_path = "/absolute/path".to_string();
        let result = config.resolve_worktree_path("feature", None);
        assert_eq!(result, "/absolute/path/feature");
    }

    #[test]
    fn test_config_serialization() {
        let config = Config::default();
        let toml_str = toml::to_string(&config).unwrap();
        assert!(toml_str.contains("default_worktree_path"));
        assert!(toml_str.contains("auto_cleanup"));
        assert!(toml_str.contains("z_integration"));
    }

    #[test]
    fn test_config_deserialization() {
        let toml_str = r#"
default_worktree_path = "/custom"
auto_cleanup = true
z_integration = false
"#;
        let config: Config = toml::from_str(toml_str).unwrap();
        assert_eq!(config.default_worktree_path, "/custom");
        assert!(config.auto_cleanup);
        assert!(!config.z_integration);
    }

    #[test]
    fn test_global_config_path_with_xdg_config_home() {
        std::env::set_var("XDG_CONFIG_HOME", "/custom/config");
        let path = Config::global_config_path().unwrap();
        assert_eq!(path, PathBuf::from("/custom/config/wkit/config.toml"));
        std::env::remove_var("XDG_CONFIG_HOME");
    }

    #[test]
    fn test_global_config_path_without_xdg_config_home() {
        std::env::remove_var("XDG_CONFIG_HOME");
        let path = Config::global_config_path();
        assert!(path.is_some());
        assert!(path.unwrap().to_string_lossy().contains("wkit/config.toml"));
    }

    #[test]
    fn test_copy_files_disabled() {
        let config = Config::default();
        let result = config.copy_files_to_worktree(Path::new("/source"), Path::new("/target"));
        assert!(result.is_ok());
        assert!(result.unwrap().is_empty());
    }

    #[test]
    fn test_find_files_by_name_logic() {
        use std::fs;
        use std::path::PathBuf;
        
        let temp_dir = std::env::temp_dir().join("wkit_test");
        let _ = fs::remove_dir_all(&temp_dir);
        fs::create_dir_all(&temp_dir).unwrap();
        
        // Create test files
        fs::create_dir_all(temp_dir.join("src")).unwrap();
        fs::create_dir_all(temp_dir.join("config")).unwrap();
        fs::write(temp_dir.join(".envrc"), "root envrc").unwrap();
        fs::write(temp_dir.join("src/.envrc"), "src envrc").unwrap();
        fs::write(temp_dir.join("config/.envrc"), "config envrc").unwrap();
        
        let mut config = Config::default();
        config.copy_files.enabled = true;
        
        let found_files = config.find_files_by_name(&temp_dir, ".envrc").unwrap();
        assert_eq!(found_files.len(), 3);
        assert!(found_files.contains(&".envrc".to_string()));
        assert!(found_files.contains(&"src/.envrc".to_string()));
        assert!(found_files.contains(&"config/.envrc".to_string()));
        
        // Clean up
        let _ = fs::remove_dir_all(&temp_dir);
    }
}
