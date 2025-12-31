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
# Or explicitly use the app command:
$ gqlxp app examples/github.graphqls

# Jump directly to a specific type:
$ gqlxp app examples/github.graphqls --select User

# Jump directly to a specific field within a type:
$ gqlxp app examples/github.graphqls --select Query.user
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
# Or explicitly:
$ gqlxp app

# Subsequent loads detect if file has changed
$ gqlxp examples/github.graphqls
Schema file has changed since last import.
Update library (y/n): y
```

### Search

Search for types and fields across your schema:
```sh
# Search in a specific schema file
$ gqlxp search examples/github.graphqls user

# Search using default schema
$ gqlxp search mutation

# Limit number of results
$ gqlxp search --limit 5 repository
```

The search command indexes your schema for fast full-text search across type names, field names, and descriptions. Indexes are automatically created when schemas are added to your library and rebuilt when schemas are updated.

#### Search Syntax

Search is implemeted using [bleve](https://github.com/blevesearch/bleve) and supports `bleve`'s [query syntax](https://blevesearch.com/docs/Query-String-Query/).

```sh
# Search for Query related to "user"
$ gqlxp search examples/github.graphqls "type:Query user"

# Search for Query with name "user" (results match _either_ type or name)
$ gqlxp search examples/github.graphqls "type:Query name:user"

# Search for Query with name "user" (results match _both_ type and name)
$ gqlxp search examples/github.graphqls "+type:Query +name:user"

# Search for Mutation with name _containing_ "user"
$ gqlxp search examples/github.graphqls "+type:Mutation +name:*user*"

# Search for Object or Interface with name starting with "repo" (slashes denote regex)
$ gqlxp search examples/github.graphqls "+type:/(Object|Interface)/ +name:repo*"

# Search for "user" in all GraphQL types and exclude all fields
$ gqlxp search examples/github.graphqls "user -type:*Field"
```

#### Search fields
Each search document contains the following fields:
- `type`: GraphQL type or field (see section below)
- `name`: Name of the type or field
- `description`: Description/docstring of type or field
- `path`: Qualified name, which matches `name` for types and `<parent-name>.<name>` for fields
    - Queries and mutations will have fixed paths `Query.<name>` and `Mutation.<name>`,
      respectively

#### Document `type`s
You can specify the standard GraphQL types for the `type` field in your query:
- `Query`
- `Mutation`
- `Object`
- `Input`
- `Enum`
- `Scalar`
- `Interface`
- `Union`
- `Directive`

In addition to the standard types, the follwing field types are defined:
- `ObjectField`
- `InputField`
- `InterfaceField`

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
    - [Schema Search](docs/search.md)
