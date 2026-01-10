# Change: Expose and Display GraphQL Directives

## Why
GraphQL schemas use directives to provide metadata and control behavior (e.g., `@deprecated`, `@auth`, `@skip`). The vektah/gqlparser already collects directives for fields, objects, arguments, and other types, but gqlxp currently doesn't expose or display this information. Users exploring schemas need to see directives to fully understand type definitions and field behaviors.

## What Changes
- Expose directives in all `gql` type wrappers (Field, Object, Interface, Argument, Enum, EnumValue, InputObject, Union, Scalar)
- Display directives inline with signatures in markdown output (e.g., `fieldName: String @deprecated(reason: "Use newField")`)
- Add a "Directives" tab in TUI panels showing:
  - Applied directives for the current type/field
  - Navigation to directive definitions
- Update signature formatting to include directive syntax

## Impact
- Affected specs: `schema-library`, `tui-search-tab`
- Affected code:
  - `gql/types.go` - Add `Directives()` methods to all type wrappers
  - `gql/format.go` - Update signature formatting to include directives
  - `gqlfmt/markdown.go` - Include directives in signatures
  - `tui/xplr/components/panels.go` - Add directives tab when directives exist
  - `tui/xplr/adapters/*.go` - Create adapters for directive list items
