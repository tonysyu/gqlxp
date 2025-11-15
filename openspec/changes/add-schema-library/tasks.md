# Implementation Tasks

## 1. Core Library Package
- [x] 1.1 Create `library` package with schema storage interface
- [x] 1.2 Implement config directory creation (`~/.config/gqlxp/`)
- [x] 1.3 Implement schema storage (save/load GraphQL schema files)
- [x] 1.4 Implement metadata storage (single `schemas/metadata.json` file with schema-id keys)
- [x] 1.5 Define metadata schema (name, favorites, URL patterns)
- [x] 1.6 Add tests for library package

## 2. Schema Management
- [x] 2.1 Add schema registration (import from file path)
- [x] 2.2 Add schema listing functionality
- [x] 2.3 Add schema removal functionality
- [x] 2.4 Add schema retrieval by ID/name
- [x] 2.5 Add tests for schema management

## 3. Metadata Management
- [x] 3.1 Implement favorite types add/remove operations
- [x] 3.2 Implement URL pattern storage and retrieval
- [x] 3.3 Implement schema display name management
- [x] 3.4 Add tests for metadata operations

## 4. CLI Integration
- [x] 4.1 Add library mode flag/subcommand support
- [x] 4.2 Update main.go to support library mode
- [x] 4.3 Maintain backward compatibility with file-path mode
- [x] 4.4 Add tests for CLI integration

## 5. TUI Integration
- [x] 5.1 Add schema selector screen for library mode
- [x] 5.2 Add favorites indicator in type lists
- [x] 5.3 Add library management UI commands (favorites toggling)

## 6. Documentation
- [x] 6.1 Update README.md with library mode usage
- [x] 6.2 Document metadata schema format
- [x] 6.3 Add architecture documentation for library package
- [x] 6.4 Update docs/index.md with new documentation references
