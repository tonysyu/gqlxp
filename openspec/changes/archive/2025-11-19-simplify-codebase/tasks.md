# Implementation Tasks

## 1. High Priority Simplifications

### 1.1 Library Package API Reduction
- [x] 1.1.1 Make `ConfigDir()` private (`configDir()`) in `library/config.go`
- [x] 1.1.2 Make `SchemasDir()` private (`schemasDir()`) in `library/config.go`
- [x] 1.1.3 Make `MetadataFile()` private (`metadataFile()`) in `library/config.go`
- [x] 1.1.4 Update internal references and tests
- [x] 1.1.5 Run `just test` to verify changes

### 1.2 Navigation Package Encapsulation
- [x] 1.2.1 Make `PanelStack` private (`panelStack`) in `tui/xplr/navigation/stack.go`
- [x] 1.2.2 Make `TypeSelector` private (`typeSelector`) in `tui/xplr/navigation/type_selector.go`
- [x] 1.2.3 Update `NavigationManager` to use private types
- [x] 1.2.4 Update tests in `tui/xplr/navigation/`
- [x] 1.2.5 Run `just test` to verify changes

### 1.3 Utils Reorganization - GraphQL Formatters (SKIPPED)
- Decided not to implement this reorganization

### 1.4 Utils Reorganization - Test Helpers (SKIPPED)
- Decided not to implement this reorganization

## 2. Medium Priority API Cleanup (NOT IMPLEMENTED)

### 2.1 Adapters Package Simplification (NOT IMPLEMENTED)
- Decided not to implement remaining simplifications

### 2.2 GQL Package API Reduction (NOT IMPLEMENTED)
- Decided not to implement remaining simplifications

## 3. Validation (NOT IMPLEMENTED)

### 3.1 Test Coverage (NOT IMPLEMENTED)
- Decided not to implement remaining validation steps

### 3.2 Build Verification (NOT IMPLEMENTED)
- Decided not to implement remaining validation steps

## 4. Documentation

### 4.1 Code Documentation
- [x] 4.1.1 Update package documentation if needed
- [x] 4.1.2 Ensure new files have appropriate package comments

### 4.2 Architecture Documentation
- [x] 4.2.1 Update `docs/architecture.md` if package structure changed significantly
- [x] 4.2.2 Update `docs/coding.md` to reference minimal public API principle
