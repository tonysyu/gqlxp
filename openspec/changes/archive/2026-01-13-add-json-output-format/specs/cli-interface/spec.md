## ADDED Requirements

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
