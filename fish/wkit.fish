# Fish shell integration for wkit
# Place this file in ~/.config/fish/functions/ or source it from your config.fish

function wkit-switch
    set -l target_path (wkit switch $argv)
    if test $status -eq 0
        cd "$target_path"
        echo "✓ Switched to worktree: $target_path"
    end
end

function wkit-checkout
    set -l target_path (wkit checkout $argv)
    if test $status -eq 0
        cd "$target_path"
        echo "✓ Checked out and switched to worktree: $target_path"
    end
end

# Aliases for convenience
alias ws="wkit-switch"
alias wc="wkit-checkout"

# Tab completion for wkit commands
function __wkit_complete_worktrees
    wkit list 2>/dev/null | tail -n +3 | awk '{print $2}' | grep -v '^$'
end

# Tab completion for remote branches
function __wkit_complete_remote_branches
    git branch -r 2>/dev/null | sed 's/^[ \t]*//' | grep -v '^origin/HEAD' | grep -v '^$'
end

complete -c wkit -n '__fish_use_subcommand' -a 'list add remove switch checkout' -d 'wkit commands'
complete -c wkit -n '__fish_seen_subcommand_from switch remove' -a '(__wkit_complete_worktrees)' -d 'Available worktrees'
complete -c wkit -n '__fish_seen_subcommand_from checkout' -a '(__wkit_complete_remote_branches)' -d 'Available remote branches'