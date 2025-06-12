use clap::{Parser, Subcommand};
use anyhow::Result;

mod worktree;
mod config;

use worktree::WorktreeManager;
use config::Config;

#[derive(Parser)]
#[command(name = "wkit")]
#[command(about = "Convenient git worktree management toolkit")]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// List all worktrees
    List,
    /// Add a new worktree
    Add {
        /// Branch name
        branch: String,
        /// Path for the worktree (optional)
        path: Option<String>,
    },
    /// Remove a worktree
    Remove {
        /// Worktree path or name
        worktree: String,
    },
    /// Switch to a worktree
    Switch {
        /// Worktree path or name
        worktree: String,
    },
    /// Configuration management
    Config {
        #[command(subcommand)]
        config_cmd: ConfigCommands,
    },
}

#[derive(Subcommand)]
enum ConfigCommands {
    /// Show current configuration
    Show,
    /// Set a configuration value
    Set {
        /// Configuration key (default_worktree_path, auto_cleanup, z_integration)
        key: String,
        /// Configuration value
        value: String,
    },
    /// Initialize a local configuration file
    Init,
}

fn main() {
    let cli = Cli::parse();
    let manager = WorktreeManager::new();

    let result = match cli.command {
        Commands::List => cmd_list(&manager),
        Commands::Add { branch, path } => cmd_add(&manager, &branch, path.as_deref()),
        Commands::Remove { worktree } => cmd_remove(&manager, &worktree),
        Commands::Switch { worktree } => cmd_switch(&manager, &worktree),
        Commands::Config { config_cmd } => cmd_config(config_cmd),
    };

    if let Err(e) = result {
        eprintln!("Error: {}", e);
        std::process::exit(1);
    }
}

fn cmd_list(manager: &WorktreeManager) -> Result<()> {
    let worktrees = manager.list_worktrees()?;
    if worktrees.is_empty() {
        println!("No worktrees found.");
        return Ok(());
    }

    println!("{:<30} {:<20} {:<12}", "PATH", "BRANCH", "HEAD");
    println!("{}", "-".repeat(65));
    
    for wt in worktrees {
        let path_str = wt.path.to_string_lossy();
        let head_short = if wt.head.len() > 10 {
            &wt.head[..10]
        } else {
            &wt.head
        };
        
        println!("{:<30} {:<20} {:<12}", 
            path_str, 
            wt.branch, 
            head_short
        );
    }
    Ok(())
}

fn cmd_add(manager: &WorktreeManager, branch: &str, path: Option<&str>) -> Result<()> {
    let config = Config::load()?;
    let target_path = config.resolve_worktree_path(branch, path);
    
    manager.add_worktree(branch, Some(&target_path))?;
    println!("✓ Created worktree for branch '{}' at '{}'", branch, target_path);
    Ok(())
}

fn cmd_remove(manager: &WorktreeManager, worktree: &str) -> Result<()> {
    // Try to find the worktree by name first
    if let Some(wt) = manager.find_worktree_by_name(worktree)? {
        let path_str = wt.path.to_string_lossy();
        manager.remove_worktree(&path_str)?;
        println!("✓ Removed worktree at '{}'", path_str);
    } else {
        // Try direct path removal
        manager.remove_worktree(worktree)
            .map_err(|e| anyhow::anyhow!("Worktree '{}' not found. {}", worktree, e))?;
        println!("✓ Removed worktree at '{}'", worktree);
    }
    Ok(())
}

fn cmd_switch(manager: &WorktreeManager, worktree: &str) -> Result<()> {
    let wt = manager.find_worktree_by_name(worktree)?
        .ok_or_else(|| anyhow::anyhow!("Worktree '{}' not found", worktree))?;
    
    let path_str = wt.path.to_string_lossy();
    println!("{}", path_str);
    Ok(())
}

fn cmd_config(config_cmd: ConfigCommands) -> Result<()> {
    match config_cmd {
        ConfigCommands::Show => {
            let config = Config::load()?;
            println!("Current configuration:");
            println!("  default_worktree_path: {}", config.default_worktree_path);
            println!("  auto_cleanup: {}", config.auto_cleanup);
            println!("  z_integration: {}", config.z_integration);
        }
        ConfigCommands::Set { key, value } => {
            let mut config = Config::load()?;
            let value_clone = value.clone();
            
            match key.as_str() {
                "default_worktree_path" => {
                    config.default_worktree_path = value;
                }
                "auto_cleanup" => {
                    config.auto_cleanup = value.parse()
                        .map_err(|_| anyhow::anyhow!("Invalid boolean value: {}", value))?;
                }
                "z_integration" => {
                    config.z_integration = value.parse()
                        .map_err(|_| anyhow::anyhow!("Invalid boolean value: {}", value))?;
                }
                _ => {
                    anyhow::bail!("Unknown configuration key: {}. Valid keys: default_worktree_path, auto_cleanup, z_integration", key);
                }
            }
            
            config.save_global()?;
            println!("✓ Configuration updated: {} = {}", key, value_clone);
        }
        ConfigCommands::Init => {
            let config = Config::default();
            config.save_local()?;
            println!("✓ Created local configuration file: .wkit.toml");
        }
    }
    Ok(())
}