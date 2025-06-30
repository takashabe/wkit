# Fish Shell Integration Examples

This directory contains example Fish shell functions and configurations for integrating `wkit` into your Fish workflow.

## Installation

Copy the functions you want to use to your Fish configuration:

```bash
# Copy all functions
cp functions/*.fish ~/.config/fish/functions/

# Or copy specific functions
cp functions/ws.fish ~/.config/fish/functions/
```

## Available Functions

### `ws` - Switch worktree with fuzzy search
Switches to a different worktree using fzf for selection.

```fish
ws
```

### `wl` - List worktrees
Lists all worktrees in a formatted table.

```fish
wl
```

### `wa` - Add new worktree
Creates a new worktree. Automatically uses the wkit configuration.

```fish
wa feature/new-feature
```

### `wst` - Show worktree status
Shows the status of all worktrees.

```fish
wst
```

## Customization Examples

### Using JSON output with jq

```fish
function my_worktree_list
    wkit list --format=json | jq -r '.[] | "\(.branch)\t\(.path)"'
end
```

### Custom worktree switcher with preview

```fish
function ws_preview
    wkit list --format=json | \
    jq -r '.[] | "\(.branch)\t\(.path)"' | \
    fzf --preview 'echo {} | cut -f2 | xargs -I {} ls -la {}' | \
    cut -f2 | \
    read -l selected
    and cd $selected
end
```

### Worktree prompt integration

```fish
function fish_prompt
    # Your existing prompt code...
    
    # Add worktree info
    set -l worktree_info (wkit list --format=json 2>/dev/null | jq -r '.[] | select(.path == "(root)" or .path == ".") | .branch' 2>/dev/null)
    if test -n "$worktree_info"
        echo -n " [$worktree_info]"
    end
    
    # Rest of your prompt...
end
```

## Tips

- The functions use `wkit` commands internally, so ensure `wkit` is in your PATH
- Customize the functions to match your workflow preferences
- Consider adding abbreviations for frequently used commands:
  ```fish
  abbr -a wt 'wkit'
  ```