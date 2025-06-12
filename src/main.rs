use clap::{Parser, Subcommand};
use anyhow::Result;

mod worktree;
mod config;
mod z_integration;

use worktree::WorktreeManager;
use config::Config;
use z_integration::ZIntegration;

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
    /// Show git status of all worktrees
    Status,
    /// Clean up unnecessary worktrees
    Clean {
        /// Skip confirmation prompt
        #[arg(short, long)]
        force: bool,
    },
    /// Sync worktree with main branch
    Sync {
        /// Worktree path or name to sync
        worktree: Option<String>,
        /// Use rebase instead of merge
        #[arg(short, long)]
        rebase: bool,
    },
    /// Z-style frecency-based worktree jumping
    Z {
        /// Query string to search for worktrees
        query: Option<String>,
        /// List all matches instead of jumping
        #[arg(short, long)]
        list: bool,
        /// Clean up non-existent entries
        #[arg(short, long)]
        clean: bool,
        /// Add current directory to z database
        #[arg(short, long)]
        add: bool,
    },
}

#[derive(Subcommand)]
enum ConfigCommands {
    /// Show current configuration
    Show,
    /// Set a configuration value
    Set {
        /// Configuration key (default_worktree_path, auto_cleanup, z_integration, default_sync_strategy, main_branch)
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
        Commands::Status => cmd_status(&manager),
        Commands::Clean { force } => cmd_clean(&manager, force),
        Commands::Sync { worktree, rebase } => cmd_sync(&manager, worktree.as_deref(), rebase),
        Commands::Z { query, list, clean, add } => cmd_z(query.as_deref(), list, clean, add),
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
    
    // Add to z database if z_integration is enabled
    if config.z_integration {
        let z_integration = ZIntegration::new();
        z_integration.create_smart_alias(&target_path, branch)?;
        println!("  Added to z database with smart alias");
    }
    
    Ok(())
}

fn cmd_remove(manager: &WorktreeManager, worktree: &str) -> Result<()> {
    let config = Config::load()?;
    let removed_path;
    
    // Try to find the worktree by name first
    if let Some(wt) = manager.find_worktree_by_name(worktree)? {
        let path_str = wt.path.to_string_lossy();
        removed_path = wt.path.clone();
        manager.remove_worktree(&path_str)?;
        println!("✓ Removed worktree at '{}'", path_str);
    } else {
        // Try direct path removal
        removed_path = std::path::PathBuf::from(worktree);
        manager.remove_worktree(worktree)
            .map_err(|e| anyhow::anyhow!("Worktree '{}' not found. {}", worktree, e))?;
        println!("✓ Removed worktree at '{}'", worktree);
    }
    
    // Remove from z database if z_integration is enabled
    if config.z_integration {
        let z_integration = ZIntegration::new();
        z_integration.remove_worktree(&removed_path)?;
        println!("  Removed from z database");
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
            println!("  default_sync_strategy: {}", config.default_sync_strategy);
            println!("  main_branch: {}", config.main_branch);
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
                "default_sync_strategy" => {
                    if value != "merge" && value != "rebase" {
                        anyhow::bail!("Invalid sync strategy: {}. Valid values: merge, rebase", value);
                    }
                    config.default_sync_strategy = value;
                }
                "main_branch" => {
                    config.main_branch = value;
                }
                _ => {
                    anyhow::bail!("Unknown configuration key: {}. Valid keys: default_worktree_path, auto_cleanup, z_integration, default_sync_strategy, main_branch", key);
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

fn cmd_status(manager: &WorktreeManager) -> Result<()> {
    let worktrees = manager.list_worktrees()?;
    if worktrees.is_empty() {
        println!("No worktrees found.");
        return Ok(());
    }

    println!("{:<30} {:<20} {:<12} {:<15}", "PATH", "BRANCH", "HEAD", "STATUS");
    println!("{}", "-".repeat(80));
    
    for wt in worktrees {
        let path_str = wt.path.to_string_lossy();
        let head_short = if wt.head.len() > 10 {
            &wt.head[..10]
        } else {
            &wt.head
        };
        
        let status = manager.get_worktree_status(&wt.path)?;
        let status_str = if status.is_clean {
            "Clean".to_string()
        } else {
            format!("{}M {}A {}D", status.modified, status.added, status.deleted)
        };
        
        println!("{:<30} {:<20} {:<12} {:<15}", 
            path_str, 
            wt.branch, 
            head_short,
            &status_str
        );
        
        if !status.is_clean {
            if status.modified > 0 {
                println!("  📝 {} modified files", status.modified);
            }
            if status.added > 0 {
                println!("  ➕ {} added files", status.added);
            }
            if status.deleted > 0 {
                println!("  ❌ {} deleted files", status.deleted);
            }
            if status.untracked > 0 {
                println!("  ❓ {} untracked files", status.untracked);
            }
        }
    }
    Ok(())
}

fn cmd_clean(manager: &WorktreeManager, force: bool) -> Result<()> {
    let config = Config::load()?;
    let unnecessary_worktrees = manager.find_unnecessary_worktrees_with_main(&config.main_branch)?;
    
    if unnecessary_worktrees.is_empty() {
        println!("No unnecessary worktrees found.");
        return Ok(());
    }

    println!("Found {} unnecessary worktree(s):", unnecessary_worktrees.len());
    for (wt, reason) in &unnecessary_worktrees {
        println!("  {} - {}", wt.path.display(), reason);
    }

    if !force {
        print!("\nRemove these worktrees? (y/N): ");
        use std::io::{self, Write};
        io::stdout().flush()?;
        
        let mut input = String::new();
        io::stdin().read_line(&mut input)?;
        
        if input.trim().to_lowercase() != "y" {
            println!("Cancelled.");
            return Ok(());
        }
    }

    for (wt, _) in unnecessary_worktrees {
        let path_str = wt.path.to_string_lossy();
        manager.remove_worktree(&path_str)?;
        println!("✓ Removed worktree at '{}'", path_str);
    }

    Ok(())
}

fn cmd_sync(manager: &WorktreeManager, worktree: Option<&str>, rebase: bool) -> Result<()> {
    let config = Config::load()?;
    let target_worktree = if let Some(name) = worktree {
        manager.find_worktree_by_name(name)?
            .ok_or_else(|| anyhow::anyhow!("Worktree '{}' not found", name))?
    } else {
        // Use current directory as worktree
        let current_dir = std::env::current_dir()?;
        let worktrees = manager.list_worktrees()?;
        worktrees.into_iter()
            .find(|wt| wt.path == current_dir)
            .ok_or_else(|| anyhow::anyhow!("Current directory is not a worktree"))?
    };

    // Use rebase flag if provided, otherwise use config default
    let use_rebase = rebase || (config.default_sync_strategy == "rebase");
    let sync_strategy = if use_rebase { "rebase" } else { "merge" };
    println!("Syncing worktree '{}' with {} branch using {}...", 
             target_worktree.branch, config.main_branch, sync_strategy);

    manager.sync_worktree_with_branch(&target_worktree, &config.main_branch, use_rebase)?;
    println!("✓ Successfully synced worktree '{}'", target_worktree.branch);
    
    Ok(())
}

fn cmd_z(query: Option<&str>, list: bool, clean: bool, add: bool) -> Result<()> {
    let z_integration = ZIntegration::new();
    
    if clean {
        let removed_paths = z_integration.cleanup_worktrees()?;
        if removed_paths.is_empty() {
            println!("No stale entries found in z database");
        } else {
            println!("Cleaned up {} stale entries from z database", removed_paths.len());
            for path in removed_paths {
                println!("  Removed: {}", path.display());
            }
        }
        return Ok(());
    }
    
    if add {
        let current_dir = std::env::current_dir()?;
        z_integration.add_worktree(&current_dir, None)?;
        println!("✓ Added current directory to z database: {}", current_dir.display());
        return Ok(());
    }
    
    if let Some(query_str) = query {
        let matches = z_integration.search_worktrees(query_str)?;
        
        if matches.is_empty() {
            eprintln!("No matching worktrees found for: {}", query_str);
            std::process::exit(1);
        }
        
        if list {
            println!("Matching worktrees (sorted by frecency):");
            for (i, entry) in matches.iter().enumerate() {
                let score = entry.frecency_score();
                println!("{:2}: {:8.2} {}", i + 1, score, entry.path.display());
            }
        } else {
            // Jump to the best match
            let best_match = &matches[0];
            println!("{}", best_match.path.display());
        }
    } else {
        // List all entries if no query provided
        let entries = z_integration.read_entries()?;
        if entries.is_empty() {
            println!("No entries in z database");
        } else {
            println!("Z database entries (sorted by frecency):");
            let mut sorted_entries = entries;
            sorted_entries.sort_by(|a, b| b.frecency_score().partial_cmp(&a.frecency_score()).unwrap());
            
            for (i, entry) in sorted_entries.iter().take(20).enumerate() {
                let score = entry.frecency_score();
                println!("{:2}: {:8.2} {}", i + 1, score, entry.path.display());
            }
        }
    }
    
    Ok(())
}