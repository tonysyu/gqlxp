# Development Commands

For project architecture and feature details, see [Architecture](architecture.md).

## Available Commands
Use `just` to see all available commands defined in the justfile.

## Building and Running
- `just build` - Build the application
- `just run <schema-file>` - Run the application with a GraphQL schema file
- `just run-with-log <schema-file>` - Run with debug logging enabled (logs to debug.log)

## Testing
- `just test` - Run all tests
- `just "test -v"` - Run all tests with verbose output
- `just test ./gql` - Run only GraphQL parsing tests
- `just test ./tests/acceptance` - Run only acceptance tests
- `just test "./gql -v"` - Run GraphQL tests with verbose output
- `just test-coverage` - Generate HTML coverage report

## Code Quality
- `just lint-fix` - Format code, tidy modules, and run static analysis
- `just verify` - Run tests, lint, and fix all code quality issues

## Documentation
Use the `go doc` command to get documentation for this module. You can pipe output to
[bat](https://github.com/sharkdp/bat) for a syntax-highlighted version:
```
go doc gql | bat -l go
go doc tui | bat -l go
go doc tui/adapters | bat -l go
go doc tui/components | bat -l go
go doc tui/config | bat -l go
```
