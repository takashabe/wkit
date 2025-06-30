# Fish prompt integration for wkit
# Shows current worktree information in your prompt

function __wkit_prompt_info -d "Get current worktree info for prompt"
    # Check if we're in a git repository
    if not git rev-parse --git-dir >/dev/null 2>&1
        return
    end
    
    # Get the current worktree info
    set -l git_dir (git rev-parse --git-dir 2>/dev/null)
    set -l worktree_path (dirname "$git_dir")
    
    # Check if this is a worktree (not the main repo)
    if test -f "$git_dir/gitdir"
        set -l worktree_name (basename "$worktree_path")
        echo -n " ⎇ $worktree_name"
    else if test "$git_dir" = ".git"
        # In main worktree
        echo -n " ⎇ main"
    end
end

function __wkit_prompt_status -d "Get worktree status for prompt"
    # Only run if in a git repository
    if not git rev-parse --git-dir >/dev/null 2>&1
        return
    end
    
    set -l dirty (git status --porcelain 2>/dev/null | wc -l | string trim)
    if test "$dirty" != "0"
        echo -n " ✗"
    end
end

# Example prompt function that includes worktree info
function fish_prompt_with_wkit -d "Example prompt with wkit integration"
    set -l last_status $status
    set -l normal (set_color normal)
    set -l red (set_color red)
    set -l blue (set_color blue)
    set -l green (set_color green)
    set -l yellow (set_color yellow)
    
    # User and host
    echo -n $blue(whoami)$normal@$green(hostname -s)$normal:
    
    # Current directory
    echo -n $yellow(prompt_pwd)$normal
    
    # Git branch if in a git repo
    if git rev-parse --git-dir >/dev/null 2>&1
        set -l branch (git branch --show-current 2>/dev/null)
        if test -n "$branch"
            echo -n " on "$blue$branch$normal
        end
    end
    
    # Worktree info
    echo -n $green(__wkit_prompt_info)$normal
    echo -n $red(__wkit_prompt_status)$normal
    
    # Prompt character
    if test $last_status -eq 0
        echo -n " \$ "
    else
        echo -n $red" \$ "$normal
    end
end

# Function to enable wkit prompt integration
function wkit_prompt_enable -d "Enable wkit prompt integration"
    # Save the current prompt function
    if functions -q fish_prompt
        functions -c fish_prompt __wkit_original_prompt
    end
    
    # Create new prompt that includes wkit info
    function fish_prompt
        set -l last_status $status
        set -l prompt_output
        
        # Call original prompt if it exists
        if functions -q __wkit_original_prompt
            set prompt_output (__wkit_original_prompt)
            # Remove trailing newline if present
            set prompt_output (string trim -r "$prompt_output")
        else
            # Basic prompt if no original
            set prompt_output (prompt_pwd)" \$ "
        end
        
        # Insert wkit info before the final prompt character
        set -l wkit_info (set_color green)(__wkit_prompt_info)(set_color normal)
        set -l wkit_status (set_color red)(__wkit_prompt_status)(set_color normal)
        
        # Find the last line of the prompt
        set -l lines (string split \n "$prompt_output")
        set -l last_line $lines[-1]
        
        # Insert wkit info into the last line
        if string match -q "*\$*" "$last_line"
            set last_line (string replace -r '(.*?)(\$\s*)$' "\$1$wkit_info$wkit_status \$2" "$last_line")
        else if string match -q "*>*" "$last_line"
            set last_line (string replace -r '(.*?)(>\s*)$' "\$1$wkit_info$wkit_status \$2" "$last_line")
        else
            set last_line "$last_line$wkit_info$wkit_status "
        end
        
        # Reconstruct prompt
        if test (count $lines) -gt 1
            for i in (seq 1 (math (count $lines) - 1))
                echo $lines[$i]
            end
        end
        echo -n "$last_line"
        
        # Preserve exit status
        return $last_status
    end
    
    echo "✓ wkit prompt integration enabled"
    echo "  To disable: wkit_prompt_disable"
end

# Function to disable wkit prompt integration
function wkit_prompt_disable -d "Disable wkit prompt integration"
    if functions -q __wkit_original_prompt
        functions -c __wkit_original_prompt fish_prompt
        functions -e __wkit_original_prompt
        echo "✓ wkit prompt integration disabled"
    else
        echo "wkit prompt integration was not enabled"
    end
end