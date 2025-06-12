#!/usr/bin/env fish

# wkit Fish Shell Integration Installer
# This script helps you set up wkit with Fish shell

set -l SCRIPT_DIR (dirname (realpath (status -f)))
set -l FISH_CONFIG_DIR ~/.config/fish
set -l FISH_FUNCTIONS_DIR $FISH_CONFIG_DIR/functions

function print_header
    echo
    echo "🐟 wkit Fish Shell Integration Installer"
    echo "========================================"
    echo
end

function print_step
    echo "📦 $argv"
end

function print_success
    echo "✅ $argv"
end

function print_error
    echo "❌ $argv" >&2
end

function print_info
    echo "ℹ️  $argv"
end

function check_requirements
    print_step "Checking requirements..."
    
    # Check if Fish is being used
    if not test "$SHELL" = (which fish)
        print_error "Fish shell is not your default shell"
        echo "Please run: chsh -s "(which fish)
        return 1
    end
    
    # Check if wkit binary exists
    if not command -q wkit
        print_error "wkit binary not found in PATH"
        echo "Please build and install wkit first:"
        echo "  cargo build --release"
        echo "  cp target/release/wkit /usr/local/bin/"
        return 1
    end
    
    print_success "Requirements check passed"
end

function create_directories
    print_step "Creating Fish configuration directories..."
    
    if not test -d $FISH_CONFIG_DIR
        mkdir -p $FISH_CONFIG_DIR
        print_success "Created $FISH_CONFIG_DIR"
    end
    
    if not test -d $FISH_FUNCTIONS_DIR
        mkdir -p $FISH_FUNCTIONS_DIR
        print_success "Created $FISH_FUNCTIONS_DIR"
    end
end

function install_functions
    print_step "Installing Fish functions..."
    
    # Copy main functions
    if test -f $SCRIPT_DIR/fish/functions/wkit.fish
        cp $SCRIPT_DIR/fish/functions/wkit.fish $FISH_FUNCTIONS_DIR/
        print_success "Installed wkit.fish functions"
    else
        print_error "wkit.fish not found in $SCRIPT_DIR/fish/functions/"
        return 1
    end
    
    # Copy prompt functions
    if test -f $SCRIPT_DIR/fish/functions/wkit_prompt.fish
        cp $SCRIPT_DIR/fish/functions/wkit_prompt.fish $FISH_FUNCTIONS_DIR/
        print_success "Installed wkit_prompt.fish functions"
    else
        print_error "wkit_prompt.fish not found in $SCRIPT_DIR/fish/functions/"
        return 1
    end
end

function setup_config
    print_step "Setting up Fish configuration..."
    
    set -l config_file $FISH_CONFIG_DIR/config.fish
    set -l wkit_config "# wkit configuration
# Load wkit functions automatically
if command -q wkit
    # wkit functions are loaded automatically from functions directory
    echo \"wkit Fish integration loaded. Type 'wkit --help' to get started.\"
end"
    
    if not test -f $config_file
        echo $wkit_config > $config_file
        print_success "Created config.fish with wkit setup"
    else
        if not grep -q "wkit" $config_file
            echo >> $config_file
            echo $wkit_config >> $config_file
            print_success "Added wkit configuration to existing config.fish"
        else
            print_info "wkit configuration already exists in config.fish"
        end
    end
end

function show_usage
    print_step "Installation complete! Here's how to use wkit:"
    echo
    echo "🎯 Basic Commands:"
    echo "  wkit list              - List all worktrees"
    echo "  wkit add <branch>      - Add new worktree"
    echo "  wkit remove <worktree> - Remove worktree"
    echo "  wkit switch <worktree> - Switch to worktree"
    echo "  wkit config show       - Show configuration"
    echo
    echo "🚀 Fish Integration:"
    echo "  ws <worktree>          - Quick switch alias"
    echo "  wl                     - Quick list alias"
    echo "  wa <branch>            - Quick add alias"
    echo "  wst                    - Show worktree status"
    echo "  wkit-add-quick <name>  - Quickly create worktree from current branch"
    echo
    echo "🎨 Prompt Integration:"
    echo "  wkit_prompt_enable     - Show worktree info in prompt"
    echo "  wkit_prompt_disable    - Hide worktree info from prompt"
    echo
    echo "⚙️  Configuration:"
    echo "  wkit config init       - Create local .wkit.toml"
    echo "  wkit config set default_worktree_path ../worktrees"
    echo
    echo "📚 Tab Completion:"
    echo "  All commands have tab completion for worktrees, branches, and config keys"
    echo
    echo "🔄 To reload: exec fish"
end

function main
    print_header
    
    if not check_requirements
        exit 1
    end
    
    create_directories
    
    if not install_functions
        exit 1
    end
    
    setup_config
    show_usage
    
    echo
    print_success "Installation successful!"
    print_info "Restart your shell or run 'exec fish' to start using wkit Fish integration"
end

# Run installer
main $argv