# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`wkit` is a Go CLI tool for convenient Git worktree management with Fish shell integration. It provides easy commands to create, list, remove, and switch between Git worktrees with Fish shell enhancements.

## Core Architecture

- **Go CLI Core**: Main binary built with `cobra` for command parsing
- **Configuration System**: TOML-based configuration with local (`.wkit.toml`) and global (`~/.config/wkit/config.toml`) precedence
- **Worktree Management**: Git worktree wrapper with enhanced functionality in `internal/worktree/manager.go`
- **Fish Integration**: Comprehensive shell integration with functions, aliases, tab completion, and prompt integration

## Development Commands

```bash
# Build the project
go build -o wkit .

# Run tests
go test ./...

# Install to system (after building)
sudo cp wkit /usr/local/bin/

# Install Fish integration
./install.fish
```

## Module Structure

- `main.go`: CLI command definitions and handlers using cobra
- `internal/config/config.go`: Configuration management with TOML serialization/deserialization
- `internal/worktree/manager.go`: Core Git worktree operations and management
- `fish/`: Fish shell integration files (functions, completions, aliases)

## Configuration System

Uses a hierarchical configuration system:
1. Command-line arguments (highest priority)
2. Local `.wkit.toml` in current directory
3. Global `~/.config/wkit/config.toml`
4. Built-in defaults (lowest priority)

Configuration keys: `default_worktree_path`, `auto_cleanup`, `default_sync_strategy`, `main_branch`, `copy_files`

## Fish Integration

The `fish/` directory contains comprehensive Fish shell integration:
- Functions with tab completion for all commands
- Aliases (`ws`, `wl`, `wa`, `wst`) for common operations
- Prompt integration to show current worktree info
- Smart worktree creation and management functions
- Installation script (`install.fish`) for automated setup