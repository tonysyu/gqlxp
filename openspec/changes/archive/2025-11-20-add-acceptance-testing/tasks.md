# Implementation Tasks

## 1. Foundation - Test Harness Core

### 1.1 Create Acceptance Test Package
- [x] 1.1.1 Create `tests/acceptance/` directory
- [x] 1.1.2 Add package documentation explaining purpose and usage patterns
- [x] 1.1.3 Verify package structure with `go list ./tests/acceptance`

### 1.2 Implement Setup Utilities
- [x] 1.2.1 Create `harness.go` with `Harness` struct extending `testModel` pattern
- [x] 1.2.2 Implement `New()` constructor with functional options (`WithSchema`, `WithWindowSize`)
- [x] 1.2.3 Add schema loading utilities (parse schema string, load from fixture)
- [x] 1.2.4 Add model initialization with proper window size and initial state
- [x] 1.2.5 Add documentation for setup utilities
- [x] 1.2.6 Run `just test ./tests/acceptance` to verify compilation

### 1.3 Implement Explorer Helpers
- [x] 1.3.1 Add type cycling methods (`CycleTypeForward()`, `CycleTypeBackward()`, `SwitchToType(gqlType)`)
- [x] 1.3.2 Add panel navigation methods (`NavigateToNextPanel()`, `NavigateToPreviousPanel()`)
- [x] 1.3.3 Add item selection methods (`SelectItem(name)`, `SelectItemAtIndex(idx)`)
- [x] 1.3.4 Add overlay interaction methods (`OpenOverlay()`, `CloseOverlay()`)
- [x] 1.3.5 Add documentation for each explorer method
- [x] 1.3.6 Run `just test ./tests/acceptance` to verify helpers

### 1.4 Implement Screen Verification Helpers
- [x] 1.4.1 Add panel content assertions (`AssertPanelContains(panelIdx, text)`, `AssertPanelEquals(panelIdx, expected)`)
- [x] 1.4.2 Add breadcrumb assertions (`AssertBreadcrumbsShow(text)`, `AssertBreadcrumbsEmpty()`)
- [x] 1.4.3 Add overlay assertions (`AssertOverlayVisible()`, `AssertOverlayContains(text)`)
- [x] 1.4.4 Add general view assertions (`AssertViewContains(text)`, `AssertCurrentType(gqlType)`)
- [x] 1.4.5 Use `matryer/is` for assertion failures with clear error messages
- [x] 1.4.6 Run `just test ./tests/acceptance` to verify assertions

## 2. Example Tests - Demonstrate Patterns

### 2.1 Navigation Workflow Tests
- [x] 2.1.1 Create `workflows_test.go` in `tests/acceptance/`
- [x] 2.1.2 Write test: Navigate from Query fields to Object type (using `SelectItem()`, `NavigateToNextPanel()`, `AssertPanelContains()`)
- [x] 2.1.3 Write test: Navigate through multiple panels with breadcrumbs (verify breadcrumb updates with `AssertBreadcrumbsShow()`)
- [x] 2.1.4 Write test: Navigate backward through panel stack (using `NavigateToPreviousPanel()`)
- [x] 2.1.5 Run tests with `just test ./tests/acceptance -v`

### 2.2 Type Switching Tests
- [x] 2.2.1 Write test: Cycle forward through all GraphQL types (using `CycleTypeForward()`, `AssertCurrentType()`)
- [x] 2.2.2 Write test: Cycle backward through GraphQL types (using `CycleTypeBackward()`)
- [x] 2.2.3 Write test: Switch directly to specific type (using `SwitchToType()`)
- [x] 2.2.4 Run tests with `just test ./tests/acceptance -v`

### 2.3 Overlay Interaction Tests
- [x] 2.3.1 Write test: Open overlay and verify content (using `OpenOverlay()`, `AssertOverlayVisible()`, `AssertOverlayContains()`)
- [x] 2.3.2 Write test: Close overlay and return to main view (using `CloseOverlay()`)
- [x] 2.3.3 Write test: Verify overlay shows correct details for different item types
- [x] 2.3.4 Run tests with `just test ./tests/acceptance -v`

### 2.4 Complex Workflow Tests
- [x] 2.4.1 Write test: Full exploration workflow combining navigation, type switching, and overlay interactions
- [x] 2.4.2 Write test: Multi-panel navigation with type cycling (verify breadcrumbs and panel content throughout)
- [x] 2.4.3 Write test: Edge cases (empty schemas, single-item lists, switching types with empty results)
- [x] 2.4.4 Run tests with `just test ./tests/acceptance -v`

## 3. Integration and Validation

### 3.1 Test Coverage Verification
- [x] 3.1.1 Run `just test` to ensure all tests pass
- [x] 3.1.2 Run `just test-coverage` to verify coverage metrics
- [x] 3.1.3 Review coverage report for acceptance package
- [x] 3.1.4 Ensure no regressions in existing tests

### 3.2 Build and Lint Verification
- [x] 3.2.1 Run `just build` to verify compilation
- [x] 3.2.2 Run `just lint-fix` to ensure code quality
- [x] 3.2.3 Run `just verify` for full validation
- [x] 3.2.4 Fix any issues identified by linters

## 4. Documentation

### 4.1 Update Testing Documentation
- [x] 4.1.1 Add acceptance testing section to `docs/coding.md`
- [x] 4.1.2 Document when to use acceptance vs unit tests
- [x] 4.1.3 Provide examples of good acceptance test scenarios
- [x] 4.1.4 Update `docs/development.md` with acceptance test commands

### 4.2 Add Inline Documentation
- [x] 4.2.1 Add package-level documentation to `tests/acceptance/`
- [x] 4.2.2 Document public methods in harness builders
- [x] 4.2.3 Add usage examples in code comments
- [x] 4.2.4 Verify documentation with `go doc tests/acceptance`

### 4.3 Update Architecture Documentation
- [x] 4.3.1 Add `tests/acceptance` to package list in `docs/architecture.md`
- [x] 4.3.2 Describe acceptance testing approach in architecture overview
- [x] 4.3.3 Link to coding.md for detailed guidance

## 5. Optional - Golden File Support (Phase 2)

### 5.1 Golden File Infrastructure
- [x] 5.1.1 Add golden file utilities to `tests/acceptance/golden.go`
- [x] 5.1.2 Implement snapshot comparison with normalized views
- [x] 5.1.3 Add `-update` flag support for regenerating golden files
- [x] 5.1.4 Create `testdata/` directory for golden files

### 5.2 Golden File Tests
- [x] 5.2.1 Convert selected workflow tests to use golden files
- [x] 5.2.2 Add tests for visual regression detection
- [x] 5.2.3 Document golden file testing approach
- [x] 5.2.4 Run tests with `just test ./tests/acceptance -update` to generate baselines

## Notes
- Tasks can be parallelized within sections (e.g., 2.1-2.4 can run in parallel)
- Phase 2 (golden files) is optional and can be deferred based on initial results
- Each task should result in verifiable, working code with passing tests
