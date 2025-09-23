# igq

`igq` is an interactive GraphQL query explorer TUI (Terminal User Interface) application that parses GraphQL schema files and provides an interactive terminal interface for exploring GraphQL queries and mutations. The application uses the Bubble Tea framework to create a multi-panel interface where users can navigate through GraphQL schema definitions and toggle between Query and Mutation fields.

## Schema Requirements

The application expects a GraphQL schema file at `examples/github.graphqls`. The parser extracts both Query and Mutation type definitions, including field information such as arguments, return types, and descriptions.

## Documentation

- [Documentation Index](docs/index.md)
    - [Development Commands](docs/development.md)
    - [Architecture](docs/architecture.md)
