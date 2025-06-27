# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`wkit` is a CLI tool for convenient Git worktree management with Fish shell integration. It provides easy commands to create, list, remove, and switch between Git worktrees, along with optional z-style frecency-based navigation and Fish shell enhancements.

The project has both Go and Rust implementations. The Go implementation is the current focus for active development.

## Core Architecture

### Go Implementation
- **CLI Core**: Main binary built with `cobra` for command parsing
- **Configuration System**: TOML-based configuration with local (`.wkit.toml`) and global (`~/.config/wkit/config.toml`) precedence
- **Worktree Management**: Git worktree wrapper with enhanced functionality in `internal/worktree/`
- **Git Operations**: Unified Git command executor in `internal/git/`
- **Command Modules**: Modular command structure in `internal/cmd/`

### Legacy Rust Implementation
- **Rust CLI Core**: Main binary built with `clap` for command parsing, uses `anyhow` for error handling
- **Z Integration**: Frecency-based directory jumping similar to the `z` command in `src/z_integration.rs`

## Development Commands

### Go Implementation (Primary)
```bash
# Build the project
go build -o wkit .

# Run tests
go test ./...

# Run with debug logging
./wkit --help

# Install to system (after building)
sudo cp wkit /usr/local/bin/

# Install Fish integration
./install.fish
```

### Rust Implementation (Legacy)
```bash
# Build the project
cargo build --release

# Run tests
cargo test

# Run with debug logging
RUST_LOG=debug cargo run -- list

# Install to system (after building)
sudo cp target/release/wkit /usr/local/bin/
```

## Module Structure

### Go Implementation
- `main.go`: CLI root command and initialization
- `internal/cmd/`: Individual command implementations (add, list, remove, etc.)
- `internal/config/`: Configuration management with TOML serialization/deserialization
- `internal/worktree/`: Core Git worktree operations and management
- `internal/git/`: Unified Git command executor with error handling
- `integration_test.go`: Integration tests for CLI functionality

### Rust Implementation (Legacy)
- `src/main.rs`: CLI command definitions and handlers using clap
- `src/config.rs`: Configuration management with TOML serialization/deserialization
- `src/worktree.rs`: Core Git worktree operations and management
- `src/z_integration.rs`: Frecency-based worktree navigation system
- `tests/integration_tests.rs`: Integration tests using `assert_cmd`

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