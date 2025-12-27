# code-quality Specification Deltas

## MODIFIED Requirements

### Requirement: Dependency Direction
The codebase SHALL maintain proper dependency direction following the package hierarchy: `cmd → cli → tui → library → search → gql → utils`, with lower-level packages not depending on higher-level packages.

#### Scenario: Utils to domain dependency
- **WHEN** creating code in a `utils/` package
- **THEN** it SHALL NOT depend on domain packages like `tui/`, `gql/`, `library/`, `search/`, or `cli/`

#### Scenario: Test helpers dependency
- **WHEN** test helper code requires domain types or packages
- **THEN** it SHALL be placed in the domain package rather than creating a reverse dependency from `utils/`

#### Scenario: Search package dependencies
- **WHEN** creating code in the `search/` package
- **THEN** it SHALL only import from `gql` and `utils` packages, not from `library`, `tui`, `cli`, or `cmd`

#### Scenario: Library package using search
- **WHEN** the `library` package needs search functionality
- **THEN** it SHALL import the `search` package as allowed by the hierarchy
