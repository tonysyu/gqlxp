# favorites-ui Specification

## Purpose
Defines correct behavior for toggling and displaying favorites in the schema explorer TUI.

## MODIFIED Requirements

### Requirement: Favorite Types Management
The system SHALL support marking type names and field names as favorites with context-aware storage.

#### Scenario: Add favorite field from top-level panel
- **GIVEN** the current panel displays a top-level GQL type (Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, or Directive)
- **WHEN** a user marks an item as favorite
- **THEN** the item's `RefName()` value is added to the schema's favorites list
- **AND** for Query and Mutation fields, this stores the field name rather than the return type name

#### Scenario: Add favorite type from non-top-level panel
- **GIVEN** the current panel is NOT a top-level GQL type panel
- **WHEN** a user marks an item as favorite
- **THEN** the item's `TypeName()` value is added to the schema's favorites list

#### Scenario: Preserve selection after favorite toggle
- **GIVEN** a user has selected an item in a panel
- **WHEN** the user toggles favorite on that item
- **THEN** the panel refreshes with updated favorite indicators
- **AND** the same item remains selected after the refresh

#### Scenario: Display favorite indicator in panel title
- **GIVEN** a panel displays content for a type or field
- **WHEN** that type or field name exists in the favorites list
- **THEN** the panel title SHALL be prefixed with the favorite icon (★)

#### Scenario: Display favorite indicator in list items
- **GIVEN** a panel displays a list of items
- **WHEN** an item's type name exists in the favorites list
- **THEN** the item title SHALL be prefixed with the favorite icon (★)
