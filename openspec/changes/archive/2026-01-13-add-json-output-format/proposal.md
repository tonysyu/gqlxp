# Change: Add JSON Output Format to Show Command

## Why
CLI users need machine-readable output from `gqlxp show` for automation, scripting, and integration with other tools. The current markdown-only output requires parsing, making programmatic usage difficult.

## What Changes
- Add `--json` flag to `gqlxp show` command
- Generate structured JSON output with custom simplified format matching markdown content
- Pretty-print JSON by default with 2-space indentation
- Auto-disable pager when JSON output is requested
- Create new `gqlfmt/json.go` module for JSON generation

## Impact
- Affected specs: cli-interface
- Affected code:
  - `cli/show.go` - Add --json flag handling
  - `gqlfmt/` - New json.go module for JSON generation
  - Tests for JSON output functionality
