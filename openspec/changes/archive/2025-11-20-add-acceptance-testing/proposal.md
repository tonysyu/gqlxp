# Change: Add Acceptance Testing

## Why
The project currently has good unit test coverage for individual components (gql parsing, navigation logic, UI components), but lacks high-level acceptance tests that verify complete user workflows. Without acceptance tests:
- Refactoring risks breaking user-facing behavior without detection
- Multi-step interactions (navigation, type switching, panel management) are not systematically verified
- Rendered output (what users actually see) is not comprehensively tested
- Confidence in end-to-end functionality is limited to manual testing

Adding acceptance tests will improve confidence in refactoring, catch integration issues early, and provide living documentation of expected user workflows.

## What Changes
- **Add acceptance test framework**: Create functional, domain-specific helpers for exploring the TUI and verifying screen output
- **Implement test harness**: Build helpers organized by application functionality (navigation, type switching, screen verification)
- **Add example acceptance tests**: Demonstrate testing patterns for core workflows (navigation, type switching, panel interactions, overlay display)
- **Document testing approach**: Update docs to explain when and how to write acceptance tests

Specific additions:
- New package: `tests/acceptance` with test harness and functional helpers
- Test harness components:
  - **Setup utilities**: Schema loading, window configuration, test model initialization
  - **Explorer helpers**: Navigation (cycle types, switch panels, select items), interactions (open overlays, navigate panel stack)
  - **Screen verification helpers**: Panel content assertions, breadcrumb checks, overlay visibility, rendered output validation
- Example tests in `tests/acceptance/workflows_test.go`:
  - Navigate through Query fields to Object types
  - Cycle through GraphQL type categories
  - Open and verify overlay details
  - Multi-panel navigation workflows
- Optional (phase 2): Golden file support for snapshot testing of rendered views

## Impact
- **Breaking changes**: None - purely additive
- **Affected specs**: Introduces new `acceptance-testing` capability spec
- **Affected code**:
  - New: `tests/acceptance/harness.go` - Test framework
  - New: `tests/acceptance/workflows_test.go` - Example tests
  - Update: `docs/coding.md` - Add acceptance testing guidance
  - Update: `docs/development.md` - Document acceptance test commands

## Alternatives Considered
1. **teatest** (official Charm library): Rejected because it's still experimental, adds a dependency, and provides less control over the testing API
2. **catwalk** (data-driven testing): Rejected because data-file approach doesn't align with preference for code-based helper functions
3. **BDD-style Given/When/Then builders**: Rejected in favor of functional, domain-specific helpers that map directly to application features (explorer navigation, screen verification) rather than abstract test phases
4. **Golden files only**: Considered but deferred to phase 2 - will add incrementally after establishing the core framework
