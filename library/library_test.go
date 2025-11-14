package library_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/library"
)

// setupTestLibrary creates a temporary config directory for testing.
func setupTestLibrary(t *testing.T) (string, func()) {
	t.Helper()

	// Create temporary directory
	tmpDir := t.TempDir()

	// Set environment variable to override config directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	cleanup := func() {
		os.Setenv("HOME", oldHome)
	}

	return tmpDir, cleanup
}

// createTestSchema creates a temporary schema file for testing.
func createTestSchema(t *testing.T, dir string, content string) string {
	t.Helper()

	schemaFile := filepath.Join(dir, "test-schema.graphqls")
	err := os.WriteFile(schemaFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test schema: %v", err)
	}

	return schemaFile
}

func TestValidateSchemaID(t *testing.T) {
	is := is.New(t)

	t.Run("valid schema ID with lowercase letters", func(t *testing.T) {
		err := library.ValidateSchemaID("github-api")
		is.NoErr(err)
	})

	t.Run("valid schema ID with numbers", func(t *testing.T) {
		err := library.ValidateSchemaID("api-v2")
		is.NoErr(err)
	})

	t.Run("invalid schema ID with uppercase letters", func(t *testing.T) {
		err := library.ValidateSchemaID("GitHub-API")
		is.True(err != nil)
	})

	t.Run("invalid schema ID with spaces", func(t *testing.T) {
		err := library.ValidateSchemaID("github api")
		is.True(err != nil)
	})

	t.Run("invalid schema ID with special characters", func(t *testing.T) {
		err := library.ValidateSchemaID("github_api")
		is.True(err != nil)
	})

	t.Run("empty schema ID", func(t *testing.T) {
		err := library.ValidateSchemaID("")
		is.True(err != nil)
	})
}

func TestSanitizeSchemaID(t *testing.T) {
	is := is.New(t)

	t.Run("converts uppercase to lowercase", func(t *testing.T) {
		result := library.SanitizeSchemaID("GitHub-API")
		is.Equal(result, "github-api")
	})

	t.Run("replaces spaces with hyphens", func(t *testing.T) {
		result := library.SanitizeSchemaID("github api")
		is.Equal(result, "github-api")
	})

	t.Run("replaces underscores with hyphens", func(t *testing.T) {
		result := library.SanitizeSchemaID("github_api")
		is.Equal(result, "github-api")
	})

	t.Run("removes leading and trailing hyphens", func(t *testing.T) {
		result := library.SanitizeSchemaID("-github-api-")
		is.Equal(result, "github-api")
	})

	t.Run("collapses multiple hyphens", func(t *testing.T) {
		result := library.SanitizeSchemaID("github---api")
		is.Equal(result, "github-api")
	})

	t.Run("handles complex input", func(t *testing.T) {
		result := library.SanitizeSchemaID("  GitHub's API v2.0! ")
		is.Equal(result, "github-s-api-v2-0")
	})
}

func TestLibrary_Add(t *testing.T) {
	is := is.New(t)
	tmpDir, cleanup := setupTestLibrary(t)
	defer cleanup()

	lib := library.NewLibrary()
	schemaContent := `type Query { hello: String }`
	schemaFile := createTestSchema(t, tmpDir, schemaContent)

	t.Run("add new schema successfully", func(t *testing.T) {
		err := lib.Add("test-schema", "Test Schema", schemaFile)
		is.NoErr(err)

		// Verify schema was added
		schema, err := lib.Get("test-schema")
		is.NoErr(err)
		is.Equal(schema.ID, "test-schema")
		is.Equal(schema.Metadata.DisplayName, "Test Schema")
		is.Equal(string(schema.Content), schemaContent)
	})

	t.Run("add schema with duplicate ID fails", func(t *testing.T) {
		err := lib.Add("test-schema", "Another Schema", schemaFile)
		is.True(err != nil)
	})

	t.Run("add schema with invalid ID fails", func(t *testing.T) {
		err := lib.Add("Invalid_ID", "Invalid Schema", schemaFile)
		is.True(err != nil)
	})

	t.Run("add schema with non-existent file fails", func(t *testing.T) {
		err := lib.Add("missing-schema", "Missing Schema", "/nonexistent/file.graphqls")
		is.True(err != nil)
	})
}

func TestLibrary_Get(t *testing.T) {
	is := is.New(t)
	tmpDir, cleanup := setupTestLibrary(t)
	defer cleanup()

	lib := library.NewLibrary()
	schemaContent := `type Query { hello: String }`
	schemaFile := createTestSchema(t, tmpDir, schemaContent)

	err := lib.Add("test-schema", "Test Schema", schemaFile)
	is.NoErr(err)

	t.Run("get existing schema", func(t *testing.T) {
		schema, err := lib.Get("test-schema")
		is.NoErr(err)
		is.Equal(schema.ID, "test-schema")
		is.Equal(schema.Metadata.DisplayName, "Test Schema")
		is.Equal(schema.Metadata.SourceFile, schemaFile)
		is.Equal(string(schema.Content), schemaContent)
	})

	t.Run("get non-existent schema fails", func(t *testing.T) {
		_, err := lib.Get("nonexistent-schema")
		is.True(err != nil)
	})
}

