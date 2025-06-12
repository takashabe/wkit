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
        # Check if the last line is a path (auto-switch output)
        set -l last_line (echo "$wkit_output" | tail -n 1)
        
        # If last line looks like a path, switch to it
        if test -d "$last_line"
            cd "$last_line"
            echo "✓ Automatically switched to: "(basename "$last_line")" at $last_line"
        end
    end
    
    return $exit_code
end

# Note: Abbreviations are handled by conf.d/wkit.fish