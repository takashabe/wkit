# wkit tab completion for Fish shell

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