# interface-navigation Specification

## Purpose
Enable users to discover and navigate to interface definitions from Object types by providing an "Interfaces" subtab in Object type panels.

## ADDED Requirements

### Requirement: Object Type Interface Display
Object type panels SHALL display an "Interfaces" subtab listing all interfaces implemented by the object.

#### Scenario: Object with multiple interfaces
- **WHEN** viewing an Object type that implements multiple interfaces (e.g., `User implements Node & Timestamped`)
- **THEN** the panel SHALL display a "Fields" tab and an "Interfaces" tab
- **AND** the "Interfaces" tab SHALL contain list items for each implemented interface
- **AND** the list items SHALL be ordered alphabetically by interface name

#### Scenario: Object with single interface
- **WHEN** viewing an Object type that implements exactly one interface
- **THEN** the panel SHALL display both "Fields" and "Interfaces" tabs

#### Scenario: Object with no interfaces
- **WHEN** viewing an Object type that implements no interfaces
- **THEN** the panel SHALL display only the "Fields" tab without an "Interfaces" tab

#### Scenario: Navigate to interface from object
- **WHEN** selecting an interface in the "Interfaces" tab and pressing Enter
- **THEN** a new panel SHALL open displaying the selected interface's details
- **AND** the interface panel SHALL show its fields

### Requirement: Interface Resolution
The system SHALL resolve interface names to navigable type definitions for Objects.

#### Scenario: Resolve interface name to definition
- **WHEN** an Object lists an interface name in its `Interfaces()` output
- **THEN** the system SHALL use TypeResolver to resolve the name to the full Interface TypeDef
- **AND** SHALL create a typeDefItem for that interface

#### Scenario: Handle unresolvable interface
- **WHEN** an interface name cannot be resolved by TypeResolver
- **THEN** the system SHALL create a simple non-navigable item with the interface name
- **AND** SHALL not crash or show an error to the user

### Requirement: Tab Ordering and Visibility
Interface tabs SHALL be positioned consistently and follow existing tab navigation patterns.

#### Scenario: Tab order for Objects
- **WHEN** an Object type panel has both Fields and Interfaces
- **THEN** tabs SHALL be ordered: "Fields", "Interfaces"
- **AND** the "Fields" tab SHALL be active by default

#### Scenario: Tab navigation consistency
- **WHEN** navigating between tabs using Shift+H/Shift+L
- **THEN** navigation SHALL follow the existing tab navigation behavior
- **AND** SHALL show the appropriate help text for multi-tab panels
