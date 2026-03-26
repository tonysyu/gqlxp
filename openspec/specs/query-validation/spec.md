# query-validation Specification

## Purpose
Defines the `parse` subcommand for validating GraphQL operation documents against a schema from the command line.
## Requirements

### Requirement: Parse command validates a GraphQL operation
The system SHALL provide a `parse` subcommand that validates a GraphQL operation document against a schema and reports any errors.

#### Scenario: Valid operation from file
- **WHEN** user runs `gqlxp parse query.graphql` with a syntactically correct operation that matches the schema
- **THEN** the command exits with code 0
- **AND** no errors are printed

#### Scenario: Syntax error in operation
- **WHEN** user runs `gqlxp parse query.graphql` and the file contains a GraphQL syntax error
- **THEN** the command exits with code 1
- **AND** each error is printed to stdout in the format `<source>:<line>:<col>: <message>`

#### Scenario: Unknown field in operation
- **WHEN** user runs `gqlxp parse query.graphql` and the operation references a field that does not exist in the schema
- **THEN** the command exits with code 1
- **AND** an error message identifying the unknown field is printed with its location

#### Scenario: Multiple errors reported
- **WHEN** a GraphQL operation contains multiple validation errors
- **THEN** all errors are printed, one per line
- **AND** the command exits with code 1

### Requirement: Parse command accepts input from stdin
The system SHALL read the operation from stdin when no file argument is provided.

#### Scenario: Valid operation from stdin
- **WHEN** user runs `gqlxp parse` (no file arg) and pipes a valid operation to stdin
- **THEN** the command validates the operation and exits with code 0

#### Scenario: Invalid operation from stdin
- **WHEN** user runs `gqlxp parse` (no file arg) and pipes an invalid operation to stdin
- **THEN** errors are printed using `<stdin>` as the source in the error location prefix
- **AND** the command exits with code 1

### Requirement: Parse command uses standard schema resolution
The system SHALL accept `--schema`/`-s` flag for schema selection, defaulting to the configured default schema when omitted.

#### Scenario: Parse with explicit schema flag
- **WHEN** user runs `gqlxp parse --schema <schema-id> query.graphql`
- **THEN** the operation is validated against the specified schema

#### Scenario: Parse without schema flag
- **WHEN** user runs `gqlxp parse query.graphql` without a schema flag
- **THEN** the operation is validated against the default schema from config
- **AND** if no default schema is configured, an error is returned