func TestLibrary_List(t *testing.T) {
	is := is.New(t)
	tmpDir, cleanup := setupTestLibrary(t)
	defer cleanup()

	lib := library.NewLibrary()
	schemaContent := `type Query { hello: String }`
	schemaFile := createTestSchema(t, tmpDir, schemaContent)

	t.Run("list empty library", func(t *testing.T) {
		schemas, err := lib.List()
		is.NoErr(err)
		is.Equal(len(schemas), 0)
	})

	t.Run("list library with schemas", func(t *testing.T) {
		err := lib.Add("schema-1", "Schema One", schemaFile)
		is.NoErr(err)

		err = lib.Add("schema-2", "Schema Two", schemaFile)
		is.NoErr(err)

		schemas, err := lib.List()
		is.NoErr(err)
		is.Equal(len(schemas), 2)

		// Verify schema info
		schemaMap := make(map[string]string)
		for _, s := range schemas {
			schemaMap[s.ID] = s.DisplayName
		}

		is.Equal(schemaMap["schema-1"], "Schema One")
		is.Equal(schemaMap["schema-2"], "Schema Two")
	})
}

func TestLibrary_Remove(t *testing.T) {
	is := is.New(t)
	tmpDir, cleanup := setupTestLibrary(t)
	defer cleanup()

	lib := library.NewLibrary()
	schemaContent := `type Query { hello: String }`
	schemaFile := createTestSchema(t, tmpDir, schemaContent)

	err := lib.Add("test-schema", "Test Schema", schemaFile)
	is.NoErr(err)

	t.Run("remove existing schema", func(t *testing.T) {
		err := lib.Remove("test-schema")
		is.NoErr(err)

		// Verify schema was removed
		_, err = lib.Get("test-schema")
		is.True(err != nil)

		// Verify library is empty
		schemas, err := lib.List()
		is.NoErr(err)
		is.Equal(len(schemas), 0)
	})

	t.Run("remove non-existent schema fails", func(t *testing.T) {
		err := lib.Remove("nonexistent-schema")
		is.True(err != nil)
	})
}

func TestLibrary_UpdateMetadata(t *testing.T) {
	is := is.New(t)
	tmpDir, cleanup := setupTestLibrary(t)
	defer cleanup()

	lib := library.NewLibrary()
	schemaContent := `type Query { hello: String }`
	schemaFile := createTestSchema(t, tmpDir, schemaContent)

	err := lib.Add("test-schema", "Test Schema", schemaFile)
	is.NoErr(err)

	t.Run("update metadata successfully", func(t *testing.T) {
		schema, err := lib.Get("test-schema")
		is.NoErr(err)

		// Update display name
		schema.Metadata.DisplayName = "Updated Schema Name"

		err = lib.UpdateMetadata("test-schema", schema.Metadata)
		is.NoErr(err)

		// Verify update
		updatedSchema, err := lib.Get("test-schema")
		is.NoErr(err)
		is.Equal(updatedSchema.Metadata.DisplayName, "Updated Schema Name")
	})

	t.Run("update metadata for non-existent schema fails", func(t *testing.T) {
		metadata := library.SchemaMetadata{DisplayName: "Test"}
		err := lib.UpdateMetadata("nonexistent-schema", metadata)
		is.True(err != nil)
	})
}

func TestLibrary_AddFavorite(t *testing.T) {
	is := is.New(t)
	tmpDir, cleanup := setupTestLibrary(t)
	defer cleanup()

	lib := library.NewLibrary()
	schemaContent := `type Query { hello: String }`
	schemaFile := createTestSchema(t, tmpDir, schemaContent)

	err := lib.Add("test-schema", "Test Schema", schemaFile)
	is.NoErr(err)

	t.Run("add favorite type", func(t *testing.T) {
		err := lib.AddFavorite("test-schema", "Query")
		is.NoErr(err)

		schema, err := lib.Get("test-schema")
		is.NoErr(err)
		is.Equal(len(schema.Metadata.Favorites), 1)
		is.Equal(schema.Metadata.Favorites[0], "Query")
	})

	t.Run("add multiple favorites", func(t *testing.T) {
		err := lib.AddFavorite("test-schema", "User")
		is.NoErr(err)

		err = lib.AddFavorite("test-schema", "Post")
		is.NoErr(err)

		schema, err := lib.Get("test-schema")
		is.NoErr(err)
		is.Equal(len(schema.Metadata.Favorites), 3)
	})

	t.Run("add duplicate favorite is idempotent", func(t *testing.T) {
		err := lib.AddFavorite("test-schema", "Query")
		is.NoErr(err)

		schema, err := lib.Get("test-schema")
		is.NoErr(err)
		is.Equal(len(schema.Metadata.Favorites), 3) // Still 3, not 4
	})

	t.Run("add favorite to non-existent schema fails", func(t *testing.T) {
		err := lib.AddFavorite("nonexistent-schema", "Query")
		is.True(err != nil)
	})
}

