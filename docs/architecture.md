# Architecture

For build and development commands, see [Development Commands](development.md).

## Package Structure
- **`cmd/gqlxp`**: Application entry point
- **`cli`**: CLI setup, user prompts, and output formatting
- **`library`**: Schema storage in `~/.config/gqlxp/schemas/`
- **`gql`**: GraphQL schema parsing and type resolution
- **`tui`**: Bubble Tea-based terminal interface
  - **`tui/libselect`**: Library selection mode
  - **`tui/xplr`**: Schema exploration mode
  - **`tui/overlay`**: Detail view overlay
  - **`tui/adapters`**: GraphQL AST to UI conversion
- **`tests/acceptance`**: End-to-end workflow tests
- **`tests/fitness`**: Architectural constraint tests (package hierarchy and dependency restrictions)
- **`utils`**: Shared utilities (testing, text manipulation)

## Key Interfaces
- **`ListItem`**: Interactive list items with `Open() Panel` method
- **`Panel`**: Resizable Bubble Tea model with `SetSize(width, height int)`

## Dependencies
- **Bubble Tea/Bubbles/Lipgloss**: TUI framework and components
- **Glamour**: Markdown rendering
- **vektah/gqlparser/v2**: GraphQL parsing
- **urfave/cli/v3**: CLI framework

Dependency usage is restricted by package. See `tests/fitness/dependency_test.go`.
