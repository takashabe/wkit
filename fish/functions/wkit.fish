# Fish shell integration for wkit
# Enhanced version with better completion and error handling

# Main switch function with error handling
function wkit-switch -d "Switch to a worktree"
    if test (count $argv) -eq 0
        echo "Error: Please specify a worktree name" >&2
        return 1
    end
    
    set -l target_path (wkit switch $argv 2>&1)
    set -l exit_code $status
    
    if test $exit_code -eq 0
        cd "$target_path"
        echo "✓ Switched to worktree: "(basename "$target_path")" at $target_path"
    else
        echo "Error: $target_path" >&2
        return $exit_code
    end
end

# List worktrees with preview
function wkit-list-preview -d "List worktrees with status preview"
    set -l current_dir (pwd)
    
    for worktree in (wkit list 2>/dev/null | tail -n +3)
        set -l path (echo $worktree | awk '{print $1}')
        set -l branch (echo $worktree | awk '{print $2}')
        
        if test -d "$path"
            cd "$path" 2>/dev/null
            set -l status_info (git status --porcelain 2>/dev/null | wc -l | string trim)
            set -l ahead_behind (git rev-list --left-right --count HEAD...@{u} 2>/dev/null)
            
            echo -n "$branch: "
            if test -n "$status_info" -a "$status_info" != "0"
                echo -n "[$status_info uncommitted] "
            end
            if test -n "$ahead_behind"
                echo "$ahead_behind" | read ahead behind
                if test "$ahead" != "0" -o "$behind" != "0"
                    echo -n "[↑$ahead ↓$behind] "
                end
            end
            echo "$path"
        end
    end
    
    cd "$current_dir"
end

# Quick add worktree
function wkit-add-quick -d "Quickly add a worktree from current branch"
    if test (count $argv) -eq 0
        echo "Usage: wkit-add-quick <new-branch-name> [base-branch]" >&2
        return 1
    end
    
    set -l new_branch $argv[1]
    set -l base_branch (test (count $argv) -ge 2; and echo $argv[2]; or echo "HEAD")
    
    # Create new branch and worktree
    git worktree add -b "$new_branch" "../$new_branch" "$base_branch"
    if test $status -eq 0
        echo "✓ Created worktree '$new_branch' based on $base_branch"
        echo "  Use 'ws $new_branch' to switch to it"
    end
end

# Worktree status function
function wkit-status -d "Show status of all worktrees"
    wkit status
end

# Sync worktree function
function wkit-sync -d "Sync current or specified worktree with main branch"
    if test (count $argv) -eq 0
        wkit sync
    else
        wkit sync $argv
    end
end

# Clean worktrees function
function wkit-clean -d "Clean up unnecessary worktrees interactively"
    wkit clean
end

# Clean up function
function wkit-cleanup -d "Remove worktrees with deleted branches"
    for worktree in (wkit list 2>/dev/null | tail -n +3 | awk '{print $2}')
        if not git show-ref --verify --quiet "refs/heads/$worktree" 2>/dev/null
            echo "Branch '$worktree' no longer exists. Remove worktree? (y/N)"
            read -l response
            if test "$response" = "y"
                wkit remove "$worktree"
            end
        end
    end
end

# Convenience aliases
alias ws="wkit-switch"
alias wl="wkit list"
alias wa="wkit add"
alias wr="wkit remove"
alias wst="wkit-status"
alias wsy="wkit-sync"
alias wcl="wkit-clean"

# Enhanced tab completion
complete -c wkit -f

# Subcommand completions
complete -c wkit -n '__fish_use_subcommand' -a 'list' -d 'List all worktrees'
complete -c wkit -n '__fish_use_subcommand' -a 'add' -d 'Add a new worktree'
complete -c wkit -n '__fish_use_subcommand' -a 'remove' -d 'Remove a worktree'
complete -c wkit -n '__fish_use_subcommand' -a 'switch' -d 'Switch to a worktree'
complete -c wkit -n '__fish_use_subcommand' -a 'config' -d 'Manage configuration'
complete -c wkit -n '__fish_use_subcommand' -a 'status' -d 'Show git status of all worktrees'
complete -c wkit -n '__fish_use_subcommand' -a 'clean' -d 'Clean up unnecessary worktrees'
complete -c wkit -n '__fish_use_subcommand' -a 'sync' -d 'Sync worktree with main branch'
complete -c wkit -n '__fish_use_subcommand' -a 'z' -d 'Frecency-based worktree jumping'

