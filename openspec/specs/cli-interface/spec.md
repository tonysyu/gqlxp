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

