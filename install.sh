#!/usr/bin/env bash

# wkit Installation Script
# Downloads or builds and installs wkit binary and Fish shell integration

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_header() {
    echo -e "\n${BLUE}🔧 wkit Installation Script${NC}"
    echo "=================================="
    echo
}

print_step() {
    echo -e "${YELLOW}📦 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}" >&2
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case "$os" in
        darwin)
            case "$arch" in
                x86_64) PLATFORM="x86_64-apple-darwin" ;;
                arm64) PLATFORM="aarch64-apple-darwin" ;;
                *) print_error "Unsupported architecture: $arch"; exit 1 ;;
            esac
            BINARY_NAME="wkit"
            ARCHIVE_EXT="tar.gz"
            ;;
        linux)
            case "$arch" in
                x86_64) PLATFORM="x86_64-unknown-linux-gnu" ;;
                *) print_error "Unsupported architecture: $arch"; exit 1 ;;
            esac
            BINARY_NAME="wkit"
            ARCHIVE_EXT="tar.gz"
            ;;
        mingw*|msys*|cygwin*)
            case "$arch" in
                x86_64) PLATFORM="x86_64-pc-windows-msvc" ;;
                *) print_error "Unsupported architecture: $arch"; exit 1 ;;
            esac
            BINARY_NAME="wkit.exe"
            ARCHIVE_EXT="zip"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
    
    print_info "Detected platform: $PLATFORM"
}

# Download binary from GitHub releases
download_binary() {
    print_step "Downloading wkit binary from GitHub releases..."
    
    local repo="takashabe/wkit"
    local asset_name="wkit-${PLATFORM}.${ARCHIVE_EXT}"
    local download_url="https://github.com/${repo}/releases/latest/download/${asset_name}"
    local temp_dir=$(mktemp -d)
    local temp_file="${temp_dir}/${asset_name}"
    
    print_info "Downloading from: $download_url"
    
    if command_exists curl; then
        curl -L -o "$temp_file" "$download_url"
    elif command_exists wget; then
        wget -O "$temp_file" "$download_url"
    else
        print_error "Neither curl nor wget found. Please install one of them."
        return 1
    fi
    
    if [ ! -f "$temp_file" ]; then
        print_error "Download failed"
        return 1
    fi
    
    # Extract archive
    cd "$temp_dir"
    case "$ARCHIVE_EXT" in
        tar.gz)
            tar -xzf "$asset_name"
            ;;
        zip)
            if command_exists unzip; then
                unzip "$asset_name"
            else
                print_error "unzip not found. Please install unzip."
                return 1
            fi
            ;;
    esac
    
    if [ ! -f "$BINARY_NAME" ]; then
        print_error "Binary not found in downloaded archive"
        return 1
    fi
    
    # Make binary executable
    chmod +x "$BINARY_NAME"
    
    # Move to a location accessible by install_downloaded_binary
    DOWNLOADED_BINARY_PATH="${temp_dir}/${BINARY_NAME}"
    
    print_success "Binary downloaded successfully"
}

# Check requirements for building from source
check_build_requirements() {
    print_step "Checking build requirements..."
    
    local missing_deps=()
    
    if ! command_exists cargo; then
        missing_deps+=("cargo (Rust toolchain)")
    fi
    
    if ! command_exists git; then
        missing_deps+=("git")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_error "Missing required dependencies for building from source:"
        for dep in "${missing_deps[@]}"; do
            echo "  - $dep"
        done
        echo
        echo "Please install the missing dependencies and try again."
        echo "To install Rust: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
        return 1
    fi
    
    print_success "Build requirements satisfied"
    return 0
}

# Check basic requirements (for download method)
check_basic_requirements() {
    print_step "Checking basic requirements..."
    
    if ! command_exists curl && ! command_exists wget; then
        print_error "Neither curl nor wget found. Please install one of them."
        return 1
    fi
    
    print_success "Basic requirements satisfied"
    return 0
}

# Build the binary
build_binary() {
    print_step "Building wkit binary..."
    
    if [ ! -f "Cargo.toml" ]; then
        print_error "Cargo.toml not found. Please run this script from the wkit project directory."
        exit 1
    fi
    
    cargo build --release
    if [ $? -eq 0 ]; then
        print_success "Binary built successfully"
    else
        print_error "Failed to build binary"
        exit 1
    fi
}

# Install binary (from build)
install_binary() {
    print_step "Installing wkit binary..."
    
    local binary_path="target/release/wkit"
    local install_dir="/usr/local/bin"
    
    if [ ! -f "$binary_path" ]; then
        print_error "Binary not found at $binary_path"
        exit 1
    fi
    
    # Check if we can write to install directory
    if [ -w "$install_dir" ]; then
        cp "$binary_path" "$install_dir/"
    else
        print_info "Installing to $install_dir requires sudo privileges"
        sudo cp "$binary_path" "$install_dir/"
    fi
    
    # Verify installation
    if command_exists wkit; then
        print_success "Binary installed to $install_dir/wkit"
    else
        print_error "Installation verification failed"
        exit 1
    fi
}

