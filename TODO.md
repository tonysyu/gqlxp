# TODO: Coding Standards Improvements

This file tracks opportunities to improve code quality based on [docs/coding.md](docs/coding.md) standards.

## High Priority: Testing Conventions

### Replace `is.True()` with `is.Equal()` for better error messages

- [x] `tui/overlay/model_test.go:33, 59` - Replace `is.True(cmd == nil)` with `is.Equal(cmd, nil)`
- [x] `tui/model_test.go:79` - Replace `is.True(cmd == nil)` with `is.Equal(cmd, nil)`
- [x] `tui/libselect/model_test.go:87, 105` - Replace `is.True(cmd == nil)` with `is.Equal(cmd, nil)`
- [x] `tui/xplr/components/panels_test.go:173` - Fix tautology `is.True(cmd != nil || cmd == nil)`
- [x] `gql/parse_test.go:570` - Replace `is.True(result == nil)` with `is.Equal(result, nil)`
- [x] `gql/parse_error_handling_test.go:38, 42, 185, 186, 210, 214, 248, 252, 258` - Replace `is.True(err == nil)` with `is.NoErr(err)`
- [x] `tui/adapters/items_test.go:534` - Replace `is.True(panel == nil)` with `is.Equal(panel, nil)`
- [x] `library/config_test.go:18` - Replace `is.True(configDir != "")` with better assertion
- [x] `library/config_test.go:41` - Keep `is.True(info.IsDir())` (boolean values are fine)
- [x] `cli/prompt_test.go:23` - Keep `is.True(err != nil)` (proper error check exists)

### Migrate from `t.Error/t.Errorf` to `is` library

- [x] `tests/fitness/dependency_test.go:127, 129, 131, 133` - Use `is` library
- [x] `tests/fitness/hierarchy_test.go:146, 148, 150` - Use `is` library
- [x] `cli/show_test.go:205, 215, 251, 289` - Use `is` library
- [x] `cli/app_test.go:42, 45` - Use `is` library
- [x] `tui/adapters/schema_test.go:121, 124` - Use `is` library
- [x] `tui/xplr/selection_test.go:66` - Use `is` library
- [x] `tests/acceptance/workflows_test.go:272` - Use `is` library

### Review black-box vs white-box testing

Consider converting these white-box tests (`package foo`) to black-box tests (`package foo_test`) if they're testing public APIs:

- [ ] `cli/show_test.go` - Review if should use `package cli_test`
- [ ] `cli/prompt_test.go` - Review if should use `package cli_test`
- [ ] `cli/app_test.go` - Review if should use `package cli_test`
- [ ] `tui/model_test.go` - Review if should use `package tui_test`
- [ ] `tui/overlay/model_test.go` - Review if should use `package overlay_test`
- [ ] `tui/adapters/items_test.go` - Review if should use `package adapters_test`
- [ ] `tui/adapters/schema_test.go` - Review if should use `package adapters_test`
- [ ] `tui/xplr/focus_test.go` - Review if should use `package xplr_test`
- [ ] `tui/xplr/model_test.go` - Review if should use `package xplr_test`
- [ ] `tui/xplr/selection_test.go` - Review if should use `package xplr_test`
- [ ] `tui/xplr/breadcrumbs_test.go` - Review if should use `package xplr_test`
- [ ] `tui/xplr/components/panels_test.go` - Review if should use `package components_test`
- [ ] `tui/xplr/navigation/stack_test.go` - Review if should use `package navigation_test`
- [ ] `tui/xplr/navigation/manager_test.go` - Review if should use `package navigation_test`
- [ ] `tui/xplr/navigation/type_selector_test.go` - Review if should use `package navigation_test`

## Medium Priority: Line of Sight

Apply guard clauses to keep happy path left-aligned and reduce nesting:

### `tui/xplr/model.go`

- [ ] Lines 214-255: Extract nested key handling logic into separate methods
- [ ] Lines 332-382: Refactor `handleNormal` to reduce nesting depth

### `cli/library.go`

- [ ] Lines 228-291: Apply guard clauses in `addCommand` to reduce if-else chains
- [ ] Lines 320-330: Use early return in `removeCommand` confirmation check

### `library/library.go`

- [ ] Lines 231-267: Simplify `Get` method with early returns
- [ ] Lines 294-310: Simplify `List` method nested conditionals
- [ ] Lines 400-403: Return immediately when metadata match found

### `gqlfmt/markdown.go`

- [ ] Lines 12-60: Use guard clauses instead of if-else chains in `GenerateMarkdown`

## Low Priority: Minimize Public API

### Review struct field visibility

- [ ] `library/types.go` - All fields in `SchemaMetadata`, `Schema`, `SchemaInfo`, `UserConfig` are public. If only read externally, make private with getters.
- [ ] `search/indexer.go:23` - Review if `BleveIndexer` needs to be public
- [ ] `search/searcher.go:11` - Review if `BleveSearcher` needs to be public

### Review exported constants

- [ ] `tui/config/config.go:9-18` - Review if constants (`VisiblePanelCount`, `HelpHeight`, `NavbarHeight`, etc.) need to be exported

## Minor: Error Handling

- [ ] `library/library.go:545, 558` - Consider logging silently ignored indexing errors for debugging
- [ ] `library/library.go:218-220, 491-493` - Add context to cleanup error handling
