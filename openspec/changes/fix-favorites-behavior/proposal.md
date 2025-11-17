# Proposal: Fix Favorites Behavior

## Problem
The favorites feature has several behavioral issues that make it confusing and inconsistent:

1. **Incorrect favoriting target**: In top-level panels (Query, Mutation, etc.), favoriting a field stores the return type instead of the field itself
2. **Inconsistent selection behavior**: After toggling favorite, the selected item doesn't remain stable during panel refresh
3. **Missing visual feedback**: Panels displaying favorited types/fields don't show any indication that the content is favorited

## Solution
Implement context-aware favoriting that:
- Stores field names (via `RefName()`) when used in top-level GQL type panels
- Stores type names when used in other panels
- Preserves selection after favorite toggle
- Displays favorite indicator (★) in panel titles for favorited content

## Scope
This change modifies the TUI favorites interaction and display logic within `tui/xplr/`. No changes to library storage format are required.

## Dependencies
None - this is a standalone fix to existing favorites functionality.

## Risks
- Low: Changes are isolated to favorites toggle and display logic
- Breaking change: Existing favorites may need re-favoriting since stored values will change from type names to field names for Query/Mutation fields
