package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tonysyu/gqlxp/library"
)

func loadSchemaFromFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}
	return content, nil
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
	suggested := strings.TrimSuffix(basename, ext)
	suggested = library.SanitizeSchemaID(suggested)

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
