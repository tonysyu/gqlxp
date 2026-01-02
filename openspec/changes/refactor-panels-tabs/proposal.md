# Change: Refactor panels to use tab-based navigation for type relationships

## Why
The current panel implementation special-cases `resultType` with dedicated navigation logic (`focusOnResultType`, cursor up/down handling). This makes it difficult to extend with additional relationships like Implementations, Interfaces, back-references, etc. A tab-based interface provides a cleaner, more extensible pattern for navigating between different aspects of a GraphQL type.

## What Changes
- Replace the special-case `resultType` handling with a tab-based interface similar to the Bubble Tea tabs example
- Support navigating between tabs (e.g., "Result Type", "Input Arguments") using Shift-H (left) and Shift-L (right)
- Make the architecture extensible to support future relationships (Implementations, Interfaces implemented, back-references)
- Remove `focusOnResultType`, `resultType` field, and related cursor navigation logic from Panel struct
- Add tab state management (active tab index, tab labels, tab content)

## Impact
- Affected specs: `tui-architecture`
- Affected code:
  - `tui/xplr/components/panels.go` - Complete rewrite of result type and tab handling
  - `tui/xplr/components/panels_test.go` - Update tests for new tab-based navigation
  - `tui/adapters/items.go` - Update `OpenPanel()` to provide tab data instead of single result type
  - `tui/xplr/model.go` - Add Shift-H/Shift-L keybindings for tab navigation
