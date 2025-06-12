# Fish shell integration for wkit
# Place this file in ~/.config/fish/functions/ or source it from your config.fish

function wkit-switch
    set -l target_path (wkit switch $argv)
    if test $status -eq 0
        cd "$target_path"
        echo "✓ Switched to worktree: $target_path"
    end
end

# Alias for convenience
alias ws="wkit-switch"

# Tab completion for wkit commands
function __wkit_complete_worktrees
    wkit list 2>/dev/null | tail -n +3 | awk '{print $2}' | grep -v '^$'
end

complete -c wkit -n '__fish_use_subcommand' -a 'list add remove switch' -d 'wkit commands'
complete -c wkit -n '__fish_seen_subcommand_from switch remove' -a '(__wkit_complete_worktrees)' -d 'Available worktrees'