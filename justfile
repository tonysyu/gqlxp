# List available commands
default:
    @echo 'Usage: just [OPTIONS] [ARGUMENTS]...'
    @just -l

# Run the igq tui
[group('app')]
run schemaPath:
    go run ./cmd/igq {{schemaPath}}

# Run tests (defaults to all tests in projects)
[group('code')]
test target=tests:
    go test {{target}}
tests := "./..."

# Run code formatters and linters
[group('code')]
lint-fix:
    go fmt ./...
    go mod tidy
    go vet ./...

# Run tests tests, lint, and fix
[group('code')]
verify: && test lint-fix
    @echo "Testing, linting, and fixing"

# Build executable
[group('deploy')]
build:
    go build -o dist/igq ./cmd/igq
