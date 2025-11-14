package library

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Library manages schema storage and metadata.
type Library interface {
	// Add adds a new schema to the library from a file path.
	Add(id string, displayName string, sourcePath string) error

	// Get retrieves a schema by ID.
	Get(id string) (*Schema, error)

	// List returns all schemas in the library.
	List() ([]SchemaInfo, error)

	// Remove removes a schema and its metadata.
	Remove(id string) error

	// UpdateMetadata updates the metadata for a schema.
	UpdateMetadata(id string, metadata SchemaMetadata) error

	// AddFavorite adds a type to the favorites list.
	AddFavorite(id string, typeName string) error

	// RemoveFavorite removes a type from the favorites list.
	RemoveFavorite(id string, typeName string) error

	// SetURLPattern sets a URL pattern for a type.
	SetURLPattern(id string, typePattern string, urlPattern string) error
}

// FileLibrary implements Library using file-based storage.
type FileLibrary struct{}

// NewLibrary creates a new Library instance.
func NewLibrary() Library {
	return &FileLibrary{}
}

// ValidateSchemaID checks if a schema ID is valid and returns an error with suggestions if not.
func ValidateSchemaID(id string) error {
	if id == "" {
		return fmt.Errorf("schema ID cannot be empty")
	}

	// Valid format: lowercase letters, numbers, hyphens
	validPattern := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !validPattern.MatchString(id) {
		sanitized := SanitizeSchemaID(id)
		return fmt.Errorf("invalid schema ID '%s': must contain only lowercase letters, numbers, and hyphens (suggested: '%s')", id, sanitized)
	}

	return nil
}

// SanitizeSchemaID converts a string to a valid schema ID format.
func SanitizeSchemaID(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace invalid characters with hyphens
	re := regexp.MustCompile(`[^a-z0-9-]+`)
	s = re.ReplaceAllString(s, "-")

	// Remove leading/trailing hyphens
	s = strings.Trim(s, "-")

	// Replace multiple consecutive hyphens with single hyphen
	re = regexp.MustCompile(`-+`)
	s = re.ReplaceAllString(s, "-")

	return s
}

// schemaFilePath returns the path to a schema file.
func schemaFilePath(id string) (string, error) {
	schemasDir, err := SchemasDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(schemasDir, id+".graphqls"), nil
}

// loadAllMetadata loads all metadata from the metadata.json file.
func loadAllMetadata() (map[string]SchemaMetadata, error) {
	metadataFile, err := MetadataFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(metadataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]SchemaMetadata), nil
		}
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var metadata map[string]SchemaMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata file: %w", err)
	}

	return metadata, nil
}

// saveAllMetadata saves all metadata to the metadata.json file atomically.
func saveAllMetadata(metadata map[string]SchemaMetadata) error {
	metadataFile, err := MetadataFile()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Atomic write: write to temp file, then rename
	tempFile := metadataFile + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp metadata file: %w", err)
	}

	if err := os.Rename(tempFile, metadataFile); err != nil {
		os.Remove(tempFile) // Clean up temp file on error
		return fmt.Errorf("failed to rename temp metadata file: %w", err)
	}

	return nil
}

// Add implements Library.Add.
func (l *FileLibrary) Add(id string, displayName string, sourcePath string) error {
	if err := ValidateSchemaID(id); err != nil {
		return err
	}

	// Ensure config directory exists
	if err := InitConfigDir(); err != nil {
		return fmt.Errorf("failed to initialize config directory: %w", err)
	}

	// Check if schema already exists
	schemaFile, err := schemaFilePath(id)
	if err != nil {
		return err
	}

	if _, err := os.Stat(schemaFile); err == nil {
		return fmt.Errorf("schema with ID '%s' already exists", id)
	}

	// Read source schema file
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source schema file: %w", err)
	}

	// Write schema file
	if err := os.WriteFile(schemaFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write schema file: %w", err)
	}

	// Load existing metadata
	allMetadata, err := loadAllMetadata()
	if err != nil {
		return err
	}

	// Add new metadata
	now := time.Now()
	allMetadata[id] = SchemaMetadata{
		DisplayName: displayName,
		SourceFile:  sourcePath,
		Favorites:   []string{},
		URLPatterns: make(map[string]string),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Save metadata
	if err := saveAllMetadata(allMetadata); err != nil {
		// Try to clean up schema file on metadata save error
		os.Remove(schemaFile)
		return err
	}

	return nil
}

