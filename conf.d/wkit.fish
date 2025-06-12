# wkit Fish shell integration initialization

# Check if wkit binary is available
if command -q wkit
    # Set up primary aliases (safe names)
    alias wkit-add="wkit-add-auto"
    alias wkit-sw="wkit-switch"
    alias wkit-ls="wkit list"
    alias wkit-rm="wkit remove"
    alias wkit-st="wkit-status"
    alias wkit-sy="wkit-sync"
    alias wkit-cl="wkit-clean"
    
    # Optional short aliases (set WKIT_SHORT_ALIASES=1 to enable)
    if set -q WKIT_SHORT_ALIASES
        alias wa="wkit-add-auto"
        alias ws="wkit-switch"
        alias wl="wkit list"
        alias wr="wkit remove"
        alias wst="wkit-status"
        alias wsy="wkit-sync"
        alias wcl="wkit-clean"
    end
    
    # Set up abbreviations for common workflows
    abbr -a wsc 'wkit switch (wkit list | tail -n +3 | fzf | awk "{print \$2}")'
    abbr -a wzl 'wkit z --list'
    abbr -a wzc 'wkit z --clean'
else
    echo "wkit binary not found. Please install wkit first:"
    echo "  git clone https://github.com/takashabe/wkit.git"
    echo "  cd wkit"
    echo "  ./install.sh"
end