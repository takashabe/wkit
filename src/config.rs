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
        // レガシー設定のマイグレーションを実行
        let _ = Self::migrate_legacy_config();

        // Try to load from local config first, then global config
        if let Ok(config) = Self::load_from_path(".wkit.toml") {
            return Ok(config);
        }

        if let Some(global_path) = Self::global_config_path() {
            if global_path.exists() {
                return Self::load_from_path(&global_path);
            }
        }

        // Return default config if no config files exist
        Ok(Self::default())
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

        for file_path in &self.copy_files.files {
            let source_file = source_dir.join(file_path);
            let target_file = target_dir.join(file_path);

            if !source_file.exists() {
                continue;
            }

            // Create parent directories if needed
            if let Some(parent) = target_file.parent() {
                fs::create_dir_all(parent)
                    .with_context(|| format!("Failed to create directory: {}", parent.display()))?;
            }

            // Skip if target file already exists
            if target_file.exists() {
                continue;
            }

            fs::copy(&source_file, &target_file)
                .with_context(|| format!("Failed to copy {} to {}",
                    source_file.display(), target_file.display()))?;

            copied_files.push(file_path.clone());
        }

        Ok(copied_files)
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
}
