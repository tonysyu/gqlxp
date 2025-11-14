# Coding Best Practices

## Minimize Public API

**Default to private.** Use lowercase names for types, functions, and variables. Only capitalize when external packages require access.

- Use accessor methods instead of exposing struct fields
- Periodically review public APIs for privatization opportunities

## Testing

- Prefer `package foo_test` (black-box) for public API tests
- Use `package foo` (white-box) only when testing private implementation
- Use `"github.com/matryer/is"` instead of `t.Error`/`t.Errorf`

## Error Handling

- Return errors for recoverable failures (avoid panics)
- Use `fmt.Errorf` with context
- Validate inputs early
