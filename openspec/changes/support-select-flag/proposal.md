# Proposal: Support Select Flag

## Overview
Add a `--select` flag to the `gqlxp app` command that opens the TUI with a specific type or field pre-selected. This allows users to jump directly to a type or field of interest instead of manually navigating through the TUI.

## Motivation
Users working with large GraphQL schemas often need to repeatedly navigate to specific types or fields during development. Currently, they must:
1. Launch the TUI
2. Cycle through type categories (Query, Mutation, Object, etc.)
3. Search or scroll through the list to find their target
4. For fields, navigate into the type and search again

This proposal enables direct navigation via CLI argument: `gqlxp app schema.graphqls --select Query.user` or `gqlxp app schema.graphqls --select User`.

## Scope
This change affects:
- **CLI**: Add `--select` flag to app command with validation
- **TUI Initialization**: Support starting with pre-selected type/field
- **Navigation**: Add logic to locate and select items by name

## User Impact
- **Positive**: Faster workflow for developers repeatedly accessing specific types/fields
- **No breaking changes**: Existing behavior unchanged when flag not used
- **Error handling**: Invalid selections fall back to normal TUI behavior

## Implementation Strategy
1. Parse `--select` argument format (Type or Type.field)
2. Pass selection target to TUI initialization
3. Determine which GQL type category contains the target
4. Locate item in panel list and select it
5. For Type.field format, navigate forward and select field

## Related Specs
- `cli-interface`: CLI command structure and flags
- `tui-architecture`: TUI initialization and navigation

## Success Criteria
- `gqlxp app schema.graphqls --select TypeName` opens with TypeName selected
- `gqlxp app schema.graphqls --select TypeName.fieldName` opens with field selected
- Invalid selections gracefully fall back to default behavior
- No impact on existing functionality when flag not used
