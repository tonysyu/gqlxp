## Context

gqlxp has TUI and non-interactive CLI commands (`search`, `show`) for exploring GraphQL schemas. It lacks a way to validate operations, requiring developers to use external tools or write custom scripts. The `gqlparser/v2` library already used for schema parsing also provides operation validation, making this a low-cost addition.

## Goals / Non-Goals

**Goals:**
- Add a `parse` subcommand that validates a GraphQL operation against a schema
- Report errors with line/column location
- Accept input from a file argument or stdin
- Exit non-zero on invalid operations (enabling scripting/CI use)

**Non-Goals:**
- Formatting or pretty-printing the operation
- Partial validation (always validate fully against the schema)
- Watching files for changes

## Decisions

**Input: file argument or stdin**
The command accepts an optional file path argument. If omitted, reads from stdin. This matches the Unix convention and enables both `gqlxp parse query.graphql` and `cat query.graphql | gqlxp parse`.

Alternative considered: require explicit `-` for stdin. Rejected — inferring stdin when no arg is given is more ergonomic and consistent with tools like `jq`.

**Validation via `gqlparser/v2/validator`**
The `gqlparser/v2` library is already a direct dependency used for schema parsing. Its `validator.Validate(schema *ast.Schema, document *ast.QueryDocument)` function returns a `gqlerrors.List`, providing both syntax and semantic validation (unknown fields, type mismatches, missing required arguments, etc.).

Alternative considered: implement custom validation. Rejected — `gqlparser/v2` already covers the required cases with battle-tested logic.

**Error output format: `<source>:<line>:<col>: <message>`**
Errors are printed to stdout one per line in this format, matching common compiler/linter conventions. `<source>` is the file path or `<stdin>`. Exits with code 1 on any error.

Alternative considered: always use `--json` structured output. Rejected — plain text is more useful as the default; `--json` can be added later if needed.

**Schema resolution: reuse `resolveSchemaFromArgument()`**
The existing `--schema/-s` flag and `resolveSchemaFromArgument()` function (used by `search` and `show`) handle file paths, library IDs, and the default schema. No new logic needed.

## Risks / Trade-offs

- `gqlparser/v2` may not catch every possible GraphQL spec violation, but it covers the common cases that matter in practice.
- Reading from stdin blocks if the user forgets to pipe input. This is standard Unix behavior and not a special concern.

## Migration Plan

No migration needed — this is a new additive command with no effect on existing behavior.

## Open Questions

- Should a `--json` flag for structured error output be included in the initial implementation? (Can be deferred to a follow-up.)
