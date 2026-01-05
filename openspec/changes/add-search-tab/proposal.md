# Proposal: Add Search Tab

## Overview
Add a new "Search" tab to the TUI explorer (`tui/xplr/`) to enable full-text search within the loaded schema, similar to the CLI search functionality but integrated into the interactive interface.

## Why
Users currently have to exit the TUI and use `gqlxp search` to find types and fields. Integrating search directly into the TUI provides:
- Seamless workflow without context switching
- Interactive result navigation within the same interface
- Consistent user experience across CLI and TUI

## Scope
- Add "Search" as a new `GQLType` constant after "Directive"
- Create search input component using `charmbracelet/bubbles/textinput`
- Display search results in the first panel
- Integrate with existing `search` package (Bleve indexing)
- Maintain consistent tab navigation patterns (Ctrl+T, Ctrl+R for type switching)

## Out of Scope
- Modifying search indexing logic (reuses existing `search` package)
- Advanced search features (filters, operators beyond existing Bleve capabilities)
- Search history or saved searches

## Related Specs
- `schema-search` - Existing search indexing and query functionality
- `tui-architecture` - TUI structure and tab-based navigation
- `interface-navigation` - Type selector and navigation patterns

## Risks and Mitigation
- **Risk**: Search indexing may block TUI responsiveness
  - **Mitigation**: Reuse existing background indexing; show loading state while index builds
- **Risk**: Input component conflicts with existing keyboard shortcuts
  - **Mitigation**: Search input only receives input when Search tab is active
