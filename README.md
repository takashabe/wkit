# wkit

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Rust](https://img.shields.io/badge/rust-stable-orange.svg)](https://www.rust-lang.org/)

A Fish-friendly CLI tool for convenient Git worktree management.

## Features

- **Worktree management** - Create, list, remove, and switch between Git worktrees
- **Flexible configuration** - Local (`.wkit.toml`) and global (`~/.config/wkit/config.toml`) configuration with XDG support
- **Fish shell integration** - Tab completion, aliases, and prompt integration
- **Batch operations** - Clean up multiple worktrees efficiently
- **Sync with main branch** - Keep worktrees up-to-date with merge or rebase
- **Z-style navigation** - Frecency-based worktree jumping

## Installation

### From Source

```bash
# Clone and build
git clone https://github.com/takashabe/wkit.git
cd wkit
cargo build --release

# Install binary
sudo cp target/release/wkit /usr/local/bin/

# Install Fish integration (optional)
./install.fish
```

### Requirements

- Rust toolchain (1.70+)
- Git (2.25+)
- Fish shell (3.0+) - optional but recommended

## Usage

### Basic Commands

```bash
# List all worktrees
wkit list

# Add a new worktree
wkit add feature-branch
wkit add feature-branch custom/path

# Remove a worktree
wkit remove feature-branch

# Switch to a worktree (outputs path)
cd $(wkit switch main)

# Show status of all worktrees
wkit status

# Clean up worktrees
wkit clean

# Sync worktree with main branch
wkit sync                    # current worktree
wkit sync feature-branch     # specific worktree
wkit sync --rebase          # use rebase instead of merge

# Z-style jumping
wkit z proj                 # jump to most frecent match
wkit z --list              # list all entries
wkit z --add               # add current directory
```

### Configuration

```bash
# Show current configuration
wkit config show

# Set configuration values
wkit config set default_worktree_path .git/.wkit-worktrees
wkit config set auto_cleanup true
wkit config set main_branch main

# Create local configuration file
wkit config init
```

### Fish Shell Integration

With Fish integration installed, you get:

**Aliases:**
- `ws <name>` - Switch to worktree
- `wl` - List worktrees
- `wa <branch>` - Add worktree
- `wst` - Show status

**Functions:**
- `wkit-add-quick` - Create branch and worktree
- `wkit-status` - Detailed status view
- `wkit-cleanup` - Clean up deleted branches
- `wkit_prompt_enable/disable` - Toggle prompt integration

## Configuration

Configuration files:
- Local: `.wkit.toml` (project-specific)
- Global: `~/.config/wkit/config.toml` (or `$XDG_CONFIG_HOME/wkit/config.toml`)

### Options

```toml
# Default path for new worktrees
default_worktree_path = ".git/.wkit-worktrees"

# Automatically clean up deleted branches
auto_cleanup = false

# Enable z integration
z_integration = true

# Default sync strategy (merge or rebase)
default_sync_strategy = "merge"

# Main branch name
main_branch = "main"

# Copy files to new worktrees
[copy_files]
enabled = false
files = [".envrc", ".env.local", "compose.override.yaml"]
```

### Precedence

1. Command-line arguments
2. Local `.wkit.toml`
3. Global config
4. Built-in defaults

## License

MIT License - see [LICENSE](LICENSE) file for details.
