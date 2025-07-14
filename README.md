# wkit

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)

A CLI tool for convenient Git worktree management.

## Features

- **Worktree management** - Create, list, remove, and switch between Git worktrees
- **Flexible configuration** - Local (`.wkit.yaml`) and global (`~/.config/wkit/config.yaml`) configuration with XDG support
- **Batch operations** - Clean up multiple worktrees efficiently
- **Sync with main branch** - Keep worktrees up-to-date with merge or rebase

## Installation

### From Source

```bash
# Clone and build
git clone https://github.com/takashabe/wkit.git
cd wkit
go build -o wkit .

# Install binary
sudo cp wkit /usr/local/bin/
```

### Shell Integration (Optional)

See [examples/fish/](examples/fish/) for shell integration examples that you can customize and add to your configuration.

### Requirements

- Go (1.21+)
- Git (2.25+)

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

```

### Configuration

```bash
# Show current configuration
wkit config show

# Set configuration values
wkit config set wkit_root .git/.wkit-worktrees
wkit config set auto_cleanup true
wkit config set main_branch main

# Create local configuration file
wkit config init
```

### Structured Output

`wkit` supports JSON output for easy integration with scripts and tools:

```bash
# Get worktree list as JSON
wkit list --format=json
```

### Shell Integration Examples

See [examples/fish/](examples/fish/) for shell integration examples including:

- Fuzzy worktree switching with fzf
- Custom aliases and functions
- Prompt integration
- JSON output parsing examples

## Configuration

Configuration files:
- Local: `.wkit.yaml` (project-specific)
- Global: `~/.config/wkit/config.yaml` (or `$XDG_CONFIG_HOME/wkit/config.yaml`)

### Options

```yaml
# Default path for new worktrees
wkit_root: ".git/.wkit-worktrees"

# Automatically clean up deleted branches
auto_cleanup: false

# Default sync strategy (merge or rebase)
default_sync_strategy: "merge"

# Main branch name
main_branch: "main"

# Copy files to new worktrees
copy_files:
  enabled: false
  files:
    - ".envrc"
    - ".env.local"
    - "compose.override.yaml"
```

### Precedence

1. Command-line arguments
2. Local `.wkit.yaml`
3. Global config
4. Built-in defaults

## License

MIT License - see [LICENSE](LICENSE) file for details.
