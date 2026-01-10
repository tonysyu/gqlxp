# Implementation Tasks

## 1. Expose Directives in Type Wrappers
- [x] 1.1 Add `Directives()` method to Field type returning `[]*Directive`
- [x] 1.2 Add `Directives()` method to Object type
- [x] 1.3 Add `Directives()` method to Interface type
- [x] 1.4 Add `Directives()` method to Argument type
- [x] 1.5 Add `Directives()` method to Enum type
- [x] 1.6 Add `Directives()` method to EnumValue type
- [x] 1.7 Add `Directives()` method to InputObject type
- [x] 1.8 Add `Directives()` method to Union type
- [x] 1.9 Add `Directives()` method to Scalar type
- [x] 1.10 Add helper function to convert ast.DirectiveList to wrapped Directives

## 2. Update Signature Formatting
- [x] 2.1 Extend `formatFieldStringWithWidth()` to include directives in field signatures
- [x] 2.2 Extend `getArgumentString()` to include directives in argument signatures
- [x] 2.3 Add helper function to format directive list (e.g., `@deprecated(reason: "old") @custom`)
- [x] 2.4 Update Field.FormatSignature() to include directives inline
- [x] 2.5 Update Argument.FormatSignature() to include directives inline

## 3. Display Directives in Markdown
- [x] 3.1 Update `GenerateFieldMarkdown()` to include directives in signatures
- [x] 3.2 Update `GenerateTypeDefMarkdown()` to include directives for type definitions
- [x] 3.3 Update `FormatFieldDefinitionsWithDescriptions()` to show field directives

## 4. Add Directives Tab in TUI
- [x] 4.1 Create directive list adapter in `tui/xplr/adapters/`
- [x] 4.2 Extend Panel creation logic to detect applied directives
- [x] 4.3 Add "Directives" tab when type/field has applied directives
- [x] 4.4 Support navigation to directive definitions from applied directives
- [x] 4.5 Update tab ordering (put Directives tab in appropriate position)

## 5. Testing and Documentation
- [x] 5.1 Add tests for directive exposure in gql/types_test.go
- [x] 5.2 Add tests for directive formatting in gql/format_test.go
- [x] 5.3 Add tests for markdown generation with directives
- [x] 5.4 Test TUI navigation to directives tab
- [x] 5.5 Update test fixtures if needed (schema files with directives)
- [x] 5.6 Run `just test` to ensure all tests pass
