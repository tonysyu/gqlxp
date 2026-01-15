# cli-interface Specification Delta

## ADDED Requirements

### Requirement: Search Command JSON Output
The system SHALL support JSON output format for the search command via a `--json` flag.

#### Scenario: Search with JSON flag
- **WHEN** user runs `gqlxp search --json <query>`
- **THEN** search results are output as a JSON array to stdout
- **AND** the pager is automatically disabled
- **AND** no styled text formatting is applied

#### Scenario: Search with JSON and schema flags
- **WHEN** user runs `gqlxp search --schema <schema-id> --json <query>`
- **THEN** search results from the specified schema are output as JSON

#### Scenario: JSON output structure for search results
- **WHEN** search results are output as JSON
- **THEN** the output is a JSON array of result objects
- **AND** each result contains: path, kind, name, description, score

#### Scenario: JSON pretty-printing
- **WHEN** JSON output is generated
- **THEN** it is formatted with 2-space indentation
- **AND** it is human-readable with line breaks

#### Scenario: Empty results as JSON
- **WHEN** a search with `--json` returns no matches
- **THEN** an empty JSON array `[]` is output
