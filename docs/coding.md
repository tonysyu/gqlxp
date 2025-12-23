# Coding Best Practices

## Minimize Public API

**Default to private.** Use lowercase names for types, functions, and variables. Only capitalize when external packages require access.

- **Field access patterns:**
  - Fields with only getters: Keep private, provide getter methods
  - Fields with setters: Make public (avoid getter/setter boilerplate)
  - Fields requiring validation: Keep private, provide methods that validate
- Periodically review public APIs for privatization opportunities

## Testing

### Unit Tests

- Prefer `package foo_test` (black-box) for public API tests
- Use `package foo` (white-box) only when testing private implementation
- Use `"github.com/matryer/is"` instead of `t.Error`/`t.Errorf`

### Acceptance Tests

**Use acceptance tests for end-to-end user workflows.** Located in `tests/acceptance/`, these tests verify complete interactions through the TUI using the test harness.

When to use acceptance tests:
- Multi-step navigation workflows (panel navigation, breadcrumb verification)
- Type switching and state transitions across the application
- Overlay interactions and screen content verification
- Complete user journeys that span multiple components

When to use unit tests:
- Single component behavior
- Edge cases and error handling
- Internal implementation details
- Pure business logic

See `tests/acceptance/workflows_test.go` for examples.

## Error Handling

- Return errors for recoverable failures (avoid panics)
- Use `fmt.Errorf` with context
- Validate inputs early

## Line of Sight

**Keep the happy path left-aligned.** Handle errors and edge cases early with guard clauses.

- Main logic should have minimal nesting
- Return early for error cases and validation failures
- Avoid deeply nested if statements
