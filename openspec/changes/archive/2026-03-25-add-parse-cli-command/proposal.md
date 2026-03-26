## Why

Developers need to validate GraphQL operations against a schema from the terminal without launching the interactive TUI. A `parse` command enables scripting, CI validation, and quick sanity checks.

## What Changes

- New `parse` subcommand: `gqlxp parse [--schema/-s <schema>] [<operation-file>]`
- Reads a GraphQL operation from a file path argument or stdin
- Validates the operation against the resolved schema
- Reports errors with location info (line/column), or indicates success
- Exits non-zero on validation errors (script-friendly)

## Capabilities

### New Capabilities
- `query-validation`: Validates a GraphQL operation (query/mutation/subscription) against a schema and reports syntax errors, unknown fields, and type violations

### Modified Capabilities

## Impact

- New CLI command registered in the existing CLI app
- Reuses existing schema resolution (`--schema/-s` flag pattern) and `gqlparser/v2` (already a dependency)
