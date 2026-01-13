# Implementation Tasks

## 1. JSON Generation Module
- [x] 1.1 Create `gqlfmt/json.go` with `GenerateJSON` function
- [x] 1.2 Implement JSON struct types for Fields, TypeDefs, and Directives
- [x] 1.3 Add helpers for converting GraphQL types to JSON structs
- [x] 1.4 Write unit tests for JSON generation with various type examples
- [x] 1.5 Add `kind` field to all JSON outputs (Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, Directive)

## 2. CLI Integration
- [x] 2.1 Add `--json` boolean flag to show command in `cli/show.go`
- [x] 2.2 Update `printType` function to handle JSON output path
- [x] 2.3 Auto-disable pager when `--json` flag is set
- [x] 2.4 Ensure JSON outputs to stdout (no markdown rendering)

## 3. Testing & Validation
- [x] 3.1 Add integration tests for `gqlxp show --json` command
- [x] 3.2 Test JSON output for all type categories (Query, Mutation, Object, Enum, etc.)
- [x] 3.3 Verify pretty-printing with proper indentation
- [x] 3.4 Run `just test` to ensure all tests pass

## 4. Documentation
- [x] 4.1 Update show command usage text to mention `--json` flag
- [x] 4.2 Add JSON output examples to command description
