# Implementation Tasks

## 1. High Priority Simplifications

### 1.1 Library Package API Reduction
- [ ] 1.1.1 Make `ConfigDir()` private (`configDir()`) in `library/config.go`
- [ ] 1.1.2 Make `SchemasDir()` private (`schemasDir()`) in `library/config.go`
- [ ] 1.1.3 Make `MetadataFile()` private (`metadataFile()`) in `library/config.go`
- [ ] 1.1.4 Update internal references and tests
- [ ] 1.1.5 Run `just test` to verify changes

### 1.2 Navigation Package Encapsulation
- [ ] 1.2.1 Make `PanelStack` private (`panelStack`) in `tui/xplr/navigation/stack.go`
- [ ] 1.2.2 Make `TypeSelector` private (`typeSelector`) in `tui/xplr/navigation/type_selector.go`
- [ ] 1.2.3 Update `NavigationManager` to use private types
- [ ] 1.2.4 Update tests in `tui/xplr/navigation/`
- [ ] 1.2.5 Run `just test` to verify changes

### 1.3 Utils Reorganization - GraphQL Formatters
- [ ] 1.3.1 Create `tui/adapters/formatting.go` for GraphQL-specific formatters
- [ ] 1.3.2 Move `GqlCode()` from `utils/text/formatters.go` to new file
- [ ] 1.3.3 Move `GqlDocString()` from `utils/text/formatters.go` to new file
- [ ] 1.3.4 Move `H1()` from `utils/text/formatters.go` to new file
- [ ] 1.3.5 Update imports in `tui/adapters/items.go` and `tui/xplr/components/items.go`
- [ ] 1.3.6 Run `just test` to verify changes

### 1.4 Utils Reorganization - Test Helpers
- [ ] 1.4.1 Create `tui/xplr/components/testhelpers.go` for component test utilities
- [ ] 1.4.2 Move `RenderMinimalPanel()` from `utils/testx/components.go` to new file
- [ ] 1.4.3 Update imports in component tests
- [ ] 1.4.4 Delete `utils/testx/components.go` if now empty
- [ ] 1.4.5 Run `just test` to verify changes

## 2. Medium Priority API Cleanup

### 2.1 Adapters Package Simplification
- [ ] 2.1.1 Remove `ParseSchemaString()` from `tui/adapters/schema.go`
- [ ] 2.1.2 Update test files to use `ParseSchema([]byte(content))` directly
  - [ ] Update `tui/xplr/model_test.go`
  - [ ] Update `tui/model_test.go`
  - [ ] Update `tui/xplr/breadcrumbs_test.go`
  - [ ] Update `tui/adapters/items_test.go` (if used)
- [ ] 2.1.3 Make `NewSchemaView()` private (`newSchemaView()`) in `tui/adapters/schema.go`
- [ ] 2.1.4 Run `just test` to verify changes

### 2.2 GQL Package API Reduction
- [ ] 2.2.1 Make `Field.DefaultValue()` private in `gql/types.go`
- [ ] 2.2.2 Make `Argument.DefaultValue()` private in `gql/types.go`
- [ ] 2.2.3 Update `gql/helpers.go` to use private methods
- [ ] 2.2.4 Update tests in `gql/` package
- [ ] 2.2.5 Run `just test` to verify changes

## 3. Validation

### 3.1 Test Coverage
- [ ] 3.1.1 Run `just test` for full test suite
- [ ] 3.1.2 Run `just test-coverage` to ensure coverage maintained
- [ ] 3.1.3 Run `just verify` (test + lint + format)

### 3.2 Build Verification
- [ ] 3.2.1 Run `just build` to ensure clean compilation
- [ ] 3.2.2 Manually test the TUI with a sample schema file
- [ ] 3.2.3 Manually test library selection mode

## 4. Documentation

### 4.1 Code Documentation
- [ ] 4.1.1 Update package documentation if needed
- [ ] 4.1.2 Ensure new files have appropriate package comments

### 4.2 Architecture Documentation
- [ ] 4.2.1 Update `docs/architecture.md` if package structure changed significantly
- [ ] 4.2.2 Update `docs/coding.md` to reference minimal public API principle
