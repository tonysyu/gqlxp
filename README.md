# gqlxp

`gqlxp` is an interactive GraphQL query explorer TUI (Terminal User Interface) for exploring GraphQL schema files. Built with Bubble Tea, it provides a multi-panel interface for navigating through all GraphQL type definitions including Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, and Directive types.

![Demo of gqlxp use](demo.gif "gqlxp-demo")

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
5. Ctrl+T/Ctrl+R cycles through GQL type categories
6. Space bar opens detail overlay for focused item

## Usage

Currently, there's no installer or executable distribution. Install from source using:

```sh
$ git clone https://github.com/tonysyu/gqlxp
$ just install
```

Then open a schema file to explore:
```sh
$ gqlxp examples/github.graphqls
```

### Schema library

Schemas are automatically saved to your library on first use:
```sh
# Load schema file (prompts for library details on first use)
$ gqlxp examples/github.graphqls
Enter schema ID (lowercase letters, numbers, hyphens) [github]: github-api
Enter display name [github-api]: GitHub GraphQL API

# Open library selector (when no arguments provided)
$ gqlxp

# Subsequent loads detect if file has changed
$ gqlxp examples/github.graphqls
Schema file has changed since last import.
Update library (y/n): y
```

### Local development
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
