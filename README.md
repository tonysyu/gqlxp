# gqlxp

`gqlxp` is an interactive GraphQL query explorer TUI (Terminal User Interface) for exploring GraphQL schema files. Built with Bubble Tea, it provides a multi-panel interface for navigating through all GraphQL type definitions including Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, and Directive types.

## Features

### GraphQL Type Exploration
Supports exploring all GraphQL schema types:
- **9 Type Categories**: Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, Directive
- **Type Cycling**: Ctrl+T (or `}`) cycles forward, Ctrl+R (or `{`) cycles backward through types
- **Auto-Loading**: Panels auto-populate when switching types

### Interactive Navigation
- **Panel Focus**: Tab/Shift+Tab (or `]`/`[`) to navigate between panels
- **Auto-Open**: Selecting items automatically opens details in adjacent panel
- **Detail Overlay**: Space bar shows full item details in centered overlay
- **Multi-Panel**: Supports up to 6 panels horizontally
- **Breadcrumbs**: Shows navigation path when panels scroll off-screen

## Navigation Flow
1. Application parses GraphQL schema from provided file path
2. Displays Query fields by default in main panel
3. Selecting items auto-opens details in adjacent panel
4. Tab/Shift+Tab navigates between panels
5. Ctrl+T/Ctrl+R cycles through 9 GQL type categories
6. Space bar opens detail overlay for focused item

## Usage

Currently, there's no installer or executable distribution. Build from source using:

```sh
$ git clone https://github.com/tonysyu/gqlxp
$ just build
$ ./dist/gqlxp examples/github.graphqls
```

### Library Mode

Persist schemas in a local library for quick access with favorites and documentation links:

```sh
# Add a schema to the library
$ ./dist/gqlxp library add github-api examples/github.graphqls

# List schemas in library
$ ./dist/gqlxp library list

# Interactive schema selector
$ ./dist/gqlxp --library

# Explore specific schema from library
$ ./dist/gqlxp --library github-api

# Remove schema from library
$ ./dist/gqlxp library remove github-api
```

**Library Features:**
- **Schema Selector**: Run `gqlxp --library` to open interactive schema picker
- **Favorites**: Press `f` to favorite/unfavorite types for quick identification (marked with ★)
- **Persistent Storage**: Schemas stored in `~/.config/gqlxp/schemas/` on macOS/Linux or `%APPDATA%\gqlxp\schemas\` on Windows

For local development commands:
```sh
$ just build  # Build executable to dist/gqlxp
$ just run {{PATH_TO_GRAPHQL_SCHEMA_FILE}}  # Run with schema file
$ just test  # Run all tests
$ just verify  # Run tests, lint, and fix
```

## Developer Documentation

- [Documentation Index](docs/index.md)
    - [Development Commands](docs/development.md)
    - [Architecture](docs/architecture.md)
    - [Coding Best Practices](docs/coding.md)
    - [Schema Library](docs/schema-library.md)
