# Tasks: Add JSON Output to Search Command

## Implementation Tasks

1. [ ] Add `--json` flag to search command in `cli/search.go`
2. [ ] Create JSON output function for search results
3. [ ] Integrate JSON output into search command action (bypass pager and styled output)
4. [ ] Add tests for JSON output functionality
5. [ ] Run `just test` to verify changes

## Validation

- `gqlxp search --json user` outputs valid JSON array
- `gqlxp search --json --limit 5 user` respects limit in JSON output
- JSON output is pretty-printed with 2-space indentation
- Pager is automatically disabled for JSON output
