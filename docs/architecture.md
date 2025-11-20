# Architecture

For build and development commands, see [Development Commands](development.md).

## Package Structure
- **`cmd/gqlxp`**: Main application entry point
  - `main()` - Reads schema file or loads from library and starts TUI
  - Library subcommands (add, list, remove)
- **`library`**: Schema library with persistent storage
  - `Library` interface - Schema management operations
  - `FileLibrary` - File-based storage implementation
  - Metadata storage in `~/.config/gqlxp/schemas/`
- **`gql`**: GraphQL schema parsing and type extraction
  - `ParseSchema()` - Parses GraphQL schema and extracts all type definitions
  - `GraphQLSchema` struct - Contains maps for Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, and Directive types
  - `NamedToTypeDef()` - Resolves type names to their definitions
- **`tui`**: Terminal user interface built on Bubble Tea
  - **`tui/libselect`**: Library selection mode for choosing schemas
  - **`tui/xplr`**: Schema exploration mode
    - `Model` - Main explorer coordinating UI components and navigation
    - **`tui/xplr/components`**: UI components for explorer
      - `Panel` - Panel displaying lists with auto-open behavior
      - `ListItem` interface - Interactive list items with `Open()` method
      - `SimpleItem` - Basic ListItem implementation
    - **`tui/xplr/navigation`**: Navigation state management
      - `NavigationManager` - Coordinates panel stack, breadcrumbs, and type selection
  - **`tui/overlay`**: Detail view overlay for displaying full item information
  - **`tui/adapters`**: Converts GraphQL AST types to UI components (shared)
  - **`tui/config`**: Configuration and styling (shared)
- **`tests/acceptance`**: Acceptance tests for end-to-end workflows
  - Test harness with domain-specific helpers for TUI interaction
  - Workflow tests verifying navigation, type switching, and overlay behavior
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
