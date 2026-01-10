## ADDED Requirements
### Requirement: Directives Tab Display
Panels SHALL display a "Directives" tab when the selected type or field has applied directives, showing both directive applications and allowing navigation to directive definitions.

#### Scenario: Field with applied directives
- **WHEN** a field panel is opened for a field with applied directives (e.g., `@deprecated`)
- **THEN** the panel SHALL include a "Directives" tab alongside existing tabs
- **AND** the tab SHALL display applied directives with their arguments

#### Scenario: Type with applied directives
- **WHEN** a type definition panel is opened for a type with applied directives
- **THEN** the panel SHALL include a "Directives" tab
- **AND** the tab SHALL display applied directives at the type level

#### Scenario: Navigation to directive definition
- **WHEN** a user selects an applied directive from the Directives tab
- **THEN** selecting it SHALL navigate to the directive definition
- **AND** SHALL open the directive's detail panel

#### Scenario: No directives applied
- **WHEN** a type or field has no applied directives
- **THEN** the panel SHALL NOT include a Directives tab

#### Scenario: Tab ordering
- **WHEN** a panel has multiple tabs including Directives
- **THEN** the Directives tab SHALL appear after primary relationship tabs (Type, Inputs, Fields, etc.)
