# Architecture

For build and development commands, see [Development Commands](development.md).

## Package Structure
- **`cmd/igq`**: Main application entry point
  - `main()` - Reads schema file and starts TUI
- **`gql`**: GraphQL schema parsing and type extraction
  - `ParseSchema()` - Parses GraphQL schema and extracts all type definitions
  - `GraphQLSchema` struct - Contains maps for Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, and Directive types
  - `NamedToTypeDefinition()` - Resolves type names to their definitions
  - Helper functions: `GetTypeString()`, `CollectAndSortMapValues()`, `GetFieldDefinitionString()`, etc.
- **`tui`**: Terminal user interface built on Bubble Tea
  - `mainModel` - Root model managing panels, navigation, and GQL type toggling
  - **`tui/components`**: Reusable UI components
    - `Panel` interface - Generic panel abstraction
    - `ListItem` interface - Interactive list items with `Open()` method
    - `ListPanel` - Panel displaying lists with auto-open behavior
    - `SimpleItem` - Basic ListItem implementation
  - **`tui/adapters`**: Converts GraphQL AST types to UI components
    - `Adapt*ToItems()` functions - Convert schema types to ListItems
    - `fieldItem`, `typeDefItem` - ListItem adapters for GraphQL types

## Key Interfaces
- **`ListItem`**: Extends `list.DefaultItem` with `Open() Panel` method for interactive items
- **`Panel`**: Implements `tea.Model` with `SetSize(width, height int)` for resizable content

## Dependencies
- **Bubble Tea**: TUI framework for terminal applications
- **Bubbles**: Pre-built components (list, help) for Bubble Tea
- **Lipgloss**: Terminal styling and layout
- **graphql-go**: GraphQL parsing and AST manipulation
