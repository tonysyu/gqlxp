# cli-interface Deltas

## ADDED Requirements

### Requirement: Library Command
The system SHALL provide a `library` command that consolidates schema library management functionality.

#### Scenario: Library command help
- **WHEN** user runs `gqlxp library --help`
- **THEN** help text is displayed showing available subcommands

## MODIFIED Requirements

### Requirement: Root Command Default Behavior
The system SHALL maintain backward compatibility by continuing to launch the TUI when no subcommand is specified.

#### Scenario: Root command without arguments
- **WHEN** user runs `gqlxp` with no arguments
- **THEN** the library selector TUI is opened (unchanged from current behavior)

#### Scenario: Root command with schema file
- **WHEN** user runs `gqlxp <schema-file>`
- **THEN** the TUI is opened with the specified schema (unchanged from current behavior)

#### Scenario: Preserve existing subcommands
- **WHEN** user runs `gqlxp search` or `gqlxp show`
- **THEN** the respective subcommand is executed (unchanged from current behavior)

