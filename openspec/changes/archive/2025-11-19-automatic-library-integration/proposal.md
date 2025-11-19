# Automatic Library Integration

## Problem
Current library interaction requires explicit CLI flags and subcommands, creating friction in the user workflow. Users must manually manage `library add` commands and remember schema IDs. The library is optional rather than integral to the experience.

## Proposed Solution
Make the library the central workflow by automatically integrating file-based schema loading with library storage. When a schema file is provided, the system checks for existing library entries by file path and hash, prompts for missing metadata, and ensures all schemas are in the library before launching the TUI.

## Expected Outcome
- Seamless library integration without manual `library add` commands
- Automatic detection of schema file changes
- Single source of truth for schema metadata
- Simplified CLI interface focusing on schema files
- All schema exploration sessions backed by library storage

## Scope
This change affects:
- CLI argument parsing and startup flow
- Library metadata storage (add file hash and path)
- User prompting for schema IDs and display names
- Schema loading logic

Out of scope:
- TUI interface changes
- Library file format migration (will enhance existing metadata structure)
- Existing library management commands (keep for manual operations)
