# List available commands
default:
    @echo 'Usage: just [OPTIONS] [ARGUMENTS]...'
    @just -l

empty := ""

# Run gqlxp
[group('app')]
run +ARGS=empty:
    go run ./cmd/gqlxp {{ARGS}}

# Run the gqlxp tui with logging to debug.log file
[group('app')]
run-with-log +ARGS=empty:
    GQLXP_LOGFILE=debug.log go run ./cmd/gqlxp {{ARGS}}

# Document public signatures for package
[group('code')]
doc pkg:
    go doc {{pkg}} | bat -l go

# Show all documenttion for package
[group('code')]
doc-all pkg:
    go doc -all {{pkg}} | bat -l go

# Run tests (defaults to all tests in projects)
[group('code')]
test target=tests:
    go test {{target}}
tests := "./..."

# Run tests and generate coverage report
[group('code')]
test-coverage:
    go test -coverpkg=./... -coverprofile=./build/coverage.out ./...
    go tool cover -html=./build/coverage.out

# Run code formatters and linters
[group('code')]
lint-fix:
    go fmt ./...
    go mod tidy
    go vet ./...
    go tool golangci-lint run

# Run tests tests, lint, and fix
[group('code')]
verify: && test lint-fix
    @echo "Testing, linting, and fixing"

# Display directory tree
[group('code')]
tree:
    tree -I openspec -d

# Build executable
[group('deploy')]
build:
    go build -o dist/gqlxp ./cmd/gqlxp

# Install executable
[group('deploy')]
install:
    go install ./cmd/gqlxp

# Record demo
[group('deploy')]
create-demo: install
    vhs demo.tape
