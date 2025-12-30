# cli-interface Specification

## Purpose
TBD - created by archiving change add-app-subcommand. Update Purpose after archive.
## Requirements
### Requirement: App Subcommand
The system SHALL provide an explicit `app` subcommand that launches the TUI with identical behavior to the root command default action.

#### Scenario: Launch app without arguments
- **WHEN** user runs `gqlxp app` with no arguments
- **THEN** the library selector TUI is opened (same as `gqlxp` with no arguments)

#### Scenario: Launch app with schema file
- **WHEN** user runs `gqlxp app <schema-file>`
- **THEN** the TUI is opened with the specified schema (same as `gqlxp <schema-file>`)

#### Scenario: App subcommand accepts log-file flag
- **WHEN** user runs `gqlxp app --log-file debug.log`
- **THEN** debug logging is enabled to the specified file

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

### Requirement: Shared Implementation
The system SHALL implement the app subcommand and root default action using shared logic to ensure consistent behavior.

#### Scenario: Code reuse
- **WHEN** the app subcommand or root default action needs to launch the TUI
- **THEN** both SHALL use the same underlying function to ensure identical behavior

#### Scenario: Flag consistency
- **WHEN** flags are added to the root command
- **THEN** those flags SHALL be available to the app subcommand

### Requirement: Library Command
The system SHALL provide a `library` command that consolidates schema library management functionality.

#### Scenario: Library command help
- **WHEN** user runs `gqlxp library --help`
- **THEN** help text is displayed showing available subcommands

