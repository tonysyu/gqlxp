## MODIFIED Requirements

### Requirement: App Subcommand
The system SHALL provide an explicit `app` subcommand that launches the TUI using explicit schema selection flags.

#### Scenario: Launch app with default schema
- **WHEN** user runs `gqlxp app --default` or `gqlxp app -d`
- **THEN** the TUI is opened with the default schema

#### Scenario: Launch app with schema flag
- **WHEN** user runs `gqlxp app --schema <schema-id>` or `gqlxp app -s <schema-id>`
- **THEN** the TUI is opened with the specified schema

#### Scenario: Launch app without schema argument
- **WHEN** user runs `gqlxp app` with no schema flags
- **THEN** the library selector TUI is opened

#### Scenario: App subcommand accepts log-file flag
- **WHEN** user runs `gqlxp app --log-file debug.log`
- **THEN** debug logging is enabled to the specified file

### Requirement: Search Command Schema Selection
The system SHALL provide explicit schema selection flags for the search command.

#### Scenario: Search with schema flag
- **WHEN** user runs `gqlxp search --schema <schema-id> <query>` or `gqlxp search -s <schema-id> <query>`
- **THEN** the search executes against the specified schema

#### Scenario: Search with default schema flag
- **WHEN** user runs `gqlxp search --default <query>` or `gqlxp search -d <query>`
- **THEN** the search executes against the default schema

#### Scenario: Search without schema flag
- **WHEN** user runs `gqlxp search <query>` without schema flags
- **THEN** an error is returned indicating schema selection is required

### Requirement: Show Command Schema Selection
The system SHALL provide explicit schema selection flags for the show command.

#### Scenario: Show with schema flag
- **WHEN** user runs `gqlxp show --schema <schema-id> <type-name>` or `gqlxp show -s <schema-id> <type-name>`
- **THEN** the type definition is displayed from the specified schema

#### Scenario: Show with default schema flag
- **WHEN** user runs `gqlxp show --default <type-name>` or `gqlxp show -d <type-name>`
- **THEN** the type definition is displayed from the default schema

#### Scenario: Show without schema flag
- **WHEN** user runs `gqlxp show <type-name>` without schema flags
- **THEN** an error is returned indicating schema selection is required

## ADDED Requirements

### Requirement: Explicit Default Schema Flag
The system SHALL provide a `--default` flag (short form `-d`) to explicitly use the default schema.

#### Scenario: Default flag without default schema set
- **WHEN** user runs a command with `--default` flag and no default schema is configured
- **THEN** an error message indicates no default schema is set and suggests using `gqlxp library default`

#### Scenario: Schema flag takes precedence
- **WHEN** user runs a command with both `--schema` and `--default` flags
- **THEN** the `--schema` flag value is used and `--default` is ignored

### Requirement: Schema Flag Accepts Paths and IDs
The system SHALL accept both file paths and schema IDs for the `--schema` flag.

#### Scenario: Schema flag with library ID
- **WHEN** user runs `gqlxp search --schema github-api <query>`
- **THEN** the command uses the schema with ID "github-api" from the library

#### Scenario: Schema flag with file path
- **WHEN** user runs `gqlxp search --schema examples/github.graphqls <query>`
- **THEN** the command loads the schema from the file path (and prompts to add to library if not present)

## REMOVED Requirements

### Requirement: Root Command Default Behavior
**Reason**: Removing positional argument support in favor of explicit flags
**Migration**: Use `gqlxp app --default` instead of `gqlxp <schema-file>`, or use `gqlxp app` with no arguments for library selector
