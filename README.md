# gq

`gq` is a GraphQL query explorer TUI (Terminal User Interface) application that parses GraphQL schema files and provides an interactive terminal interface for exploring GraphQL queries. The application uses the Bubble Tea framework to create a multi-panel interface where users can navigate through GraphQL schema definitions.

## Schema Requirements

The application expects a GraphQL schema file at `examples/github.graphqls`. The parser specifically looks for Query type definitions and extracts field information including arguments, return types, and descriptions.

## Documentation

- [Documentation Index](docs/index.md)
    - [Development Commands](docs/development.md)
    - [Architecture](docs/architecture.md)
