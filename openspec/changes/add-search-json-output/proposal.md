# Proposal: Add JSON Output to Search Command

## Summary
Add a `--json` flag to `gqlxp search` command to output search results as JSON, similar to the existing `--json` flag on `gqlxp show`.

## Motivation
- Enables programmatic consumption of search results (piping to `jq`, integration with scripts)
- Consistent with existing `show --json` behavior
- Useful for automation and tooling integrations

## Scope
- Add `--json` flag to search command
- Output search results as pretty-printed JSON array
- Disable pager automatically when JSON output is enabled
- Update cli-interface spec with new requirement

## Out of Scope
- Changes to search result structure or scoring
- Additional output formats (YAML, CSV, etc.)
