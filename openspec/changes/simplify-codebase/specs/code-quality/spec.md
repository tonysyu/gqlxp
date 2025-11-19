# code-quality Specification Deltas

## ADDED Requirements

### Requirement: Minimal Public API
The codebase SHALL minimize the public API surface by making functions, types, and methods private (lowercase) when they are only used within their defining package or subpackages.

#### Scenario: Internal helper function
- **WHEN** a function is only called within the same package
- **THEN** the function SHALL be private (lowercase name)

#### Scenario: Package-scoped type
- **WHEN** a type is only instantiated and used within the same package
- **THEN** the type SHALL be private (lowercase name)

#### Scenario: Cross-package dependency
- **WHEN** a function or type is used by other packages
- **THEN** the function or type SHALL be public (uppercase name) with clear documentation

#### Scenario: Test-only usage
- **WHEN** a function is only used in test files
- **THEN** the function SHALL be private and tests SHALL use the public API or white-box testing approach

### Requirement: Code Organization by Domain
The codebase SHALL organize code by domain concern, placing domain-specific code in appropriate packages rather than generic utility packages.

#### Scenario: Domain-specific utilities
- **WHEN** utility functions are specific to a domain (e.g., GraphQL formatting)
- **THEN** they SHALL be placed in the appropriate domain package (e.g., `tui/adapters`) rather than generic `utils/`

#### Scenario: Generic utilities
- **WHEN** utility functions are truly generic and reusable across domains
- **THEN** they SHALL be placed in `utils/` packages

#### Scenario: Test helpers
- **WHEN** test helper functions depend on specific packages or types
- **THEN** they SHALL be co-located with the code they test rather than in generic `utils/testx/`

### Requirement: Dependency Direction
The codebase SHALL maintain proper dependency direction, with lower-level packages not depending on higher-level packages.

#### Scenario: Utils to domain dependency
- **WHEN** creating code in a `utils/` package
- **THEN** it SHALL NOT depend on domain packages like `tui/`, `gql/`, or `library/`

#### Scenario: Test helpers dependency
- **WHEN** test helper code requires domain types or packages
- **THEN** it SHALL be placed in the domain package rather than creating a reverse dependency from `utils/`

### Requirement: Dead Code Removal
The codebase SHALL remove unused code, unnecessary wrappers, and redundant abstractions to maintain clarity and reduce maintenance burden.

#### Scenario: Trivial wrapper function
- **WHEN** a function is a simple wrapper with no added value
- **THEN** the wrapper SHALL be removed and callers SHALL use the underlying function directly

#### Scenario: Unused export
- **WHEN** an exported function, type, or method is never used outside its package
- **THEN** it SHALL be made private or removed if completely unused

#### Scenario: Redundant abstraction
- **WHEN** an abstraction adds complexity without providing benefit
- **THEN** the abstraction SHALL be removed in favor of simpler, direct code
