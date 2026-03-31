# gqlxp

`gqlxp` is an interactive GraphQL query explorer TUI (Terminal User Interface) for exploring GraphQL schema files. It provides a multi-panel interface for navigating through GraphQL schemas.

![Demo of gqlxp use](demo.gif "gqlxp-demo")

## TUI Components

```
                 ┌─────────────────────────────────────────────────────────────────────┐
 Type Class Nav  │ ▸ Query | Mutation | Object | Input | Enum | Scalar | Interface | … │
                 │                                                                     │
    Breadcrumbs  │ Query > user > User                                                 │
                 ├────────────────────────────────────┬────────────────────────────────┤
Type/Field Name  │ User                               │ email                          │
                 │                                    │                                │
  Panel SubTabs  │ ▸ Fields | Interfaces              │ Type                           │
                 ├────────────────────────────────────┤                                │
    Panel Items  │   name: String                     │ EmailAddress                   │
                 │ ▸ email: EmailAddress              │                                │
                 │   avatarUrl: URI                   │                                │
                 │   createdAt: Date                  │                                │
                 │                                    │                                │
                 │                                    │                                │
       Page Nav  │ ...                                │                                │
                 └────────────────────────────────────┴────────────────────────────────┘

                  Active Panel                         Detail Panel
```

- **Type Class Nav**: Navbar for GraphQL Type Classes + Search-mode
  - GraphQL Type Classes: Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, Directive
  - Search: Search all GraphQL types and fields by name, description. See [[#Search Syntax]] below
  - Currently focused Type Class will be displayed in the Active Panel
  - Keymaps:
    - `}`: Focus on next type
    - `{`: Focus on previous type
- **Active Panel**: Left panel displaying currently focused Type or Field
  - Keymaps:
    - `]`/`Tab`: Make next panel the Active Panel
    - `[`/`⇧+Tab` Make previous panel the Active Panel
    - `j`/`↓`: Next item in Panel Items list
    - `k`/`↑`: Previous item in Panel Items list
    - `l`/`→`: Next page in Panel Items list
    - `h`/`←`: Previous page in Panel Items list
    - `L`/`⇧+→`: Next Panel SubTab
    - `H`/`⇧+←`: Previous Panel SubTab
    - `Spacebar`: Open detail view of the current Type or Field in Active Panel
- **Detail Panel**: Right panel displaying currently focused Panel Item in Active Panel
- **Breadcrumbs**: Displays names of current Active Panel and any hidden Panels used to
  navigate to the current Active Panel
- **Detail Overlay** (not pictured): Shows full details of Type or Field in Active Panel
  - Keymaps:
    - `Spacebar`: Toggles display of Detail Overlay

## Quick Start

```sh
# Install gqlxp executable (defaults to $HOME/bin; see below for custom executable path)
curl -sSfL https://raw.githubusercontent.com/tonysyu/gqlxp/main/install.sh

# Add GraphQL schema to library using graphql endpoint with introspection enabled
gqlxp library add --id rick-and-morty-api https://rickandmortyapi.com/graphql

# Run gqlxp interactive schema explorer
gqlxp
```

## Installation

### Quick install (Linux/macOS)

```sh
curl -sSfL https://raw.githubusercontent.com/tonysyu/gqlxp/main/install.sh | sh
```

To install to a custom directory:
```sh
curl -sSfL https://raw.githubusercontent.com/tonysyu/gqlxp/main/install.sh | sh -s -- -b ~/.local/bin
```

### Install from source

```sh
git clone https://github.com/tonysyu/gqlxp
cd gqlxp
just install
```

## Usage

Then open a schema file to explore:
```sh
# Open library selector (app is optional here)
$ gqlxp app

# Open a specific schema file (-s is an alias for --schema)
$ gqlxp app -s examples/github.graphqls

# Open a schema from library
$ gqlxp app -s github-api

# Jump directly to a specific type:
$ gqlxp app -s github-api --select User

# Jump directly to a specific field within a type:
$ gqlxp app -s github-api --select Query.user
```

### Schema library

Schemas are automatically saved to your library on first use:
```sh
# Load schema file (prompts for library details on first use)
$ gqlxp app -s examples/github.graphqls
Enter schema ID (lowercase letters, numbers, hyphens) [github]: github-api
Enter display name [github-api]: GitHub GraphQL API

# PREFERRED: Load schema from graphql endpoint (requires introspection to download
schema)
$ gqlxp library add --id rick-and-morty-api https://rickandmortyapi.com/graphql

# Schemas added/updated with url can be updated
$ gqlxp library update --id rick-and-morty-api
Fetching schema from https://rickandmortyapi.com/graphql...
Schema 'rick-and-morty-api' is already up to date (timestamp updated)

# Set default schema for commands that omit --schema
$ gqlxp library default github-api

# Open app (shows library selector)
$ gqlxp app

# Load a schema from library by ID
$ gqlxp app -s github-api
```

👉 For more information about the schema library, see [docs/schema-library.md](docs/schema-library.md).

### Search

You can search for types and fields across your schema within the TUI app, or from the
CLI:
```sh
# Search for "user" in schema defined by file path (-s is an alias for --schema)
$ gqlxp search -s examples/github.graphqls user

# Search for "user" using schema in library named "github-api"
$ gqlxp search -s github-api user

# Search using default schema (omit --schema flag)
$ gqlxp search mutation

# Limit number of results
$ gqlxp search --limit 5 repository
```

The search command indexes your schema for fast full-text search across type names, field names, and descriptions.

👉 For search syntax, fields, and document types, see [docs/search.md](docs/search.md).

#### Indexing

Schemas are automatically indexed when added or updated. Indexes are stored in `~/.config/gqlxp/schemas/<schema-id>.bleve/` and removed when schemas are deleted.

To rebuild an index manually:
```sh
gqlxp library reindex <schema-id>
gqlxp library reindex --all  # Rebuild all indexes
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
    - [Search Syntax and Configuration](docs/search.md)
