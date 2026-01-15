# cli-interface Specification

## Purpose
Defines the command-line interface for gqlxp, using explicit flags for schema selection instead of positional arguments.
## Requirements
### Requirement: App Subcommand
The system SHALL provide an explicit `app` subcommand that launches the TUI using schema selection via flags.

#### Scenario: Launch app with schema flag
- **WHEN** user runs `gqlxp app --schema <schema-id>` or `gqlxp app -s <schema-id>`
- **THEN** the TUI is opened with the specified schema

#### Scenario: Launch app without schema flag
- **WHEN** user runs `gqlxp app` with no schema flag
- **THEN** the library selector TUI is opened

#### Scenario: App subcommand accepts log-file flag
- **WHEN** user runs `gqlxp app --log-file debug.log`
- **THEN** debug logging is enabled to the specified file

### Requirement: Search Command Schema Selection
The system SHALL use `--schema` flag for explicit schema selection, defaulting to config default when omitted.

#### Scenario: Search with schema flag
- **WHEN** user runs `gqlxp search --schema <schema-id> <query>` or `gqlxp search -s <schema-id> <query>`
- **THEN** the search executes against the specified schema

#### Scenario: Search without schema flag
- **WHEN** user runs `gqlxp search <query>` without schema flag
- **THEN** the search executes against the default schema from config
- **AND** if no default is set, an error is returned

### Requirement: Show Command Schema Selection
The system SHALL use `--schema` flag for explicit schema selection, defaulting to config default when omitted.

#### Scenario: Show with schema flag
- **WHEN** user runs `gqlxp show --schema <schema-id> <type-name>` or `gqlxp show -s <schema-id> <type-name>`
- **THEN** the type definition is displayed from the specified schema

#### Scenario: Show without schema flag
- **WHEN** user runs `gqlxp show <type-name>` without schema flag
- **THEN** the type definition is displayed from the default schema from config
- **AND** if no default is set, an error is returned

### Requirement: Schema Flag Accepts Paths and IDs
The system SHALL accept both file paths and schema IDs for the `--schema` flag.

#### Scenario: Schema flag with library ID
- **WHEN** user runs `gqlxp search --schema github-api <query>`
- **THEN** the command uses the schema with ID "github-api" from the library

#### Scenario: Schema flag with file path
- **WHEN** user runs `gqlxp search --schema examples/github.graphqls <query>`
- **THEN** the command loads the schema from the file path (and prompts to add to library if not present)

### Requirement: Library Command
The system SHALL provide a `library` command that consolidates schema library management functionality.

#### Scenario: Library command help
- **WHEN** user runs `gqlxp library --help`
- **THEN** help text is displayed showing available subcommands

### Requirement: Show Command JSON Output
The system SHALL support JSON output format for the show command via a `--json` flag.

#### Scenario: Show with JSON flag
- **WHEN** user runs `gqlxp show --json <type-name>`
- **THEN** a pretty-printed JSON representation of the type is output to stdout
- **AND** the pager is automatically disabled
- **AND** no markdown rendering is applied

#### Scenario: Show with JSON and schema flags
- **WHEN** user runs `gqlxp show --schema <schema-id> --json <type-name>`
- **THEN** the type from the specified schema is output as JSON

#### Scenario: JSON output structure for Object types
- **WHEN** user runs `gqlxp show --json User` for an Object type
- **THEN** the JSON contains fields: name, kind, description, fields, interfaces, directives
- **AND** each field includes: name, type, description, arguments, directives

#### Scenario: JSON output structure for Enum types
- **WHEN** user runs `gqlxp show --json Status` for an Enum type
- **THEN** the JSON contains fields: name, kind, description, values, directives
- **AND** each value includes: name, description, directives

#### Scenario: JSON output structure for Query/Mutation fields
- **WHEN** user runs `gqlxp show --json Query.getUser` for a field
- **THEN** the JSON contains: name, type, description, arguments, directives

#### Scenario: JSON output structure for Directives
- **WHEN** user runs `gqlxp show --json @auth` for a directive
- **THEN** the JSON contains: name, description, locations, arguments

#### Scenario: JSON pretty-printing
- **WHEN** JSON output is generated
- **THEN** it is formatted with 2-space indentation
- **AND** it is human-readable with line breaks

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

