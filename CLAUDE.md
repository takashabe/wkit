# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`wkit` is a Go CLI tool for convenient Git worktree management with Fish shell integration. It provides easy commands to create, list, remove, and switch between Git worktrees with Fish shell enhancements.

## Core Architecture

- **Go CLI Core**: Main binary built with `cobra` for command parsing
- **Configuration System**: YAML-based configuration with local (`.wkit.yaml`) and global (`~/.config/wkit/config.yaml`) precedence
- **Worktree Management**: Git worktree wrapper with enhanced functionality in `internal/worktree/manager.go`
- **Structured Output**: Supports JSON output format for easy scripting and integration (e.g., `wkit list --format=json`)
- **Fish Integration Examples**: Example functions and configurations in `examples/fish/` directory

## Development Commands

```bash
# Build the project
go build -o wkit .

# Run all tests
go test ./...

# Run specific test file
go test ./integration_test.go

# Run tests with verbose output
go test -v ./...

# Code formatting and linting
go fmt ./...
gofmt -s -l .
go vet ./...

# Install to system (after building)
sudo cp wkit /usr/local/bin/

# Copy Fish integration examples (optional)
cp examples/fish/functions/*.fish ~/.config/fish/functions/
```

## Module Structure

- `main.go`: CLI command definitions and handlers using cobra
- `internal/config/config.go`: Configuration management with YAML serialization/deserialization
- `internal/worktree/manager.go`: Core Git worktree operations and management
- `internal/cmd/`: Command implementations with JSON output support
- `examples/fish/`: Fish shell integration examples (functions, completions, aliases)

## Configuration System

Uses a hierarchical configuration system:
1. Command-line arguments (highest priority)
2. Local `.wkit.yaml` in current directory
3. Global `~/.config/wkit/config.yaml`
4. Built-in defaults (lowest priority)

Configuration keys: `wkit_root`, `auto_cleanup`, `default_sync_strategy`, `main_branch`, `copy_files`

## Shell Integration

The `examples/fish/` directory contains example Fish shell integrations that users can customize:
- Functions with tab completion for all commands
- Aliases (`ws`, `wl`, `wa`, `wst`) for common operations
- Prompt integration to show current worktree info
- Smart worktree creation and management functions
- Examples of using JSON output with tools like jq and fzf

## Testing Strategy

### Integration Tests
- Located in `integration_test.go` at the project root
- Tests CLI binary by building and executing actual commands
- Requires Git repository environment
- Skipped if Git is not available

### Unit Tests
- Located in `internal/` package subdirectories as `*_test.go`
- Tests configuration management, worktree operations, and command logic
- Run with `go test ./internal/...` for package-specific testing

## CI/CD and Release Management

### Continuous Integration (`.github/workflows/ci.yml`)
- Automated testing on Go 1.19-1.22 across Linux, macOS, and Windows
- Cross-compilation testing for multiple architectures (amd64, arm64)
- Code quality checks including `go vet` and formatting validation

### Release Automation
- **tagpr** (`tagpr.yml`): Automated version management using Conventional Commits
- **GoReleaser** (`.goreleaser.yml`): Cross-platform binary distribution
- **Automated releases**: Triggered on tag pushes to create GitHub releases with binaries

### Development Workflow
- Conventional Commits enforced for consistent release notes
- Automatic CHANGELOG.md generation
- Cross-platform testing ensures compatibility