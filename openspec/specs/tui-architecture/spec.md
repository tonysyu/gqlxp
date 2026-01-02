# tui-architecture Specification

## Purpose
TBD - created by archiving change restructure-tui-package. Update Purpose after archive.
## Requirements
### Requirement: TUI Package Organization
The TUI package SHALL organize code into subpackages by UI mode to separate concerns and improve maintainability.

#### Scenario: Library selection mode
- **WHEN** code implements library selection functionality
- **THEN** it SHALL be located in `tui/libselect/` subpackage

#### Scenario: Schema exploration mode
- **WHEN** code implements schema exploration functionality
- **THEN** it SHALL be located in `tui/xplr/` subpackage

#### Scenario: Overlay display mode
- **WHEN** code implements overlay display functionality
- **THEN** it SHALL be located in `tui/overlay/` subpackage

#### Scenario: Shared utilities
- **WHEN** code provides shared utilities (adapters, config)
- **THEN** it SHALL remain at the `tui/` package level

#### Scenario: Explorer-specific utilities
- **WHEN** code provides explorer-specific utilities (components, navigation)
- **THEN** it SHALL be located under `tui/xplr/` subpackage

### Requirement: Top-Level Model Delegation
The TUI package SHALL provide a top-level model that delegates to appropriate submode models.

#### Scenario: Mode delegation
- **WHEN** the TUI is started in a specific mode
- **THEN** the top-level model SHALL delegate to the corresponding subpackage model

#### Scenario: Mode transitions
- **WHEN** transitioning between UI modes (e.g., library selection to schema exploration)
- **THEN** the top-level model SHALL handle the transition and activate the appropriate submode

### Requirement: Backward Compatible API
The TUI package SHALL maintain backward compatibility for all public entry points.

#### Scenario: Start function
- **WHEN** `tui.Start()` is called with a schema
- **THEN** it SHALL start the schema exploration mode as before

#### Scenario: StartWithLibraryData function
- **WHEN** `tui.StartWithLibraryData()` is called with schema and metadata
- **THEN** it SHALL start the schema exploration mode with library data as before

#### Scenario: StartSchemaSelector function
- **WHEN** `tui.StartSchemaSelector()` is called
- **THEN** it SHALL start the library selection mode as before

### Requirement: Tab-Based Type Relationship Navigation
Panels SHALL use a tab-based interface to display different relationships for GraphQL types, making it extensible for future relationship types.

#### Scenario: Field panel with result type and arguments
- **WHEN** a field with arguments is opened in a panel
- **THEN** the panel SHALL display tabs for "Type" and "Inputs"
- **AND** the active tab SHALL be visually indicated
- **AND** the content area SHALL show items for the active tab

#### Scenario: Tab navigation with keyboard
- **WHEN** Shift-H is pressed in a panel with multiple tabs
- **THEN** the previous tab SHALL become active
- **AND** the content area SHALL update to show that tab's items

#### Scenario: Tab navigation wrapping
- **WHEN** Shift-L is pressed on the last tab
- **THEN** the active tab SHALL remain on the last tab (no wrapping)

#### Scenario: Single tab display
- **WHEN** a type has only one relationship (e.g., only fields, no result type)
- **THEN** the panel SHALL display only the content without showing tabs
- **OR** SHALL display a single tab if the implementation prefers consistency

#### Scenario: Extensibility for future relationships
- **WHEN** new relationship types are added (e.g., "Implementations", "Interfaces", "Back-references")
- **THEN** they SHALL be added as additional tabs using the same tab navigation pattern
- **AND** SHALL use the same Shift-H/Shift-L keyboard navigation

### Requirement: Panel Tab State Management
Panels SHALL maintain state for the currently active tab and handle focus transitions between tabs and list content.

#### Scenario: Tab and list item focus
- **WHEN** a panel has multiple tabs and list items
- **THEN** keyboard navigation within the list SHALL use existing up/down keys
- **AND** tab switching SHALL use Shift-H/Shift-L
- **AND** focus SHALL remain on the list items, not the tabs themselves

#### Scenario: Active tab persistence
- **WHEN** a panel's active tab is changed
- **THEN** the panel SHALL remember the active tab index
- **AND** SHALL restore it if the panel is revisited (implementation detail)

#### Scenario: Tab content updates
- **WHEN** the active tab changes
- **THEN** the panel's list view SHALL update to show items for the new tab
- **AND** the selection SHALL reset to the first item in the new tab

