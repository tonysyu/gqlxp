# Implementation Tasks

## 1. Expose Directives in Type Wrappers
- [ ] 1.1 Add `Directives()` method to Field type returning `[]*Directive`
- [ ] 1.2 Add `Directives()` method to Object type
- [ ] 1.3 Add `Directives()` method to Interface type
- [ ] 1.4 Add `Directives()` method to Argument type
- [ ] 1.5 Add `Directives()` method to Enum type
- [ ] 1.6 Add `Directives()` method to EnumValue type
- [ ] 1.7 Add `Directives()` method to InputObject type
- [ ] 1.8 Add `Directives()` method to Union type
- [ ] 1.9 Add `Directives()` method to Scalar type
- [ ] 1.10 Add helper function to convert ast.DirectiveList to wrapped Directives

## 2. Update Signature Formatting
- [ ] 2.1 Extend `formatFieldStringWithWidth()` to include directives in field signatures
- [ ] 2.2 Extend `getArgumentString()` to include directives in argument signatures
- [ ] 2.3 Add helper function to format directive list (e.g., `@deprecated(reason: "old") @custom`)
- [ ] 2.4 Update Field.FormatSignature() to include directives inline
- [ ] 2.5 Update Argument.FormatSignature() to include directives inline

## 3. Display Directives in Markdown
- [ ] 3.1 Update `GenerateFieldMarkdown()` to include directives in signatures
- [ ] 3.2 Update `GenerateTypeDefMarkdown()` to include directives for type definitions
- [ ] 3.3 Update `FormatFieldDefinitionsWithDescriptions()` to show field directives

## 4. Add Directives Tab in TUI
- [ ] 4.1 Create directive list adapter in `tui/xplr/adapters/`
- [ ] 4.2 Extend Panel creation logic to detect applied directives
- [ ] 4.3 Add "Directives" tab when type/field has applied directives
- [ ] 4.4 Support navigation to directive definitions from applied directives
- [ ] 4.5 Update tab ordering (put Directives tab in appropriate position)

## 5. Testing and Documentation
- [ ] 5.1 Add tests for directive exposure in gql/types_test.go
- [ ] 5.2 Add tests for directive formatting in gql/format_test.go
- [ ] 5.3 Add tests for markdown generation with directives
- [ ] 5.4 Test TUI navigation to directives tab
- [ ] 5.5 Update test fixtures if needed (schema files with directives)
- [ ] 5.6 Run `just test` to ensure all tests pass
