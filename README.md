# wkit - Convenient Git Worktree Management Toolkit

A Fish-friendly command-line tool for managing Git worktrees with enhanced productivity features.

## Features

- 🌳 **Easy worktree management** - Create, list, remove, and switch between worktrees
- ⚙️ **Flexible configuration** - Local and global settings with `.wkit.toml`
- 🐟 **Fish shell integration** - Tab completion, aliases, and prompt integration
- 🔍 **Smart worktree detection** - Find worktrees by branch name or path
- 📊 **Status overview** - See worktree status and git information at a glance

## Installation

### 🚀 Quick Install (Recommended)

```bash
# Clone and run automatic installer
git clone https://github.com/takashabe/wkit.git
cd wkit
./install.sh
```

This will:
- Build the binary from source
- Install to `/usr/local/bin/`
- Set up Fish shell integration with tab completion and aliases
- Create default configuration

### 📦 Manual Installation

#### Build from Source

```bash
# Clone the repository
git clone https://github.com/takashabe/wkit.git
cd wkit

# Build the binary
cargo build --release

# Install to system PATH
sudo cp target/release/wkit /usr/local/bin/

# Install Fish integration (optional but recommended)
./install.fish
```

#### Fish Shell Integration Only

If you already have wkit installed and just want Fish integration:

```bash
git clone https://github.com/takashabe/wkit.git
cd wkit
./install.fish
```

### 🔧 Installation Options

The installation script supports various options:

```bash
# Install only binary (skip Fish integration)
./install.sh --binary-only

# Install only Fish integration (assume binary exists)
./install.sh --fish-only

# Show all available options
./install.sh --help
```

### 📋 Requirements

- **Rust toolchain** - Install from [rustup.rs](https://rustup.rs/)
- **Git** - For cloning the repository
- **Fish shell** - Optional, for enhanced shell integration

## Usage

### Basic Commands

```bash
# List all worktrees
wkit list

# Add a new worktree
wkit add feature-branch
wkit add feature-branch /custom/path

# Remove a worktree
wkit remove feature-branch

# Switch to a worktree (outputs path for shell integration)
wkit switch feature-branch
```

### Configuration Management

```bash
# Show current configuration
wkit config show

# Set default worktree path
wkit config set default_worktree_path ../worktrees

# Create local configuration file
wkit config init
```

### Fish Shell Integration

If you installed the Fish integration, you get these extra features:

```bash
# Quick aliases
ws feature-branch    # Switch to worktree
wl                   # List worktrees
wa new-feature       # Add worktree
wst                  # Worktree status overview

# Advanced functions
wkit-add-quick new-feature main  # Create branch and worktree from main
wkit-status                      # Detailed status of all worktrees
wkit-cleanup                     # Remove worktrees with deleted branches

# Prompt integration
wkit_prompt_enable               # Show worktree info in prompt
wkit_prompt_disable              # Remove worktree info from prompt
```

### Tab Completion

With Fish integration, you get comprehensive tab completion for:
- All wkit commands and subcommands
- Worktree names for `switch` and `remove`
- Branch names for `add`
- Configuration keys for `config set`
- File paths where appropriate

## Configuration

wkit supports both local (`.wkit.toml`) and global (`~/.config/wkit/config.toml`) configuration files.

### Configuration Options

```toml
# Default path for new worktrees (relative to current directory)
default_worktree_path = ".."

# Automatically clean up deleted branches
auto_cleanup = false

# Enable z integration (planned feature)
z_integration = true
```

### Configuration Precedence

1. Command-line path argument (highest priority)
2. Local `.wkit.toml` in current directory
3. Global `~/.config/wkit/config.toml`
4. Built-in defaults (lowest priority)

## Examples

### Typical Workflow

```bash
# Create a worktree for a new feature
wkit add feature-login

# Switch to it (with Fish integration)
ws feature-login

# Work on your feature...
git commit -m "Add login functionality"

# Switch back to main
ws main

# Remove the feature worktree when done
wkit remove feature-login
```

### With Custom Configuration

```bash
# Set up custom worktree directory
wkit config set default_worktree_path ~/Dev/worktrees

# Now all new worktrees go to ~/Dev/worktrees/
wkit add feature-payment  # Creates ~/Dev/worktrees/feature-payment
```

## Architecture

- **Rust CLI** - Fast, reliable core functionality with proper error handling
- **Fish Functions** - Enhanced shell integration and convenience features
- **TOML Configuration** - Flexible, human-readable settings
- **Git Integration** - Direct git worktree command integration

## Development

```bash
# Run tests
cargo test

# Check code
cargo check

# Run with debug output
RUST_LOG=debug cargo run -- list
```

## Contributing

1. Fork the repository
2. Create a feature branch: `wkit add my-feature`
3. Make your changes and add tests
4. Ensure tests pass: `cargo test`
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Roadmap

- [x] Basic worktree management
- [x] Fish shell integration
- [x] Configuration management
- [ ] Z command integration
- [ ] Performance optimizations
- [ ] Additional shell support (bash, zsh)
- [ ] Worktree templates
- [ ] Git hooks integration