# Install downloaded binary
install_downloaded_binary() {
    print_step "Installing downloaded wkit binary..."
    
    local install_dir="/usr/local/bin"
    local target_name="wkit"
    
    if [ ! -f "$DOWNLOADED_BINARY_PATH" ]; then
        print_error "Downloaded binary not found at $DOWNLOADED_BINARY_PATH"
        exit 1
    fi
    
    # Check if we can write to install directory
    if [ -w "$install_dir" ]; then
        cp "$DOWNLOADED_BINARY_PATH" "$install_dir/$target_name"
    else
        print_info "Installing to $install_dir requires sudo privileges"
        sudo cp "$DOWNLOADED_BINARY_PATH" "$install_dir/$target_name"
    fi
    
    # Verify installation
    if command_exists wkit; then
        print_success "Binary installed to $install_dir/$target_name"
    else
        print_error "Installation verification failed"
        exit 1
    fi
}

# Install Fish integration
install_fish_integration() {
    if ! command_exists fish; then
        print_info "Fish shell not found, skipping Fish integration"
        echo "To install Fish shell integration later, run: ./install.fish"
        return 0
    fi
    
    print_step "Installing Fish shell integration..."
    
    if [ -f "install.fish" ]; then
        # Make sure install.fish is executable
        chmod +x install.fish
        
        # Run Fish installer
        if fish install.fish; then
            print_success "Fish integration installed"
        else
            print_error "Failed to install Fish integration"
            return 1
        fi
    else
        print_error "install.fish not found"
        return 1
    fi
}

# Create configuration
create_default_config() {
    print_step "Creating default configuration..."
    
    local config_dir="$HOME/.config/wkit"
    local config_file="$config_dir/config.toml"
    
    if [ ! -d "$config_dir" ]; then
        mkdir -p "$config_dir"
    fi
    
    if [ ! -f "$config_file" ]; then
        cat > "$config_file" << 'EOF'
# wkit configuration file
# Generated by install script

default_worktree_path = ".."
auto_cleanup = false
z_integration = true
default_sync_strategy = "merge"
main_branch = "main"
EOF
        print_success "Default configuration created at $config_file"
    else
        print_info "Configuration file already exists at $config_file"
    fi
}

# Show usage information
show_usage() {
    print_step "Installation complete! Here's how to use wkit:"
    echo
    echo "🎯 Basic Commands:"
    echo "  wkit list              - List all worktrees"
    echo "  wkit add <branch>      - Add new worktree"
    echo "  wkit remove <worktree> - Remove worktree"
    echo "  wkit switch <worktree> - Switch to worktree"
    echo "  wkit status            - Show git status of all worktrees"
    echo "  wkit clean [--force]   - Clean up unnecessary worktrees"
    echo "  wkit sync [worktree]   - Sync with main branch"
    echo "  wkit config show       - Show configuration"
    echo
    if command_exists fish; then
        echo "🚀 Fish Integration:"
        echo "  ws <worktree>          - Quick switch alias"
        echo "  wl                     - Quick list alias"
        echo "  wa <branch>            - Quick add alias"
        echo "  wst                    - Show worktree status"
        echo "  wsy [worktree]         - Quick sync alias"
        echo "  wcl                    - Quick clean alias"
        echo
        echo "🔄 To activate Fish integration: exec fish"
    fi
    echo
    echo "📚 Documentation: https://github.com/takashabe/wkit"
    echo "🐛 Issues: https://github.com/takashabe/wkit/issues"
}

# Main installation process
main() {
    print_header
    
    # Build from source
    if check_build_requirements; then
        build_binary
        install_binary
    else
        print_error "Build requirements are not met. Please install Rust and try again."
        echo "To install Rust: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
        exit 1
    fi
    
    create_default_config
    
    if command_exists fish; then
        install_fish_integration
    fi
    
    show_usage
    
    echo
    print_success "Installation successful!"
    if command_exists fish; then
        print_info "Restart your shell or run 'exec fish' to start using wkit"
    else
        print_info "You can now use wkit from any directory"
    fi
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --help|-h)
            echo "wkit Installation Script"
            echo
            echo "Usage: $0 [options]"
            echo
            echo "Options:"
            echo "  --help, -h        Show this help message"
            echo "  --binary-only     Install only the binary (skip Fish integration)"
            echo "  --fish-only       Install only Fish integration (assume binary exists)"
            echo
            exit 0
            ;;
        --binary-only)
            BINARY_ONLY=1
            shift
            ;;
        --fish-only)
            FISH_ONLY=1
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Run installation based on options
if [ "${FISH_ONLY}" = "1" ]; then
    print_header
    install_fish_integration
elif [ "${BINARY_ONLY}" = "1" ]; then
    print_header
    check_build_requirements
    build_binary
    install_binary
    create_default_config
    show_usage
else
    main
fi