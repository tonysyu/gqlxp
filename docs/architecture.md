# Architecture

For build and development commands, see [Development Commands](development.md).

## Package Structure
- **`cmd/gq`**: Main application entry point that reads schema file and initializes TUI
  - `main()` - Entry point that reads schema and starts TUI
- **`gql`**: GraphQL schema parsing and type extraction
  - `ParseSchema()` - Main function to parse GraphQL schema and extract Query and Mutation fields
  - `GraphQLSchema` struct - Contains Query and Mutation field maps
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
2. Creates list of GraphQL Query and Mutation fields as interactive items
3. User navigates with tab/shift+tab between panels
4. Pressing enter on list items opens details in adjacent panels
5. Supports up to 6 panels horizontally with ctrl+n/ctrl+w for add/remove
6. **Field Type Toggling**: Use ctrl+t to toggle between Query and Mutation fields (see [Features](#features) section below)

## Features

### Query and Mutation Support
The application supports exploring both GraphQL Query and Mutation types:
- **Schema Parsing**: Extracts both Query and Mutation fields from GraphQL schemas
- **Field Type Enum**: Uses `FieldType` enum to track current view mode
- **Dynamic Field Loading**: `loadFieldsPanel()` method loads appropriate fields based on current type

### Interactive Field Type Toggling
Users can switch between viewing Query and Mutation fields:
- **Keyboard Shortcut**: Ctrl+T toggles between Query and Mutation fields
- **Toggle Implementation**: `toggleFieldType()` method switches field types and reloads panel
- **Visual Feedback**: Panel title updates to show current field type

## Dependencies
- **Bubble Tea**: TUI framework for terminal applications
- **Bubbles**: Pre-built components (list, help) for Bubble Tea
- **Lipgloss**: Terminal styling and layout
- **graphql-go**: GraphQL parsing and AST manipulation