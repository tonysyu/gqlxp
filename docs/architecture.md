# Architecture

## Package Structure
- **`cmd/gq`**: Main application entry point that reads schema file and initializes TUI
- **`gql`**: GraphQL schema parsing and type extraction
  - `ParseSchema()` - Main function to parse GraphQL schema and extract Query fields
  - `buildGraphQLTypes()` - Builds map of GraphQL types from AST
  - `GetTypeString()` - Converts AST types to string representation
- **`tui`**: Terminal user interface components built on Bubble Tea
  - `mainModel` - Root application model managing panels and navigation
  - `Panel` interface - Generic panel abstraction for different content types
  - `ListItem` interface - Interactive list items that can be "opened" to show details
  - `listPanel` and `stringPanel` - Concrete panel implementations

## Key Interfaces
- **`ListItem`**: Extends `list.DefaultItem` with `Open() Panel` method for interactive items
- **`Panel`**: Implements `tea.Model` with `SetSize(width, height int)` for resizable content

## Navigation Flow
1. Application starts by parsing GraphQL schema from `examples/github.graphqls`
2. Creates list of GraphQL query fields as interactive items
3. User navigates with tab/shift+tab between panels
4. Pressing enter on list items opens details in adjacent panels
5. Supports up to 6 panels horizontally with ctrl+n/ctrl+w for add/remove

## Dependencies
- **Bubble Tea**: TUI framework for terminal applications
- **Bubbles**: Pre-built components (list, help) for Bubble Tea
- **Lipgloss**: Terminal styling and layout
- **graphql-go**: GraphQL parsing and AST manipulation