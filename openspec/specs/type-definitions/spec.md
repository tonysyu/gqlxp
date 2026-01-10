# type-definitions Specification

## Purpose
TBD - created by archiving change add-directive-display. Update Purpose after archive.
## Requirements
### Requirement: Type Directive Exposure
All GraphQL type wrappers in the `gql` package SHALL expose directives applied to the type or type member.

#### Scenario: Field directives
- **WHEN** a Field type wraps an ast.FieldDefinition with directives
- **THEN** calling `Field.Directives()` SHALL return wrapped Directive instances

#### Scenario: Object type directives
- **WHEN** an Object type wraps an ast.Definition with type-level directives
- **THEN** calling `Object.Directives()` SHALL return wrapped Directive instances

#### Scenario: Argument directives
- **WHEN** an Argument wraps an ast.ArgumentDefinition with directives
- **THEN** calling `Argument.Directives()` SHALL return wrapped Directive instances

#### Scenario: Enum value directives
- **WHEN** an EnumValue wraps an ast.EnumValueDefinition with directives
- **THEN** calling `EnumValue.Directives()` SHALL return wrapped Directive instances

#### Scenario: All directive-supporting types
- **WHEN** any type that supports directives (Field, Object, Interface, Argument, Enum, EnumValue, InputObject, Union, Scalar) is accessed
- **THEN** it SHALL provide a `Directives()` method returning `[]*Directive`

### Requirement: Directive Formatting in Signatures
Type signatures SHALL include applied directives in their formatted output.

#### Scenario: Field signature with directive
- **WHEN** formatting a field with directives using `Field.FormatSignature()`
- **THEN** the output SHALL include directive syntax (e.g., `fieldName: String @deprecated(reason: "old")`)

#### Scenario: Argument signature with directive
- **WHEN** formatting an argument with directives using `Argument.FormatSignature()`
- **THEN** the output SHALL include directive syntax (e.g., `arg: Int @internal`)

#### Scenario: Multiple directives
- **WHEN** a type member has multiple directives applied
- **THEN** all directives SHALL be included in the signature separated by spaces

#### Scenario: Directive with arguments
- **WHEN** a directive has arguments
- **THEN** the formatted directive SHALL include the arguments (e.g., `@deprecated(reason: "Use newField")`)

#### Scenario: Directive without arguments
- **WHEN** a directive has no arguments
- **THEN** the formatted directive SHALL be name-only (e.g., `@internal`)