# Config subcommand completions
complete -c wkit -n '__fish_seen_subcommand_from config; and __fish_use_subcommand' -a 'show' -d 'Show current configuration'
complete -c wkit -n '__fish_seen_subcommand_from config; and __fish_use_subcommand' -a 'set' -d 'Set configuration value'
complete -c wkit -n '__fish_seen_subcommand_from config; and __fish_use_subcommand' -a 'init' -d 'Initialize local config'

# Config set key completions
complete -c wkit -n '__fish_seen_subcommand_from config; and __fish_seen_subcommand_from set' -a 'default_worktree_path' -d 'Default path for new worktrees'
complete -c wkit -n '__fish_seen_subcommand_from config; and __fish_seen_subcommand_from set' -a 'auto_cleanup' -d 'Auto cleanup deleted branches'
complete -c wkit -n '__fish_seen_subcommand_from config; and __fish_seen_subcommand_from set' -a 'z_integration' -d 'Enable z integration'
complete -c wkit -n '__fish_seen_subcommand_from config; and __fish_seen_subcommand_from set' -a 'default_sync_strategy' -d 'Default sync strategy (merge/rebase)'
complete -c wkit -n '__fish_seen_subcommand_from config; and __fish_seen_subcommand_from set' -a 'main_branch' -d 'Main branch name'

# Worktree name completions for switch and remove
function __wkit_complete_worktrees
    wkit list 2>/dev/null | tail -n +3 | awk '{print $2"\t"$1}' | grep -v '^$'
end

complete -c wkit -n '__fish_seen_subcommand_from switch remove' -f -a '(__wkit_complete_worktrees)'

# Branch completions for add command
complete -c wkit -n '__fish_seen_subcommand_from add' -f -a "(git branch -a 2>/dev/null | sed 's/^[* ]*//' | sed 's/^remotes\///' | sort -u)"

# Path argument for add command (as second argument)
complete -c wkit -n '__fish_seen_subcommand_from add; and test (count (commandline -opc)) -eq 3' -f -a "(__fish_complete_directories)"

# Clean command options
complete -c wkit -n '__fish_seen_subcommand_from clean' -s f -l force -d 'Skip confirmation prompt'

# Sync command options
complete -c wkit -n '__fish_seen_subcommand_from sync' -s r -l rebase -d 'Use rebase instead of merge'
complete -c wkit -n '__fish_seen_subcommand_from sync' -f -a '(__wkit_complete_worktrees)' -d 'Worktree to sync'

# Z command options
complete -c wkit -n '__fish_seen_subcommand_from z' -s l -l list -d 'List all matches instead of jumping'
complete -c wkit -n '__fish_seen_subcommand_from z' -s c -l clean -d 'Clean up non-existent entries'
complete -c wkit -n '__fish_seen_subcommand_from z' -s a -l add -d 'Add current directory to z database'

# Z command completions for wz function
complete -c wz -f -a "(wkit z --list 2>/dev/null | tail -n +2 | awk '{print \$3}' | xargs -I {} basename {} 2>/dev/null)"

# Z-style jumping function
function wz -d "Jump to worktree using frecency (z-style)"
    if test (count $argv) -eq 0
        wkit z
        return
    end
    
    set -l target_path (wkit z $argv)
    set -l exit_code $status
    
    if test $exit_code -eq 0 -a -n "$target_path"
        cd "$target_path"
        echo "✓ Jumped to: "(basename "$target_path")" at $target_path"
    else
        return $exit_code
    end
end

# Enhanced z function that integrates with wkit
function z-wkit -d "Enhanced z function with wkit integration"
    # If wkit z finds a match, use it; otherwise fall back to regular z
    set -l wkit_result (wkit z $argv 2>/dev/null)
    if test $status -eq 0 -a -n "$wkit_result"
        cd "$wkit_result"
        echo "✓ wkit: Jumped to "(basename "$wkit_result")
    else
        # Fall back to regular z command
        command z $argv
    end
end

# Abbreviations for common workflows
abbr -a wsc 'wkit switch (wkit list | tail -n +3 | fzf | awk "{print \$2}")'
abbr -a wzl 'wkit z --list'
abbr -a wzc 'wkit z --clean'