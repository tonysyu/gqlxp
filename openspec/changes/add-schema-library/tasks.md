# Implementation Tasks

## 1. Core Library Package
- [ ] 1.1 Create `library` package with schema storage interface
- [ ] 1.2 Implement config directory creation (`~/.config/gqlxp/`)
- [ ] 1.3 Implement schema storage (save/load GraphQL schema files)
- [ ] 1.4 Implement metadata storage (JSON format for schema metadata)
- [ ] 1.5 Define metadata schema (name, favorites, URL patterns)
- [ ] 1.6 Add tests for library package

## 2. Schema Management
- [ ] 2.1 Add schema registration (import from file path)
- [ ] 2.2 Add schema listing functionality
- [ ] 2.3 Add schema removal functionality
- [ ] 2.4 Add schema retrieval by ID/name
- [ ] 2.5 Add tests for schema management

## 3. Metadata Management
- [ ] 3.1 Implement favorite types add/remove operations
- [ ] 3.2 Implement URL pattern storage and retrieval
- [ ] 3.3 Implement schema display name management
- [ ] 3.4 Add tests for metadata operations

## 4. CLI Integration
- [ ] 4.1 Add library mode flag/subcommand support
- [ ] 4.2 Update main.go to support library mode
- [ ] 4.3 Maintain backward compatibility with file-path mode
- [ ] 4.4 Add tests for CLI integration

## 5. TUI Integration (Optional - Future Enhancement)
- [ ] 5.1 Add schema selector screen for library mode
- [ ] 5.2 Add favorites indicator in type lists
- [ ] 5.3 Add URL opening functionality for selected types
- [ ] 5.4 Add library management UI commands

## 6. Documentation
- [ ] 6.1 Update README.md with library mode usage
- [ ] 6.2 Document metadata schema format
- [ ] 6.3 Add architecture documentation for library package
- [ ] 6.4 Update docs/index.md with new documentation references
