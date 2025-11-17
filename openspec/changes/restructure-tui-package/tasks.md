# Tasks

## 1. Create subpackage directories
- [x] Create `tui/libselect/`, `tui/xplr/`, `tui/overlay/` directories
- [x] Verify structure with `ls -la tui/`

## 2. Move overlay code to tui/overlay/
- [x] Move `tui/overlay.go` → `tui/overlay/model.go`
- [x] Move `tui/overlay_test.go` → `tui/overlay/model_test.go`
- [x] Move `tui/overlay_render_test.go` → `tui/overlay/model_render_test.go`
- [x] Update package declaration from `package tui` to `package overlay`
- [x] Update imports in moved files to use `tui/config` and `tui/components`
- [x] Run `just test` to verify overlay tests pass

## 3. Move schema selector code to tui/libselect/
- [x] Move `tui/schema_selector.go` → `tui/libselect/model.go`
- [x] Update package declaration from `package tui` to `package libselect`
- [x] Update imports to use `tui/adapters`, `tui/config`, `tui/components`
- [x] Rename `StartSchemaSelector()` → `Start()` (exported from libselect)
- [x] Run `just test` to verify no compilation errors

## 4. Move explorer code to tui/xplr/
- [x] Move `tui/model.go` → `tui/xplr/model.go`
- [x] Move `tui/favoritable_item.go` → `tui/xplr/favoritable_item.go`
- [x] Move `tui/model_test.go` → `tui/xplr/model_test.go`
- [x] Move `tui/focus_test.go` → `tui/xplr/focus_test.go`
- [x] Move `tui/breadcrumbs_test.go` → `tui/xplr/breadcrumbs_test.go`
- [x] Move `tui/components/` → `tui/xplr/components/` (explorer-specific UI components)
- [x] Move `tui/navigation/` → `tui/xplr/navigation/` (explorer-specific navigation state)
- [x] Update package declarations from `package tui` to `package xplr`
- [x] Update package declarations in components/ and navigation/ to prefix with `xplr` if needed
- [x] Update imports to use `tui/adapters`, `tui/config`, `tui/xplr/components`, `tui/xplr/navigation`, `tui/overlay`
- [x] Rename `newModel()` → `New()` (exported from xplr)
- [x] Run `just test` to verify all xplr tests pass

## 5. Create top-level delegating model
- [x] Create new `tui/model.go` with delegating model implementation
- [x] Implement state machine that delegates to libselect, xplr, or overlay
- [x] Handle transitions between modes (e.g., SchemaSelectedMsg from libselect → xplr)
- [x] Update `tui/app.go` to use new delegating model
- [x] Ensure `Start()`, `StartWithLibraryData()`, `StartSchemaSelector()` still work

## 6. Fix import paths across codebase
- [x] Update `cmd/gqlxp/main.go` imports if needed (should remain unchanged)
- [x] Update any internal imports between tui subpackages
- [x] Run `just build` to verify all imports resolve correctly

## 7. Run full test suite
- [x] Run `just test` to verify all tests pass
- [x] Run `just verify` to check linting and formatting
- [x] Fix any compilation errors or test failures

## 8. Update documentation
- [x] Update `docs/architecture.md` to reflect new package structure
- [x] Add note about subpackage organization to documentation

## Dependencies
- Tasks 2, 3, 4 can be done in parallel
- Task 5 depends on tasks 2, 3, 4 being complete
- Task 6 depends on task 5
- Tasks 7, 8 depend on task 6
