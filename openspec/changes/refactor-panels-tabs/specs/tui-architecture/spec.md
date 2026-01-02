# tui-architecture Specification Deltas

## ADDED Requirements

### Requirement: Tab-Based Type Relationship Navigation
Panels SHALL use a tab-based interface to display different relationships for GraphQL types, making it extensible for future relationship types.

#### Scenario: Field panel with result type and arguments
- **WHEN** a field with arguments is opened in a panel
- **THEN** the panel SHALL display tabs for "Result Type" and "Input Arguments"
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

## REMOVED Requirements

### Requirement: Result Type Special-Case Navigation
**Reason**: Replaced by tab-based navigation which provides a more extensible pattern.

**Migration**: The special-case fields `focusOnResultType`, `resultType`, and related cursor up/down navigation logic will be removed. Instead, result type becomes the content of the first tab ("Result Type") and arguments become the content of the second tab ("Input Arguments"). The `SetObjectType()` method will be replaced with a method to set tab data (labels and their corresponding items).

**Previous behavior**:
- Panel had `resultType` field and `focusOnResultType` bool
- Cursor up from first list item moved focus to result type
- Cursor down from result type moved focus to first list item

**New behavior**:
- Panel has tabs with labels and content
- Shift-H/Shift-L navigate between tabs
- Regular cursor up/down navigate within the active tab's content
