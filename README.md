# wkit

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Rust](https://img.shields.io/badge/rust-stable-orange.svg)](https://www.rust-lang.org/)

A powerful, Fish-friendly command-line tool for convenient Git worktree management with enhanced productivity features.

## ✨ Features

- 🌳 **Easy worktree management** - Create, list, remove, and switch between worktrees effortlessly
- ⚙️ **Flexible configuration** - Local and global settings with hierarchical `.wkit.toml` configuration
- 🐟 **Rich Fish shell integration** - Tab completion, aliases, prompt integration, and enhanced functions
- 🔍 **Smart worktree detection** - Find worktrees by branch name or path with intelligent matching
- 📊 **Comprehensive status overview** - See worktree status and git information at a glance
- 🎯 **Z-style navigation** - Quick directory jumping with frecency-based worktree navigation (planned)

## 🚀 Installation

### Quick Install (Recommended)

```bash
# Clone and install in one step
git clone https://github.com/takashabe/wkit.git
cd wkit
./install.sh
```

**What this does:**
- ✅ Builds the binary from source using Cargo
- ✅ Installs to `/usr/local/bin/` (requires sudo)
- ✅ Sets up Fish shell integration with tab completion
- ✅ Creates default configuration files

### Alternative Installation Methods

#### 📦 Using Fisher (Fish Package Manager)

```bash
# Install Fish integration only (binary must exist)
fisher install takashabe/wkit
```

#### 🔨 Manual Build & Install

```bash
# 1. Clone the repository
git clone https://github.com/takashabe/wkit.git
cd wkit

# 2. Build with Cargo
cargo build --release

# 3. Install binary (choose one)
sudo cp target/release/wkit /usr/local/bin/        # System-wide
cp target/release/wkit ~/.local/bin/               # User-only

# 4. Install Fish integration (optional)
./install.fish
```

#### ⚙️ Installation Options

```bash
./install.sh --binary-only     # Skip Fish integration
./install.sh --fish-only       # Fish integration only
./install.sh --help           # Show all options
```

### 📋 Requirements

| Component | Required | Notes |
|-----------|----------|-------|
| **Rust toolchain** | ✅ Yes | Install from [rustup.rs](https://rustup.rs/) |
| **Git** | ✅ Yes | For cloning and git worktree operations |
| **Fish shell** | ⭐ Optional | Recommended for enhanced features |

#### Minimum Versions
- Rust: 1.70+ (2021 edition)
- Fish: 3.0+ (for shell integration)
- Git: 2.25+ (for modern worktree features)

## Usage

### 📖 Basic Commands

| Command | Description | Example |
|---------|-------------|---------|
| `wkit list` | List all worktrees with status | `wkit list` |
| `wkit add <branch>` | Create worktree from existing branch | `wkit add feature-login` |
| `wkit add <branch> [path]` | Create worktree at custom path | `wkit add feature-login ../feature-login` |
| `wkit remove <name>` | Remove a worktree | `wkit remove feature-login` |
| `wkit switch <name>` | Switch to worktree (outputs path) | `wkit switch main` |

```bash
# List all worktrees with detailed status
wkit list

# Create a new worktree from existing branch
wkit add feature-login

# Create worktree at custom location
wkit add feature-login /path/to/custom/location

# Remove worktree when done
wkit remove feature-login

# Switch to different worktree (use with cd)
cd $(wkit switch main)
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

### 🐟 Fish Shell Integration

With Fish integration, you get enhanced productivity features:

#### Quick Aliases
| Alias | Command | Description |
|-------|---------|-------------|
| `ws <name>` | Switch & cd to worktree | `ws feature-login` |
| `wl` | List worktrees | `wl` |
| `wa <branch>` | Add new worktree | `wa new-feature` |
| `wst` | Worktree status overview | `wst` |

#### Advanced Functions
```bash
# Create branch and worktree in one command
wkit-add-quick new-feature main

# Detailed status of all worktrees
wkit-status

# Clean up worktrees with deleted branches
wkit-cleanup

# Prompt integration
wkit_prompt_enable    # Show worktree info in prompt
wkit_prompt_disable   # Remove worktree info from prompt
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

## 🔧 Troubleshooting

### Common Issues

**Q: `wkit` command not found after installation**
```bash
# Ensure /usr/local/bin is in your PATH
echo $PATH | grep "/usr/local/bin"

# If not, add it to your shell config
echo 'set -gx PATH /usr/local/bin $PATH' >> ~/.config/fish/config.fish
```

**Q: Fish integration not working**
```bash
# Verify Fisher installation
fisher list | grep wkit

# Or reinstall manually
./install.fish
```

**Q: Tab completion not working**
```bash
# Reload Fish configuration
source ~/.config/fish/config.fish

# Or restart your terminal
```

**Q: Permission denied when installing**
```bash
# Use sudo for system installation
sudo ./install.sh

# Or install to user directory
cargo install --path . --root ~/.local
```

### Fish Shell Requirements

- Fish shell version 3.0 or higher
- Fisher package manager (recommended but not required)

### Debugging

Enable debug logging to troubleshoot issues:
```bash
RUST_LOG=debug wkit list
```

## 🗺️ Roadmap

- [x] Basic worktree management
- [x] Fish shell integration  
- [x] Configuration management
- [x] Tab completion and aliases
- [ ] Z command integration (frecency-based navigation)
- [ ] Performance optimizations
- [ ] Additional shell support (bash, zsh)
- [ ] Worktree templates
- [ ] Git hooks integration
- [ ] Interactive mode for worktree selection
