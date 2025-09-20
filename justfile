# List available commands
default:
    @echo 'Usage: just [OPTIONS] [ARGUMENTS]...'
    @just -l

# Run the gq tui
[group('app')]
run:
    go run ./cmd/gq

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

# Build executable
[group('deploy')]
build:
    go build -o dist/gq ./cmd/gq
