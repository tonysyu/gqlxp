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

## Navigation Flow
1. Application parses GraphQL schema from provided file path
2. Displays Query fields by default in main panel
3. Selecting items auto-opens details in adjacent panel
4. Tab/Shift+Tab navigates between panels
5. Ctrl+T/Ctrl+R cycles through 9 GQL type categories
6. Space bar opens detail overlay for focused item
7. Supports up to 6 panels horizontally

## Features

### GraphQL Type Exploration
Supports exploring all GraphQL schema types:
- **9 Type Categories**: Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, Directive
- **Type Cycling**: Ctrl+T cycles forward, Ctrl+R cycles backward through types
- **Auto-Loading**: Panels auto-populate when switching types

### Interactive Navigation
- **Panel Focus**: Tab/Shift+Tab to navigate between panels
- **Auto-Open**: Selecting items automatically opens details in adjacent panel
- **Detail Overlay**: Space bar shows full item details in centered overlay
- **Multi-Panel**: Supports up to 6 panels horizontally

## Dependencies
- **Bubble Tea**: TUI framework for terminal applications
- **Bubbles**: Pre-built components (list, help) for Bubble Tea
- **Lipgloss**: Terminal styling and layout
- **graphql-go**: GraphQL parsing and AST manipulation
