# igq

`igq` is an interactive GraphQL query explorer TUI (Terminal User Interface) application that parses GraphQL schema files and provides an interactive terminal interface for exploring GraphQL queries and mutations. The application uses the Bubble Tea framework to create a multi-panel interface where users can navigate through GraphQL schema definitions and toggle between Query and Mutation fields.

## Usage

Currently, there's no installer or executable distribution. You'll have to build from
source using:

```sh
$ git clone https://github.com/tonysyu/igq
$ go build -o dist/igq ./cmd/igq
$ ./dist/igq examples/github.graphqls
```

For local development, use [justfile](https://just.systems/man/en/introduction.html) to
run dev commands:
```sh
$ just build  # creates `dist/igq`, like command above
$ just run {{PATH_TO_GRAPHQL_SCHEMA_FILE}}
```

## Developer Documentation

- [Documentation Index](docs/index.md)
    - [Development Commands](docs/development.md)
    - [Architecture](docs/architecture.md)
    - [Coding Best Practices](docs/coding.md)
