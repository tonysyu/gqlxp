# Change: Standardize CLI Schema Flags

## Why
Current CLI commands (`search`, `show`, `app`) use positional arguments for schema selection, making the interface ambiguous and harder to discover. Users must know the argument order to distinguish between schema identifiers and query parameters.

## What Changes
- Replace positional `[schema-file]` argument with explicit `--schema/-s` flag
- Add `--default/-d` flag to explicitly use the default schema
- Remove implicit schema resolution from positional arguments
- **BREAKING**: Commands require explicit schema selection via flags

## Impact
- Affected specs: `cli-interface`
- Affected code: `cli/app.go`, `cli/search.go`, `cli/show.go`
- Breaking change for existing users who rely on positional schema arguments
