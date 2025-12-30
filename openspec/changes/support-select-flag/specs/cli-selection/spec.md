# cli-selection Specification Delta

## ADDED Requirements

### Requirement: CLI Select Flag
The system SHALL provide a `--select` flag on the `app` command to open the TUI with a specific type or field pre-selected.

#### Scenario: Select type by name
- **GIVEN** a schema file with a type named "User"
- **WHEN** user runs `gqlxp app schema.graphqls --select User`
- **THEN** the TUI opens with the "User" type selected in the appropriate type category (Object, Input, etc.)
- **AND** the type's detail panel is displayed

#### Scenario: Select field within type
- **GIVEN** a schema file with type "Query" containing field "user"
- **WHEN** user runs `gqlxp app schema.graphqls --select Query.user`
- **THEN** the TUI opens with the Query type category active
- **AND** the "user" field is selected in the Query panel
- **AND** the field's detail panel is displayed

#### Scenario: Select with library selector
- **WHEN** user runs `gqlxp app --select User` without schema file argument
- **THEN** the library selector opens normally (selection is ignored without schema)

#### Scenario: Invalid type name
- **GIVEN** a schema file that does not contain type "NonExistent"
- **WHEN** user runs `gqlxp app schema.graphqls --select NonExistent`
- **THEN** the TUI opens in default state (Query category, no selection)
- **AND** no error is displayed to the user

#### Scenario: Invalid field name
- **GIVEN** a schema with type "Query" that does not contain field "missing"
- **WHEN** user runs `gqlxp app schema.graphqls --select Query.missing`
- **THEN** the TUI opens with Query type category active
- **AND** the Query type is selected but the field is not found
- **AND** no error is displayed to the user

#### Scenario: Flag available on app subcommand
- **WHEN** user runs `gqlxp app --help`
- **THEN** the help text shows the `--select` flag with usage description

### Requirement: Type Category Detection
The system SHALL automatically determine the correct GQL type category (Query, Mutation, Object, etc.) for a selected type.

#### Scenario: Detect object type
- **GIVEN** a schema with an Object type named "User"
- **WHEN** user selects "User" via `--select User`
- **THEN** the TUI switches to the Object type category
- **AND** "User" is selected in the Object panel

#### Scenario: Detect input type
- **GIVEN** a schema with an Input type named "UserInput"
- **WHEN** user selects "UserInput" via `--select UserInput`
- **THEN** the TUI switches to the Input type category
- **AND** "UserInput" is selected in the Input panel

#### Scenario: Detect enum type
- **GIVEN** a schema with an Enum type named "Status"
- **WHEN** user selects "Status" via `--select Status`
- **THEN** the TUI switches to the Enum type category
- **AND** "Status" is selected in the Enum panel

### Requirement: Initial Selection State
The system SHALL apply the selection after schema is loaded and before the first render.

#### Scenario: Selection applied on startup
- **WHEN** TUI initializes with a `--select` target
- **THEN** the selection is applied before the first screen render
- **AND** the user sees the selected item immediately

#### Scenario: Panel navigation with selection
- **GIVEN** user starts with `--select Query.user`
- **WHEN** the TUI initializes
- **THEN** the navigation manager positions panels to show Query and user field
- **AND** breadcrumbs reflect the navigation path if applicable
