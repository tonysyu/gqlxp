# acceptance-testing Specification Deltas

## ADDED Requirements

### Requirement: Functional Setup Utilities
The codebase SHALL provide setup utilities for initializing test harnesses with schema content, window configuration, and initial application state.

#### Scenario: Harness initialization
- **WHEN** creating a test harness
- **THEN** the `New()` constructor SHALL accept functional options for configuration
- **AND** options SHALL include schema content (`WithSchema()`) and window size (`WithWindowSize()`)

#### Scenario: Schema loading
- **WHEN** configuring test schema
- **THEN** utilities SHALL support loading from schema strings and test fixtures
- **AND** schema parsing SHALL match production parsing behavior

#### Scenario: Model initialization
- **WHEN** initializing the test model
- **THEN** the harness SHALL create a properly configured model with window size and initial state
- **AND** the model SHALL be ready for navigation and interaction

### Requirement: Explorer Helper Methods
The codebase SHALL provide helper methods for navigating and interacting with the TUI explorer that map directly to user actions.

#### Scenario: Type cycling
- **WHEN** testing type exploration features
- **THEN** helpers SHALL include `CycleTypeForward()`, `CycleTypeBackward()`, and `SwitchToType(gqlType)`
- **AND** methods SHALL simulate keyboard shortcuts (Ctrl+T, Ctrl+R)

#### Scenario: Panel navigation
- **WHEN** testing panel navigation
- **THEN** helpers SHALL include `NavigateToNextPanel()` and `NavigateToPreviousPanel()`
- **AND** methods SHALL simulate Tab and Shift+Tab keyboard navigation

#### Scenario: Item selection
- **WHEN** testing item selection and exploration
- **THEN** helpers SHALL include `SelectItem(name)` and `SelectItemAtIndex(idx)`
- **AND** selection SHALL trigger the same panel-opening behavior as in the real application

#### Scenario: Overlay interactions
- **WHEN** testing overlay detail views
- **THEN** helpers SHALL include `OpenOverlay()` and `CloseOverlay()`
- **AND** methods SHALL simulate spacebar and ESC key interactions

### Requirement: Screen Verification Helpers
The codebase SHALL provide helper methods for verifying rendered screen output organized by screen region and display type.

#### Scenario: Panel content verification
- **WHEN** verifying panel content
- **THEN** helpers SHALL include `AssertPanelContains(panelIdx, text)` and `AssertPanelEquals(panelIdx, expected)`
- **AND** assertions SHALL check rendered panel output from `View()`

#### Scenario: Breadcrumb verification
- **WHEN** verifying breadcrumb display
- **THEN** helpers SHALL include `AssertBreadcrumbsShow(text)` and `AssertBreadcrumbsEmpty()`
- **AND** assertions SHALL verify breadcrumb line content in rendered view

#### Scenario: Overlay verification
- **WHEN** verifying overlay display
- **THEN** helpers SHALL include `AssertOverlayVisible()` and `AssertOverlayContains(text)`
- **AND** assertions SHALL check overlay presence and content in rendered view

#### Scenario: General view assertions
- **WHEN** verifying overall application state
- **THEN** helpers SHALL include `AssertViewContains(text)` and `AssertCurrentType(gqlType)`
- **AND** assertions SHALL use normalized views for reliable comparisons

### Requirement: User-Facing Output Verification
The acceptance test framework SHALL verify rendered terminal output (what users actually see) rather than internal model state.

#### Scenario: View content assertions
- **WHEN** asserting on application output
- **THEN** assertions SHALL check the rendered view string (from `Model.View()`)
- **AND** assertions SHALL use normalized views for reliable comparisons

#### Scenario: Rendered element verification
- **WHEN** verifying UI elements are displayed
- **THEN** tests SHALL check for presence of text, breadcrumbs, panel content, and overlays in rendered output
- **AND** tests SHALL NOT directly access internal model fields for verification

#### Scenario: Visual state checks
- **WHEN** verifying visual application state
- **THEN** tests SHALL check overlay visibility, panel count, and breadcrumb display through rendered view
- **AND** checks SHALL reflect what users would observe in the terminal

### Requirement: Complete Workflow Coverage
The acceptance tests SHALL cover complete, multi-step user workflows that represent real usage patterns.

#### Scenario: Navigation workflow
- **WHEN** testing navigation features
- **THEN** tests SHALL verify multi-step workflows: loading schema → selecting items → opening panels → navigating panel stack
- **AND** tests SHALL verify breadcrumb updates, panel focus changes, and content display

#### Scenario: Type switching workflow
- **WHEN** testing GraphQL type exploration
- **THEN** tests SHALL verify cycling through all type categories (Query, Mutation, Object, Input, Enum, Scalar, Interface, Union, Directive)
- **AND** tests SHALL verify type-specific content loads correctly and panel state resets appropriately

#### Scenario: Overlay interaction workflow
- **WHEN** testing detail overlay feature
- **THEN** tests SHALL verify opening overlay, displaying item details, and closing overlay
- **AND** tests SHALL verify correct markdown rendering and return to main view state

#### Scenario: Complex multi-step workflow
- **WHEN** testing realistic usage scenarios
- **THEN** tests SHALL combine multiple features: navigation + type switching + overlay display
- **AND** tests SHALL verify state consistency across feature interactions

### Requirement: Test Organization and Reusability
The acceptance test framework SHALL promote code reuse and maintainability through well-organized, domain-specific helpers.

#### Scenario: Explorer helper organization
- **WHEN** implementing navigation and interaction helpers
- **THEN** methods SHALL be grouped by functionality (type cycling, panel navigation, item selection, overlay interaction)
- **AND** method names SHALL clearly describe the user action being simulated

#### Scenario: Screen verification organization
- **WHEN** implementing verification helpers
- **THEN** methods SHALL be organized by screen region (panels, breadcrumbs, overlays, full view)
- **AND** assertion names SHALL clearly describe what is being verified

#### Scenario: Schema test fixtures
- **WHEN** setting up test schemas
- **THEN** common schema patterns SHALL be available as reusable fixtures
- **AND** tests SHALL be able to create custom schemas for specific scenarios

#### Scenario: Test file organization
- **WHEN** organizing acceptance tests
- **THEN** tests SHALL be grouped by workflow category in `tests/acceptance/workflows_test.go`
- **AND** test names SHALL clearly describe the workflow being verified

### Requirement: Integration with Existing Test Infrastructure
The acceptance test framework SHALL integrate seamlessly with existing testing patterns and tools.

#### Scenario: Build on testModel pattern
- **WHEN** implementing the test harness
- **THEN** the harness SHALL extend the existing `testModel` wrapper from `tui/xplr/breadcrumbs_test.go`
- **AND** the harness SHALL reuse existing test utilities from `utils/testx`

#### Scenario: Use standard assertions
- **WHEN** implementing test assertions
- **THEN** assertions SHALL use `github.com/matryer/is` consistent with existing tests
- **AND** error messages SHALL follow existing assertion patterns

#### Scenario: Run with existing commands
- **WHEN** running acceptance tests
- **THEN** tests SHALL execute via standard `just test` command
- **AND** tests SHALL integrate with coverage reporting via `just test-coverage`

#### Scenario: No new dependencies
- **WHEN** implementing the framework
- **THEN** implementation SHALL use only existing project dependencies
- **AND** no external testing libraries (teatest, catwalk, ginkgo, etc.) SHALL be required
