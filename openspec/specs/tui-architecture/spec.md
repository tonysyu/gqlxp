# tui-architecture Specification

## Purpose
TBD - created by archiving change restructure-tui-package. Update Purpose after archive.
## Requirements
### Requirement: TUI Package Organization
The TUI package SHALL organize code into subpackages by UI mode to separate concerns and improve maintainability.

#### Scenario: Library selection mode
- **WHEN** code implements library selection functionality
- **THEN** it SHALL be located in `tui/libselect/` subpackage

#### Scenario: Schema exploration mode
- **WHEN** code implements schema exploration functionality
- **THEN** it SHALL be located in `tui/xplr/` subpackage

#### Scenario: Overlay display mode
- **WHEN** code implements overlay display functionality
- **THEN** it SHALL be located in `tui/overlay/` subpackage

#### Scenario: Shared utilities
- **WHEN** code provides shared utilities (adapters, config)
- **THEN** it SHALL remain at the `tui/` package level

#### Scenario: Explorer-specific utilities
- **WHEN** code provides explorer-specific utilities (components, navigation)
- **THEN** it SHALL be located under `tui/xplr/` subpackage

### Requirement: Top-Level Model Delegation
The TUI package SHALL provide a top-level model that delegates to appropriate submode models.

#### Scenario: Mode delegation
- **WHEN** the TUI is started in a specific mode
- **THEN** the top-level model SHALL delegate to the corresponding subpackage model

#### Scenario: Mode transitions
- **WHEN** transitioning between UI modes (e.g., library selection to schema exploration)
- **THEN** the top-level model SHALL handle the transition and activate the appropriate submode

### Requirement: Backward Compatible API
The TUI package SHALL maintain backward compatibility for all public entry points.

#### Scenario: Start function
- **WHEN** `tui.Start()` is called with a schema
- **THEN** it SHALL start the schema exploration mode as before

#### Scenario: StartWithLibraryData function
- **WHEN** `tui.StartWithLibraryData()` is called with schema and metadata
- **THEN** it SHALL start the schema exploration mode with library data as before

#### Scenario: StartSchemaSelector function
- **WHEN** `tui.StartSchemaSelector()` is called
- **THEN** it SHALL start the library selection mode as before

