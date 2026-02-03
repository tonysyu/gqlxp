package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tonysyu/gqlxp/library"
)

func loadSchemaFromFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}
	return content, nil
}

// resolveSchemaFromArgument resolves a schema argument to a Schema.
// The argument can be:
// 1. Empty string - use default schema from config
// 2. A schema ID that exists in the library
// 3. A file path (will be added to library if needed)
func resolveSchemaFromArgument(arg string) (*library.Schema, error) {
	lib := library.NewLibrary()

	var schemaID string

	// Empty argument - use default schema
	if arg == "" {
		defaultSchemaID, err := lib.GetDefaultSchema()
		if err != nil {
			return nil, fmt.Errorf("error getting default schema: %w", err)
		}
		if defaultSchemaID == "" {
			return nil, fmt.Errorf("no schema specified and no default schema set. Use 'gqlxp library default' to set one")
		}
		schemaID = defaultSchemaID
	} else {
		// First check if it's an existing schema ID
		if _, err := lib.Get(arg); err == nil {
			schemaID = arg
		} else {
			// Not a schema ID - try as file path
			resolvedID, _, err := resolveSchemaSource(arg)
			if err != nil {
				return nil, fmt.Errorf("invalid schema argument '%s': %w", arg, err)
			}
			schemaID = resolvedID
		}
	}

	// Load schema from library
	schema, err := lib.Get(schemaID)
	if err != nil {
		return nil, fmt.Errorf("failed to load schema '%s': %w", schemaID, err)
	}

	return schema, nil
}

func resolveSchemaSource(filePath string) (schemaID string, content []byte, err error) {
	lib := library.NewLibrary()

	// Normalize to absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Load file content
	content, err = loadSchemaFromFile(absPath)
	if err != nil {
		return "", nil, err
	}

	// Calculate file hash
	fileHash := library.CalculateFileHash(content)
	existingSchema, err := lib.FindByPath(absPath)

	// No match - register new schema
	if err != nil {
		schemaID, err := registerSchema(absPath, content)
		return schemaID, content, err
	}

	// Hash matches - use existing schema
	if existingSchema.Metadata.FileHash == fileHash {
		return existingSchema.ID, existingSchema.Content, nil
	}

	// Hash mismatch - handle update workflow
	return handleSchemaUpdate(lib, existingSchema, content)
}

func handleSchemaUpdate(lib library.Library, existingSchema *library.Schema, newContent []byte) (string, []byte, error) {
	fmt.Printf("Schema file has changed since last import.\n")
	update, err := PromptYesNo("Update library")
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user input: %w", err)
	}

	// User chose not to update
	if !update {
		fmt.Printf("Using existing library version\n")
		return existingSchema.ID, existingSchema.Content, nil
	}

	// Update library content
	if err := lib.UpdateContent(existingSchema.ID, newContent); err != nil {
		return "", nil, fmt.Errorf("failed to update library: %w", err)
	}

	fmt.Printf("Library schema '%s' updated\n", existingSchema.ID)
	return existingSchema.ID, newContent, nil
}

func registerSchema(filePath string, content []byte) (string, error) {
	lib := library.NewLibrary()

	// Generate suggested ID from filename
	basename := filepath.Base(filePath)
	ext := filepath.Ext(basename)
	suggested := library.SanitizeSchemaID(filepath.Base(filePath[:len(filePath)-len(ext)]))

	// Prompt for schema ID
	schemaID, err := PromptSchemaID(suggested)
	if err != nil {
		return "", fmt.Errorf("failed to get schema ID: %w", err)
	}

	// Prompt for display name
	displayName, err := PromptString("Enter display name", schemaID)
	if err != nil {
		return "", fmt.Errorf("failed to get display name: %w", err)
	}

	// Add to library
	if err := lib.Add(schemaID, displayName, filePath); err != nil {
		return "", fmt.Errorf("failed to add schema to library: %w", err)
	}

	fmt.Printf("Schema '%s' added to library\n", schemaID)
	return schemaID, nil
}
