# wkit Fish shell integration initialization

# Check if wkit binary is available
if command -q wkit
    # Set up aliases
    alias ws="wkit-switch"
    alias wl="wkit list"
    alias wa="wkit add"
    alias wr="wkit remove"
    alias wst="wkit-status"
    alias wsy="wkit-sync"
    alias wcl="wkit-clean"
    
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