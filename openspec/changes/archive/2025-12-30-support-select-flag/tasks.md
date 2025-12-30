# Tasks: Support Select Flag

## Implementation Tasks

### 1. Add --select flag to CLI app command
- Add `--select` StringFlag to `cli/app.go` appCommand flags
- Set flag alias to `-s` for convenience
- Add usage text: "Pre-select TYPE or TYPE.FIELD in TUI"
- Verify flag appears in `gqlxp app --help` output

**Validation**: Run `gqlxp app --help` and confirm flag is documented

### 2. Parse selection target format
- Create function `parseSelectionTarget(target string) (typeName, fieldName string)` in `cli/app.go`
- Split on "." to separate type and field: `strings.SplitN(target, ".", 2)`
- Return typeName and fieldName (empty if no field specified)
- Handle edge cases: empty string, multiple dots

**Validation**: Unit test with cases: "User", "Query.user", "", "A.B.C"

### 3. Add SelectionTarget to TUI initialization
- Create `SelectionTarget` struct in `tui/xplr/model.go`:
  ```go
  type SelectionTarget struct {
      TypeName  string
      FieldName string
  }
  ```
- Add `StartWithSelection(schema, schemaID, metadata, target)` function to `tui/app.go`
- Pass SelectionTarget through to xplr.Model initialization
- Store target in Model struct temporarily during init

**Validation**: Build succeeds, type compiles

### 4. Implement type category lookup in schema adapter
- Add method `FindTypeCategory(typeName string) (navigation.GQLType, bool)` to `adapters.SchemaView`
- Check each type collection (Objects, Inputs, Enums, etc.) for matching name
- Return corresponding GQLType constant (QueryType for "Query", ObjectType for others, etc.)
- Return `("", false)` if type not found

**Validation**: Unit test with various type names from test schema

### 5. Implement item selection by name in Panel
- Add method `SelectItemByName(name string) bool` to `components.Panel`
- Iterate through `ListModel.Items()` to find item with matching `RefName()`
- Call `ListModel.Select(index)` if found
- Return true if item found and selected, false otherwise

**Validation**: Unit test creating panel with items and selecting by name

### 6. Add ApplySelection method to xplr.Model
- Create `ApplySelection(target SelectionTarget)` method in `xplr/model.go`
- Use `FindTypeCategory` to determine which GQL type contains the type
- Call `nav.SwitchType(gqlType)` to switch to that category
- Populate panel for that type category
- Call `SelectItemByName` on current panel to select the type
- If fieldName provided, trigger navigation forward and select field in next panel
- Handle not-found cases gracefully (do nothing)

**Validation**: Manual test with example schema

### 7. Integrate selection into TUI initialization flow
- In `tui/app.go` StartWithSelection, call `model.ApplySelection(target)` after schema loaded
- Ensure selection applied before program.Run()
- In `cli/app.go`, read `--select` flag value
- Parse into typeName and fieldName
- Pass to `tui.StartWithSelection` instead of `tui.StartWithLibraryData` when flag present

**Validation**: Run `gqlxp app examples/github.graphqls --select User` and verify selection

### 8. Add acceptance test for --select flag
- Create test case in acceptance tests for `--select TypeName`
- Create test case for `--select TypeName.fieldName`
- Verify correct panel is focused and item selected
- Test invalid selections (non-existent type/field)

**Validation**: Run `just test` and verify acceptance tests pass

### 9. Update documentation
- Add `--select` flag to README.md CLI usage examples
- Add example: `gqlxp app schema.graphqls --select Query.user`
- Update `docs/index.md` if CLI documentation exists there

**Validation**: Review README.md changes

## Validation
- [x] All unit tests pass (`just test`)
- [x] Acceptance tests cover --select scenarios
- [x] `gqlxp app --help` shows --select flag
- [x] Manual test: `gqlxp app examples/github.graphqls --select User` selects User
- [x] Manual test: `gqlxp app examples/github.graphqls --select Query.user` selects field
- [x] Invalid selections gracefully fallback without errors
