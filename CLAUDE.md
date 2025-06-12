# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`wkit` is a Rust CLI tool for convenient Git worktree management with Fish shell integration. It provides easy commands to create, list, remove, and switch between Git worktrees, along with optional z-style frecency-based navigation and Fish shell enhancements.

## Core Architecture

- **Rust CLI Core**: Main binary built with `clap` for command parsing, uses `anyhow` for error handling
- **Configuration System**: TOML-based configuration with local (`.wkit.toml`) and global (`~/.config/wkit/config.toml`) precedence
- **Worktree Management**: Git worktree wrapper with enhanced functionality in `src/worktree.rs`
- **Z Integration**: Frecency-based directory jumping similar to the `z` command in `src/z_integration.rs`
- **Fish Integration**: Comprehensive shell integration with functions, aliases, tab completion, and prompt integration

## Development Commands

```bash
# Build the project
cargo build --release

# Run tests
cargo test

# Run with debug logging
RUST_LOG=debug cargo run -- list

# Install to system (after building)
sudo cp target/release/wkit /usr/local/bin/

# Install Fish integration
./install.fish
```

## Module Structure

- `src/main.rs`: CLI command definitions and handlers using clap
- `src/config.rs`: Configuration management with TOML serialization/deserialization
- `src/worktree.rs`: Core Git worktree operations and management
- `src/z_integration.rs`: Frecency-based worktree navigation system
- `fish/`: Fish shell integration files (functions, completions, aliases)
- `tests/integration_tests.rs`: Integration tests using `assert_cmd`

## Configuration System

Uses a hierarchical configuration system:
1. Command-line arguments (highest priority)
2. Local `.wkit.toml` in current directory
3. Global `~/.config/wkit/config.toml`
4. Built-in defaults (lowest priority)

Configuration keys: `default_worktree_path`, `auto_cleanup`, `z_integration`

## Fish Integration

The `fish/` directory contains comprehensive Fish shell integration:
- Functions with tab completion for all commands
- Aliases (`ws`, `wl`, `wa`, `wst`) for common operations
- Prompt integration to show current worktree info
- Smart worktree creation and management functions
- Installation script (`install.fish`) for automated setup