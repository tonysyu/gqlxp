# Change: Add Explicit App Subcommand for TUI

## Why
The TUI app functionality is currently only accessible through the root command's default action, making the command structure less explicit and harder to discover. Adding a dedicated `app` subcommand provides a clearer interface while maintaining backward compatibility with the default behavior.

## What Changes
- Add new `app` subcommand that mirrors the current default action (launches TUI)
- Root command continues to route to the same TUI functionality when no subcommand is specified
- Both `gqlxp` and `gqlxp app` launch the TUI with identical behavior
- Both `gqlxp <file>` and `gqlxp app <file>` load schema files identically

## Impact
- Affected specs: cli-interface (new capability)
- Affected code:
  - `cli/app.go` - Add new app subcommand
  - `cli/app_cmd.go` (new) - Extract TUI launch logic to reusable function
- User impact: None (backward compatible, additive change only)
