# TUI Search Tab Capability

## ADDED Requirements

### Requirement: Search Tab Type
The TUI SHALL include "Search" as a navigable GraphQL type tab, positioned after "Directive" in the type cycle.

#### Scenario: Search in type cycle
- **WHEN** cycling through types with Ctrl+T from "Directive"
- **THEN** the next type SHALL be "Search"

#### Scenario: Reverse cycle to Search
- **WHEN** cycling backwards with Ctrl+R from "Query"
- **THEN** the previous type SHALL be "Search"

#### Scenario: Search type constant
- **WHEN** the Search tab type is defined in code
- **THEN** it SHALL be added as a `SearchType` constant in `navigation.GQLType`

### Requirement: Search Input Display
The TUI SHALL display a text input field at the bottom of the screen when the Search tab is active.

#### Scenario: Input field visibility
- **WHEN** the Search tab is active
- **THEN** a text input field SHALL be visible at the bottom of the screen
- **AND** the input SHALL use `charmbracelet/bubbles/textinput` component

#### Scenario: Input field positioning
- **WHEN** the Search tab is active
- **THEN** the input field SHALL appear above the help bar
- **AND** the main panel height SHALL be adjusted to accommodate the input

#### Scenario: Input field hidden on other tabs
- **WHEN** navigating to any non-Search tab
- **THEN** the text input field SHALL be hidden
- **AND** the main panel SHALL use full available height

### Requirement: Search Query Execution
The TUI SHALL execute search queries using the existing `search` package and Bleve indexes when the user submits input.

#### Scenario: Execute search on enter
- **WHEN** the user types a query and presses Enter in the search input
- **THEN** the search SHALL execute using `search.Searcher.Search()`
- **AND** results SHALL be displayed in the main panel

#### Scenario: Empty query handling
- **WHEN** the search input is empty and Enter is pressed
- **THEN** the main panel SHALL display an empty state message
- **AND** no search query SHALL be executed

#### Scenario: Search with missing index
- **WHEN** a search is executed but the schema index doesn't exist
- **THEN** indexing SHALL start in the background
- **AND** a "Indexing schema..." message SHALL be displayed
- **AND** search SHALL execute once indexing completes

### Requirement: Search Results Display
The TUI SHALL display search results as a list of selectable items in the main panel.

#### Scenario: Result list format
- **WHEN** search results are displayed
- **THEN** each result SHALL show the path (e.g., "User.email")
- **AND** the type category (e.g., "Field")
- **AND** the description excerpt if available

#### Scenario: Empty results
- **WHEN** a search returns no results
- **THEN** the panel SHALL display "No results found for '<query>'"

#### Scenario: Result selection navigation
- **WHEN** a user selects a search result
- **THEN** the TUI SHALL navigate to the appropriate type tab
- **AND** SHALL select the type and field as specified in the result path
- **AND** SHALL open detail panels as in `xplr.ApplySelection()`

### Requirement: Search Input Focus Management
The TUI SHALL manage focus between the search input field and the results panel.

#### Scenario: Initial focus on input
- **WHEN** the Search tab becomes active
- **THEN** keyboard focus SHALL be on the search input field
- **AND** typing SHALL update the input text

#### Scenario: Focus on results after search
- **WHEN** a search is executed with Enter
- **THEN** focus SHALL move to the results panel
- **AND** arrow keys SHALL navigate the result list

#### Scenario: Return focus to input
- **WHEN** focus is on results panel and user presses "/"
- **THEN** focus SHALL return to the search input field
- **AND** existing input text SHALL be preserved

#### Scenario: Clear input with Escape
- **WHEN** focus is on the input field and Escape is pressed
- **THEN** the input field SHALL be cleared
- **AND** focus SHALL remain on the input field
