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

# Note: Aliases and completions are handled by conf.d/wkit.fish and completions/wkit.fish

# Enhanced add function with auto-switch capability
function wkit-add-auto -d "Add worktree with optional auto-switch"
    if test (count $argv) -eq 0
        echo "Error: Please specify a branch name" >&2
        return 1
    end
    
    set -l wkit_output (wkit add $argv 2>&1)
    set -l exit_code $status
    
    echo "$wkit_output"
    
    if test $exit_code -eq 0
        # Extract the last word which should be the path
        set -l words (echo "$wkit_output" | string split ' ')
        set -l potential_path $words[-1]
        
        # If the last word looks like a path, switch to it
        if test -d "$potential_path"
            cd "$potential_path"
            echo "✓ Automatically switched to: "(basename "$potential_path")" at $potential_path"
        end
    end
    
    return $exit_code
end

# Checkout existing branch function with auto-switch capability
function wkit-checkout -d "Checkout existing branch and create worktree"
    if test (count $argv) -eq 0
        echo "Error: Please specify a branch (local or remote, e.g., feature-branch, origin/feature-branch)" >&2
        return 1
    end
    
    set -l wkit_output (wkit checkout $argv 2>&1)
    set -l exit_code $status
    
    echo "$wkit_output"
    
    if test $exit_code -eq 0
        # Extract the last word which should be the path
        set -l words (echo "$wkit_output" | string split ' ')
        set -l potential_path $words[-1]
        
        # If the last word looks like a path, switch to it
        if test -d "$potential_path"
            cd "$potential_path"
            echo "✓ Automatically switched to: "(basename "$potential_path")" at $potential_path"
        end
    end
    
    return $exit_code
end

# Note: Abbreviations are handled by conf.d/wkit.fish