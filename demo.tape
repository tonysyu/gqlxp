# VHS tape file for demonstrating gqlxp
# Run with: vhs demo.tape

Output demo.gif

Set FontSize 18
Set Width 1200
Set Height 1000
Set Theme "Catppuccin Mocha"

# Show the help command
Type "gqlxp --help"
Sleep 1s
Enter
Sleep 3s

# Clear screen
Type "clear"
Enter
Sleep 1s

# Run gqlxp on an example schema
Type "gqlxp examples/github.graphqls"
Sleep 1s
Enter

# Wait for the TUI to load
Sleep 2s

# Navigate through the schema
Type "j"
Sleep 1s
Type "j"
Sleep 1s

# Open overlay
Space
Sleep 2s

# Quit overlay
Type "q"
Sleep 1s

# Navigate into and out of child type panels
Type "l"
Sleep 1s
Type "]"
Sleep 1s
Type "k"
Sleep 1s
Type "]"
Sleep 1s
Type "j"
Sleep 1s
Type "]"
Sleep 2s
Type "["
Sleep 1s
Type "["
Sleep 1s
Type "["
Sleep 1s

# Navigate through types
Type "}"
Sleep 1s
Type "}"
Sleep 1s
Type "}"
Sleep 1s
Type "}"
Sleep 1s
Type "}"
Sleep 1s
Type "}"
Sleep 1s
Type "{"
Sleep 1s
