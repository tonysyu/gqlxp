# Change: Simplify Codebase

## Why
The codebase has accumulated over-exported APIs, minor code organization issues, and some dead code through iterative development. This increases the maintenance burden and makes the public API surface larger than necessary. Simplifying now will improve long-term maintainability and make the codebase easier to understand.

## What Changes
- **Reduce public API surface**: Make internal functions private (lowercase) in library, gql, and tui packages
- **Remove dead code**: Delete unused convenience wrappers and consolidate test helpers
- **Reorganize utilities**: Move GraphQL-specific formatters from generic utils/text to tui/adapters
- **Fix dependency directions**: Move TUI-specific test helpers from utils/testx into tui packages

Specific changes:
- Library package: Make `ConfigDir()`, `SchemasDir()`, `MetadataFile()` private (only used internally)
- Adapters package: Remove `ParseSchemaString()` wrapper, make `NewSchemaView()` private
- GQL package: Make `Field.DefaultValue()`, `Argument.DefaultValue()`, and resolver methods private
- Navigation package: Make `PanelStack` and `TypeSelector` structs private (accessed via `NavigationManager`)
- Utils reorganization: Move `GqlCode()`, `GqlDocString()`, `H1()` from utils/text to tui/adapters, move `RenderMinimalPanel()` from utils/testx to tui/xplr/components

## Impact
- **Breaking changes**: **BREAKING** - Internal APIs become private (affects tests only, not external users since this is a standalone application)
- **Affected specs**: Introduces new `code-quality` capability spec
- **Affected code**:
  - `library/config.go` - 3 functions made private
  - `library/library.go` - Consider encapsulation of `CalculateFileHash()` and `SanitizeSchemaID()`
  - `tui/adapters/schema.go` - Remove 1 function, privatize 1 function
  - `gql/types.go` - Make 4 methods private
  - `tui/xplr/navigation/` - Make 2 types private
  - `utils/text/formatters.go` - Move 3 functions to tui/adapters
  - `utils/testx/components.go` - Move 1 function to tui/xplr/components
  - All test files using affected functions
