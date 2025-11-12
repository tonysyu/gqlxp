# Architecture

For build and development commands, see [Development Commands](development.md).

## Package Structure
- **`cmd/gqlxp`**: Main application entry point
  - `main()` - Reads schema file and starts TUI
- **`gql`**: GraphQL schema parsing and type extraction
  - `ParseSchema()` - Parses GraphQL schema and extracts all type definitions
  - `GraphQLSchema` struct - Contains maps for Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, and Directive types
  - `NamedToTypeDef()` - Resolves type names to their definitions
- **`tui`**: Terminal user interface built on Bubble Tea
  - `mainModel` - Root model managing panels, navigation, and GQL type toggling
  - `overlayModel` - Detail view overlay for displaying full item information
  - `breadcrumbsModel` - Navigation path display when panels scroll off-screen
  - **`tui/components`**: Reusable UI components
    - `Panel` interface - Generic panel abstraction
    - `ListItem` interface - Interactive list items with `Open()` method
    - `ListPanel` - Panel displaying lists with auto-open behavior
    - `SimpleItem` - Basic ListItem implementation
  - **`tui/adapters`**: Converts GraphQL AST types to UI components
  - **`tui/config`**: Configuration and styling
- **`utils/testx`**: Testing utilities
  - **`utils/testx/assert`**: Test assertion helpers
- **`utils/text`**: Text manipulation utilities

## Key Interfaces
- **`ListItem`**: Extends `list.DefaultItem` with `Open() Panel` method for interactive items
- **`Panel`**: Implements `tea.Model` with `SetSize(width, height int)` for resizable content

## Dependencies
- **Bubble Tea**: TUI framework for terminal applications
- **Bubbles**: Pre-built components (list, help) for Bubble Tea
- **Lipgloss**: Terminal styling and layout
- **Glamour**: Markdown rendering for detail overlays
- **vektah/gqlparser/v2**: GraphQL parsing and AST manipulation
