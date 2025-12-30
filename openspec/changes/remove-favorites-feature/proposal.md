# Proposal: Remove Favorites Feature

## Problem
The favorites feature adds complexity to the schema library and TUI without clear value:

1. **Limited utility**: Favoriting types in a GraphQL schema doesn't provide significant value over search and navigation
2. **Maintenance burden**: Requires persistent storage, UI interactions, keybindings, and metadata management
3. **Cognitive overhead**: Adds another concept users must learn without solving a real workflow problem

## Solution
Remove the favorites feature entirely, including:
- Favorites field from schema metadata storage
- Favorites keybinding and toggle logic from TUI
- Favoritable item wrapper and visual indicators
- Library methods for managing favorites

## Scope
This change affects:
- `openspec/specs/schema-library/` - Remove favorites requirements
- `library/` - Remove favorites storage and methods
- `tui/` - Remove favorites UI integration
- `tui/xplr/` - Remove favoritable item wrapper and keybinding

## Dependencies
None - this is a standalone removal of an isolated feature.

## Risks
- **Breaking change**: Users with existing favorites in metadata.json will lose that data (field will be ignored)
- **No migration needed**: Old metadata.json files with favorites field will continue to load (Go's JSON unmarshaling ignores unknown fields)
- Low impact: Favorites feature is not essential to core functionality
