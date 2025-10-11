# List available commands
default:
    @echo 'Usage: just [OPTIONS] [ARGUMENTS]...'
    @just -l

# Run the gqlxp tui
[group('app')]
run schemaPath:
    go run ./cmd/gqlxp {{schemaPath}}

# Run tests (defaults to all tests in projects)
[group('code')]
test target=tests:
    go test {{target}}
tests := "./..."

# Run tests and generate coverage report
[group('code')]
test-coverage:
    go test -coverprofile=./build/coverage.out ./...
    go tool cover -html=./build/coverage.out

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
    go build -o dist/gqlxp ./cmd/gqlxp
