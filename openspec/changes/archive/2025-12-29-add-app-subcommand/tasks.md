# Implementation Tasks

## 1. Code Implementation
- [x] 1.1 Extract TUI launch logic from root Action into shared function `executeTUICommand` in `cli/app.go`
- [x] 1.2 Create `appCommand()` function that returns `*cli.Command` for the app subcommand
- [x] 1.3 Update root command Action to call the shared `executeTUICommand` function
- [x] 1.4 Remove library subcommand from Commands list in `NewApp()` function
- [x] 1.5 Add app subcommand to Commands list in `NewApp()` function
- [x] 1.6 Ensure app subcommand inherits log-file flag from parent command

## 2. Testing
- [x] 2.1 Test `gqlxp app` without arguments opens library selector
- [x] 2.2 Test `gqlxp app <schema-file>` opens TUI with schema
- [x] 2.3 Test `gqlxp app --log-file debug.log` enables logging
- [x] 2.4 Test backward compatibility: `gqlxp` still opens library selector
- [x] 2.5 Test backward compatibility: `gqlxp <schema-file>` still opens TUI
- [x] 2.6 Verify `gqlxp library` command is removed and returns error
- [x] 2.7 Verify remaining subcommands (search, show, config) still work

## 3. Documentation
- [x] 3.1 Update README.md to show both `gqlxp` and `gqlxp app` as valid commands
- [x] 3.2 Update README.md to remove references to `gqlxp library` command
- [x] 3.3 Add migration note for users moving from `gqlxp library` to `gqlxp app`
- [x] 3.4 Update help text if needed to clarify app subcommand purpose

## 4. Quality Checks
- [x] 4.1 Run `just test` to ensure all tests pass
- [x] 4.2 Run `just verify` for linting and formatting
- [x] 4.3 Build and manually test both command forms