func TestLibrary_RemoveFavorite(t *testing.T) {
	is := is.New(t)
	tmpDir, cleanup := setupTestLibrary(t)
	defer cleanup()

	lib := library.NewLibrary()
	schemaContent := `type Query { hello: String }`
	schemaFile := createTestSchema(t, tmpDir, schemaContent)

	err := lib.Add("test-schema", "Test Schema", schemaFile)
	is.NoErr(err)

	// Add favorites first
	err = lib.AddFavorite("test-schema", "Query")
	is.NoErr(err)
	err = lib.AddFavorite("test-schema", "User")
	is.NoErr(err)

	t.Run("remove existing favorite", func(t *testing.T) {
		err := lib.RemoveFavorite("test-schema", "Query")
		is.NoErr(err)

		schema, err := lib.Get("test-schema")
		is.NoErr(err)
		is.Equal(len(schema.Metadata.Favorites), 1)
		is.Equal(schema.Metadata.Favorites[0], "User")
	})

	t.Run("remove non-existent favorite is safe", func(t *testing.T) {
		err := lib.RemoveFavorite("test-schema", "NonExistent")
		is.NoErr(err)

		schema, err := lib.Get("test-schema")
		is.NoErr(err)
		is.Equal(len(schema.Metadata.Favorites), 1)
	})

	t.Run("remove favorite from non-existent schema fails", func(t *testing.T) {
		err := lib.RemoveFavorite("nonexistent-schema", "Query")
		is.True(err != nil)
	})
}

func TestLibrary_SetURLPattern(t *testing.T) {
	is := is.New(t)
	tmpDir, cleanup := setupTestLibrary(t)
	defer cleanup()

	lib := library.NewLibrary()
	schemaContent := `type Query { hello: String }`
	schemaFile := createTestSchema(t, tmpDir, schemaContent)

	err := lib.Add("test-schema", "Test Schema", schemaFile)
	is.NoErr(err)

	t.Run("set URL pattern for specific type", func(t *testing.T) {
		err := lib.SetURLPattern("test-schema", "Query", "https://example.com/docs/${field}")
		is.NoErr(err)

		schema, err := lib.Get("test-schema")
		is.NoErr(err)
		is.Equal(schema.Metadata.URLPatterns["Query"], "https://example.com/docs/${field}")
	})

	t.Run("set wildcard URL pattern", func(t *testing.T) {
		err := lib.SetURLPattern("test-schema", "*", "https://example.com/docs/${type}")
		is.NoErr(err)

		schema, err := lib.Get("test-schema")
		is.NoErr(err)
		is.Equal(schema.Metadata.URLPatterns["*"], "https://example.com/docs/${type}")
	})

	t.Run("update existing URL pattern", func(t *testing.T) {
		err := lib.SetURLPattern("test-schema", "Query", "https://example.com/new-docs/${field}")
		is.NoErr(err)

		schema, err := lib.Get("test-schema")
		is.NoErr(err)
		is.Equal(schema.Metadata.URLPatterns["Query"], "https://example.com/new-docs/${field}")
	})

	t.Run("set URL pattern for non-existent schema fails", func(t *testing.T) {
		err := lib.SetURLPattern("nonexistent-schema", "Query", "https://example.com")
		is.True(err != nil)
	})
}

func TestConfigDir(t *testing.T) {
	is := is.New(t)
	_, cleanup := setupTestLibrary(t)
	defer cleanup()

	t.Run("returns non-empty config directory", func(t *testing.T) {
		configDir, err := library.ConfigDir()
		is.NoErr(err)
		is.True(configDir != "")
	})
}

func TestInitConfigDir(t *testing.T) {
	is := is.New(t)
	_, cleanup := setupTestLibrary(t)
	defer cleanup()

	t.Run("creates config directory structure", func(t *testing.T) {
		err := library.InitConfigDir()
		is.NoErr(err)

		// Verify schemas directory exists
		schemasDir, err := library.SchemasDir()
		is.NoErr(err)

		info, err := os.Stat(schemasDir)
		is.NoErr(err)
		is.True(info.IsDir())

		// Verify metadata file exists
		metadataFile, err := library.MetadataFile()
		is.NoErr(err)

		_, err = os.Stat(metadataFile)
		is.NoErr(err)
	})

	t.Run("calling InitConfigDir multiple times is safe", func(t *testing.T) {
		err := library.InitConfigDir()
		is.NoErr(err)

		err = library.InitConfigDir()
		is.NoErr(err)
	})
}
