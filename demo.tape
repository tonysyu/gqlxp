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
Type "gqlxp app -s examples/github.graphqls"
Sleep 1s
Enter

# Wait for the TUI to load
Sleep 2s

# Navigate through query types
Type "j"
Sleep 1s
Type "j"
Sleep 1s

# Navigate into type panel
Type "]"
Sleep 2s

# Open overlay for type, scroll down, then close overlay
Space
Sleep 1s
Type "j"
Type "j"
Type "j"
Type "j"
Sleep 1s
Type "q"
Sleep 2s

# Navigate from Type to Inputs, then back
Type "L"
Sleep 2s
Type "H"
Sleep 2s

# Navigate into type's field panel, go to Usages
Type "]"
Sleep 1s
Type "L"
Sleep 1s
Type "L"
Sleep 2s

# Navigate to Object kind
Type "}"
Sleep 1s
Type "}"
Sleep 2s

# Search Object kind for org
Type "/"
Sleep 1s
Type "org"
Sleep 1s
Enter
Sleep 1s
Enter
Sleep 1s
Type "j"
Sleep 1s
Type "j"
Sleep 1s

# Navigate to Search
Type "{"
Sleep 1s
Type "{"
Sleep 1s
Type "{"
Sleep 2s

# Search for org
Type "org"
Sleep 1s
Enter
Sleep 1s
Type "j"
Sleep 1s
Type "j"
Sleep 1s
Type "/"
Sleep 1s