// Get implements Library.Get.
func (l *FileLibrary) Get(id string) (*Schema, error) {
	schemaFile, err := schemaFilePath(id)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(schemaFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("schema '%s' not found", id)
		}
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Load metadata
	allMetadata, err := loadAllMetadata()
	if err != nil {
		return nil, err
	}

	metadata, exists := allMetadata[id]
	if !exists {
		// Return default metadata if not found
		metadata = SchemaMetadata{
			DisplayName: id,
			SourceFile:  "",
			Favorites:   []string{},
			URLPatterns: make(map[string]string),
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		}
	}

	return &Schema{
		ID:       id,
		Content:  content,
		Metadata: metadata,
	}, nil
}

// List implements Library.List.
func (l *FileLibrary) List() ([]SchemaInfo, error) {
	schemasDir, err := SchemasDir()
	if err != nil {
		return nil, err
	}

	// Check if schemas directory exists
	if _, err := os.Stat(schemasDir); os.IsNotExist(err) {
		return []SchemaInfo{}, nil
	}

	// Read all .graphqls files
	entries, err := os.ReadDir(schemasDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read schemas directory: %w", err)
	}

	// Load metadata
	allMetadata, err := loadAllMetadata()
	if err != nil {
		return nil, err
	}

	var schemas []SchemaInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".graphqls") {
			continue
		}

		id := strings.TrimSuffix(entry.Name(), ".graphqls")
		displayName := id

		if metadata, exists := allMetadata[id]; exists {
			displayName = metadata.DisplayName
		}

		schemas = append(schemas, SchemaInfo{
			ID:          id,
			DisplayName: displayName,
		})
	}

	return schemas, nil
}

// Remove implements Library.Remove.
func (l *FileLibrary) Remove(id string) error {
	schemaFile, err := schemaFilePath(id)
	if err != nil {
		return err
	}

	// Check if schema exists
	if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
		return fmt.Errorf("schema '%s' not found", id)
	}

	// Remove schema file
	if err := os.Remove(schemaFile); err != nil {
		return fmt.Errorf("failed to remove schema file: %w", err)
	}

	// Load and update metadata
	allMetadata, err := loadAllMetadata()
	if err != nil {
		return err
	}

	delete(allMetadata, id)

	if err := saveAllMetadata(allMetadata); err != nil {
		return err
	}

	return nil
}

// UpdateMetadata implements Library.UpdateMetadata.
func (l *FileLibrary) UpdateMetadata(id string, metadata SchemaMetadata) error {
	// Verify schema exists
	schemaFile, err := schemaFilePath(id)
	if err != nil {
		return err
	}

	if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
		return fmt.Errorf("schema '%s' not found", id)
	}

	// Load all metadata
	allMetadata, err := loadAllMetadata()
	if err != nil {
		return err
	}

	// Update timestamp
	metadata.UpdatedAt = time.Now()

	allMetadata[id] = metadata

	return saveAllMetadata(allMetadata)
}

// AddFavorite implements Library.AddFavorite.
func (l *FileLibrary) AddFavorite(id string, typeName string) error {
	schema, err := l.Get(id)
	if err != nil {
		return err
	}

	// Check if already a favorite
	for _, fav := range schema.Metadata.Favorites {
		if fav == typeName {
			return nil // Already a favorite
		}
	}

	schema.Metadata.Favorites = append(schema.Metadata.Favorites, typeName)
	return l.UpdateMetadata(id, schema.Metadata)
}

// RemoveFavorite implements Library.RemoveFavorite.
func (l *FileLibrary) RemoveFavorite(id string, typeName string) error {
	schema, err := l.Get(id)
	if err != nil {
		return err
	}

	var newFavorites []string
	for _, fav := range schema.Metadata.Favorites {
		if fav != typeName {
			newFavorites = append(newFavorites, fav)
		}
	}

	schema.Metadata.Favorites = newFavorites
	return l.UpdateMetadata(id, schema.Metadata)
}

// SetURLPattern implements Library.SetURLPattern.
func (l *FileLibrary) SetURLPattern(id string, typePattern string, urlPattern string) error {
	schema, err := l.Get(id)
	if err != nil {
		return err
	}

	if schema.Metadata.URLPatterns == nil {
		schema.Metadata.URLPatterns = make(map[string]string)
	}

	schema.Metadata.URLPatterns[typePattern] = urlPattern
	return l.UpdateMetadata(id, schema.Metadata)
}
