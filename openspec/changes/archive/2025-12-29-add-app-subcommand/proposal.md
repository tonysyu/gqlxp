# Change: Add Explicit App Subcommand for TUI

## Why
The TUI app functionality is currently split between the root command's default action and a redundant `library` subcommand. This creates confusion and duplicates functionality. Adding a dedicated `app` subcommand and removing `library` provides a clearer, more consistent interface.

## What Changes
- Add new `app` subcommand that launches the TUI (library selector or schema exploration)
- **BREAKING**: Remove the `library` subcommand (users should use `gqlxp app` or `gqlxp` instead)
- Root command continues to route to the same TUI functionality when no subcommand is specified (backward compatible)
- Both `gqlxp` and `gqlxp app` (with no args) open the library selector
- Both `gqlxp <file>` and `gqlxp app <file>` load schema files identically

## Impact
- Affected specs: cli-interface (new capability)
- Affected code:
  - `cli/app.go` - Add new app subcommand, remove library subcommand
- **User impact**: BREAKING - `gqlxp library` command removed (migrate to `gqlxp app` or just `gqlxp`)
