# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`gq` is a GraphQL query explorer TUI (Terminal User Interface) application that parses GraphQL schema files and provides an interactive terminal interface for exploring GraphQL queries. The application uses the Bubble Tea framework to create a multi-panel interface where users can navigate through GraphQL schema definitions.

## Development Commands

### Building and Running
- `go build -o gq ./cmd/gq` - Build the application
- `go run ./cmd/gq` - Run the application (requires `examples/github.graphqls` schema file)

### Testing
- `go test ./... -v` - Run all tests with verbose output
- `go test ./gql -v` - Run only GraphQL parsing tests

### Code Quality
- `go fmt ./...` - Format Go code
- `go vet ./...` - Run Go static analysis
- `go mod tidy` - Clean up module dependencies

## Architecture

### Package Structure
- **`cmd/gq`**: Main application entry point that reads schema file and initializes TUI
- **`gql`**: GraphQL schema parsing and type extraction
  - `ParseSchema()` - Main function to parse GraphQL schema and extract Query fields
  - `buildGraphQLTypes()` - Builds map of GraphQL types from AST
  - `GetTypeString()` - Converts AST types to string representation
- **`tui`**: Terminal user interface components built on Bubble Tea
  - `mainModel` - Root application model managing panels and navigation
  - `Panel` interface - Generic panel abstraction for different content types
  - `ListItem` interface - Interactive list items that can be "opened" to show details
  - `listPanel` and `stringPanel` - Concrete panel implementations

### Key Interfaces
- **`ListItem`**: Extends `list.DefaultItem` with `Open() Panel` method for interactive items
- **`Panel`**: Implements `tea.Model` with `SetSize(width, height int)` for resizable content

### Navigation Flow
1. Application starts by parsing GraphQL schema from `examples/github.graphqls`
2. Creates list of GraphQL query fields as interactive items
3. User navigates with tab/shift+tab between panels
4. Pressing enter on list items opens details in adjacent panels
5. Supports up to 6 panels horizontally with ctrl+n/ctrl+w for add/remove

### Dependencies
- **Bubble Tea**: TUI framework for terminal applications
- **Bubbles**: Pre-built components (list, help) for Bubble Tea
- **Lipgloss**: Terminal styling and layout
- **graphql-go**: GraphQL parsing and AST manipulation

## Schema Requirements

The application expects a GraphQL schema file at `examples/github.graphqls`. The parser specifically looks for Query type definitions and extracts field information including arguments, return types, and descriptions.