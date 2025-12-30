package tui

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/libselect"
)

// setupTestLibrary creates a temporary config directory for testing
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

// createTestSchema creates a temporary schema file for testing
func createTestSchema(t *testing.T, dir string, content string) string {
	t.Helper()

	schemaFile := filepath.Join(dir, "test-schema.graphqls")
	err := os.WriteFile(schemaFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test schema: %v", err)
	}

	return schemaFile
}

func TestModel_TransitionFromLibselectToXplr(t *testing.T) {
	is := is.New(t)

	// Setup test library
	tmpDir, cleanup := setupTestLibrary(t)
	defer cleanup()

	// Create test schema
	schemaContent := `
		type Query {
			hello: String!
			world: String!
		}

		type Post {
			id: ID!
			title: String!
		}
	`
	schemaFile := createTestSchema(t, tmpDir, schemaContent)

	// Add schema to library
	lib := library.NewLibrary()
	err := lib.Add("test-schema", "Test Schema", schemaFile)
	is.NoErr(err)

	// Create model starting in libselect mode
	model, err := newModelWithLibselect()
	is.NoErr(err)
	is.Equal(model.state, libselectView)

	// Initialize the model
	cmd := model.Init()
	is.True(cmd == nil) // libselect.Init() returns nil

	// Simulate selecting a schema by sending SchemaSelectedMsg
	parsedSchema, err := adapters.ParseSchemaString(schemaContent)
	is.NoErr(err)

	msg := libselect.SchemaSelectedMsg{
		SchemaID: "test-schema",
		Schema:   parsedSchema,
		Metadata: library.SchemaMetadata{},
	}

	// Update with the selection message
	updatedModel, _ := model.Update(msg)

	// Verify transition to xplr view
	m, ok := updatedModel.(Model)
	is.True(ok)
	is.Equal(m.state, xplrView)

	// Verify xplr model was initialized with correct data
	is.Equal(m.xplr.SchemaID, "test-schema")
	is.True(m.xplr.HasLibraryData)
}

func TestModel_XplrModeHandlesMessagesCorrectly(t *testing.T) {
	is := is.New(t)

	schemaContent := `
		type Query {
			hello: String!
		}
	`
	parsedSchema, err := adapters.ParseSchemaString(schemaContent)
	is.NoErr(err)

	// Create model starting in xplr mode
	model := newModelWithXplr(parsedSchema)
	is.Equal(model.state, xplrView)

	// Send a window size message (should be handled by xplr)
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := model.Update(msg)

	m, ok := updatedModel.(Model)
	is.True(ok)
	is.Equal(m.state, xplrView) // Should still be in xplr view

	// Verify the message was passed to xplr (xplr stores width/height)
	is.Equal(m.xplr.Width(), 100)
	is.Equal(m.xplr.Height(), 50)
}

func TestModel_LibselectModeHandlesMessagesCorrectly(t *testing.T) {
	is := is.New(t)

	// Setup test library
	_, cleanup := setupTestLibrary(t)
	defer cleanup()

	// Create model starting in libselect mode
	model, err := newModelWithLibselect()
	is.NoErr(err)
	is.Equal(model.state, libselectView)

	// Send a window size message (should be handled by libselect)
	msg := tea.WindowSizeMsg{Width: 80, Height: 40}
	updatedModel, _ := model.Update(msg)

	m, ok := updatedModel.(Model)
	is.True(ok)
	is.Equal(m.state, libselectView) // Should still be in libselect view
}

func TestModel_ForwardsWindowSizeOnTransition(t *testing.T) {
	is := is.New(t)

	// Setup test library
	tmpDir, cleanup := setupTestLibrary(t)
	defer cleanup()

	// Create test schema
	schemaContent := `
		type Query {
			hello: String!
		}
	`
	schemaFile := createTestSchema(t, tmpDir, schemaContent)

	// Add schema to library
	lib := library.NewLibrary()
	err := lib.Add("test-schema", "Test Schema", schemaFile)
	is.NoErr(err)

	// Create model starting in libselect mode
	model, err := newModelWithLibselect()
	is.NoErr(err)

	// Send window size message while in libselect mode
	windowMsg := tea.WindowSizeMsg{Width: 120, Height: 60}
	updatedModel, _ := model.Update(windowMsg)
	model, ok := updatedModel.(Model)
	is.True(ok)

	// Verify top-level model stored the window size
	is.Equal(model.width, 120)
	is.Equal(model.height, 60)

	// Transition to xplr by selecting a schema
	parsedSchema, err := adapters.ParseSchemaString(schemaContent)
	is.NoErr(err)

	schemaMsg := libselect.SchemaSelectedMsg{
		SchemaID: "test-schema",
		Schema:   parsedSchema,
		Metadata: library.SchemaMetadata{},
	}

	updatedModel, _ = model.Update(schemaMsg)
	model, ok = updatedModel.(Model)
	is.True(ok)
	is.Equal(model.state, xplrView)

	// Verify xplr model received the window size
	is.Equal(model.xplr.Width(), 120)
	is.Equal(model.xplr.Height(), 60)
}
