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
}

impl Default for Config {
    fn default() -> Self {
        Self {
            default_worktree_path: default_worktree_path(),
            auto_cleanup: false,
            z_integration: true,
        }
    }
}

fn default_worktree_path() -> String {
    "..".to_string()
}

impl Config {
    pub fn load() -> Result<Self> {
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
        dirs::config_dir().map(|dir| dir.join("wkit").join("config.toml"))
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
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_default_config() {
        let config = Config::default();
        assert_eq!(config.default_worktree_path, "..");
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
        assert_eq!(result, "../feature");
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
}