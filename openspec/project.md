# Project Context

## Purpose
`gqlxp` is an interactive GraphQL query explorer TUI (Terminal User Interface) for exploring GraphQL schema files. Provides multi-panel interface for navigating through all GraphQL type definitions including Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, and Directive types.

## Tech Stack
- **Language**: Go 1.25.1
- **TUI Framework**: Bubble Tea (terminal application framework)
- **UI Components**: Bubbles (pre-built components), Lipgloss (styling/layout), Glamour (markdown rendering)
- **GraphQL Parser**: vektah/gqlparser/v2
- **Testing**: matryer/is (assertion library)
- **Build Tool**: just (command runner)

## Project Conventions

### Code Style
- **Default to private**: Use lowercase names; only capitalize when external access required
- **Minimal public API**: Use accessor methods instead of exposing struct fields
- **Error handling**: Return errors for recoverable failures; use `fmt.Errorf` with context; validate inputs early
- **Naming**: Follow Go conventions; descriptive names over abbreviations

### Architecture Patterns
- **Package structure**: `cmd/gqlxp` (entry), `gql` (parsing), `tui` (UI), `utils` (helpers)
- **UI components**: Reusable components implementing `tea.Model` interface
- **Navigation**: Centralized `NavigationManager` coordinates panel stack and type selection
- **Adapters**: Convert GraphQL AST types to UI components
- **Interfaces**: `ListItem` (interactive items), `Panel` (resizable content)

### Testing Strategy
- **Black-box tests**: Prefer `package foo_test` for public API tests
- **White-box tests**: Use `package foo` only for private implementation testing
- **Assertions**: Use `github.com/matryer/is` instead of `t.Error`/`t.Errorf`
- **Coverage**: Run `just test-coverage` for HTML coverage reports
- **Verification**: Run `just verify` (tests + lint + fix) before commits

### Git Workflow
- **Main branch**: `main`
- **Build validation**: Always run `just test` after code changes
- **Code quality**: Run `just verify` before commits (tests, lint, format)

## Domain Context
- **GraphQL Types**: Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, Directive (9 categories)
- **Navigation flow**: Parse schema → Display types → Auto-open details → Panel navigation → Type cycling
- **Panel behavior**: Auto-populate on type switch; auto-open details in adjacent panel; breadcrumbs for off-screen panels

## Important Constraints
- **No external services**: Operates entirely on local GraphQL schema files
- **Terminal-only**: No web or GUI interface
- **Read-only**: Does not modify schema files; exploration only

## External Dependencies
- **Charm ecosystem**: Bubble Tea, Bubbles, Lipgloss, Glamour (actively maintained)
- **GraphQL parser**: vektah/gqlparser/v2 (stable, widely used)
- **No runtime dependencies**: Compiles to standalone binary
