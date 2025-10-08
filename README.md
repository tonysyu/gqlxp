# igq

`igq` is an interactive GraphQL query explorer TUI (Terminal User Interface) for exploring GraphQL schema files. Built with Bubble Tea, it provides a multi-panel interface for navigating through all GraphQL type definitions including Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, and Directive types.

## Usage

Currently, there's no installer or executable distribution. Build from source using:

```sh
$ git clone https://github.com/tonysyu/igq
$ just build
$ ./dist/igq examples/github.graphqls
```

For local development commands:
```sh
$ just build  # Build executable to dist/igq
$ just run {{PATH_TO_GRAPHQL_SCHEMA_FILE}}  # Run with schema file
$ just test  # Run all tests
```

## Developer Documentation

- [Documentation Index](docs/index.md)
    - [Development Commands](docs/development.md)
    - [Architecture](docs/architecture.md)
    - [Coding Best Practices](docs/coding.md)
