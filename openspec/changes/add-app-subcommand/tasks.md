# Implementation Tasks

## 1. Code Implementation
- [ ] 1.1 Extract TUI launch logic from root Action into shared function `executeTUICommand` in `cli/app.go`
- [ ] 1.2 Create `appCommand()` function that returns `*cli.Command` for the app subcommand
- [ ] 1.3 Update root command Action to call the shared `executeTUICommand` function
- [ ] 1.4 Add app subcommand to Commands list in `NewApp()` function
- [ ] 1.5 Ensure app subcommand inherits log-file flag from parent command

## 2. Testing
- [ ] 2.1 Test `gqlxp app` without arguments opens library selector
- [ ] 2.2 Test `gqlxp app <schema-file>` opens TUI with schema
- [ ] 2.3 Test `gqlxp app --log-file debug.log` enables logging
- [ ] 2.4 Test backward compatibility: `gqlxp` still opens library selector
- [ ] 2.5 Test backward compatibility: `gqlxp <schema-file>` still opens TUI
- [ ] 2.6 Verify existing subcommands (library, search, show, config) still work

## 3. Documentation
- [ ] 3.1 Update README.md to show both `gqlxp` and `gqlxp app` as valid commands
- [ ] 3.2 Update help text if needed to clarify app subcommand purpose

## 4. Quality Checks
- [ ] 4.1 Run `just test` to ensure all tests pass
- [ ] 4.2 Run `just verify` for linting and formatting
- [ ] 4.3 Build and manually test both command forms
