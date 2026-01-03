# Proposal: Add Interfaces Subtab

## Overview
Add an "Interfaces" subtab to Object type panels to display implemented interfaces, enabling users to navigate from Objects to their interface definitions.

## Problem
Currently, when viewing an Object type that implements interfaces (e.g., `User implements Node & Timestamped`), there is no way to see or navigate to those interface definitions. This limits schema exploration for interface-based patterns common in GraphQL.

## Proposed Solution
Extend the existing tab-based panel navigation to include an "Interfaces" tab for Object types showing all interfaces the object implements.

The Object type already exposes `Interfaces()` returning `[]string` of interface names. The TypeResolver can resolve these names to full Interface definitions for navigation.

## Scope
This change affects:
- `tui/adapters/items.go`: Modify `typeDefItem.OpenPanel()` to add "Interfaces" tab for Object types
- Tests: Update adapter tests to verify new tabs are created correctly

## Dependencies
- Existing tab navigation system (already implemented in panels)
- `gql.Object.Interfaces()` method (already available)
- TypeResolver for resolving interface names to TypeDefs

## Success Criteria
1. Object type panels display an "Interfaces" tab when the object implements one or more interfaces
2. Users can navigate from Objects to their Interface definitions
3. Tab navigation follows existing Shift+H/Shift+L keyboard patterns
4. All tests pass with new functionality

## Related Changes
This change aligns with the TUI architecture spec's extensibility goal for relationship types (Requirement: Tab-Based Type Relationship Navigation, Scenario: Extensibility for future relationships).
