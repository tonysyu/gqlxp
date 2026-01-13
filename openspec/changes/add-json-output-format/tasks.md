# Implementation Tasks

## 1. JSON Generation Module
- [ ] 1.1 Create `gqlfmt/json.go` with `GenerateJSON` function
- [ ] 1.2 Implement JSON struct types for Fields, TypeDefs, and Directives
- [ ] 1.3 Add helpers for converting GraphQL types to JSON structs
- [ ] 1.4 Write unit tests for JSON generation with various type examples

## 2. CLI Integration
- [ ] 2.1 Add `--json` boolean flag to show command in `cli/show.go`
- [ ] 2.2 Update `printType` function to handle JSON output path
- [ ] 2.3 Auto-disable pager when `--json` flag is set
- [ ] 2.4 Ensure JSON outputs to stdout (no markdown rendering)

## 3. Testing & Validation
- [ ] 3.1 Add integration tests for `gqlxp show --json` command
- [ ] 3.2 Test JSON output for all type categories (Query, Mutation, Object, Enum, etc.)
- [ ] 3.3 Verify pretty-printing with proper indentation
- [ ] 3.4 Run `just test` to ensure all tests pass

## 4. Documentation
- [ ] 4.1 Update show command usage text to mention `--json` flag
- [ ] 4.2 Add JSON output examples to command description
