package libselect_test

import (
	"fmt"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui/libselect"
	"github.com/tonysyu/gqlxp/utils/testx"
	"github.com/tonysyu/gqlxp/utils/testx/assert"
)

// mockLibrary provides a mock library for testing
type mockLibrary struct {
	schemas        []library.SchemaInfo
	getErr         error
	listErr        error
	setDefaultErr  error
	setDefaultID   string
}

func (m *mockLibrary) Add(id, displayName, sourcePath string) error {
	return nil
}

func (m *mockLibrary) AddFromContent(id, displayName string, content []byte, sourceInfo string) error {
	return nil
}

func (m *mockLibrary) Get(id string) (*library.Schema, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return &library.Schema{
		ID:      id,
		Content: []byte(`type Query { hello: String }`),
		Metadata: library.SchemaMetadata{
			DisplayName: "Test Schema",
		},
	}, nil
}

func (m *mockLibrary) List() ([]library.SchemaInfo, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.schemas, nil
}

func (m *mockLibrary) Remove(id string) error {
	return nil
}

func (m *mockLibrary) UpdateMetadata(id string, metadata library.SchemaMetadata) error {
	return nil
}

func (m *mockLibrary) SetURLPattern(id, typePattern, urlPattern string) error {
	return nil
}

func (m *mockLibrary) FindByPath(absolutePath string) (*library.Schema, error) {
	return nil, nil
}

func (m *mockLibrary) UpdateContent(id string, content []byte) error {
	return nil
}

func (m *mockLibrary) GetDefaultSchema() (string, error) {
	return "", nil
}

func (m *mockLibrary) SetDefaultSchema(id string) error {
	m.setDefaultID = id
	return m.setDefaultErr
}

func TestModel_Init(t *testing.T) {
	is := is.New(t)

	lib := &mockLibrary{
		schemas: []library.SchemaInfo{
			{ID: "test", DisplayName: "Test"},
		},
	}

	model, err := libselect.New(lib)
	is.NoErr(err)

	cmd := model.Init()
	is.Equal(cmd, nil) // Init should return nil cmd
}

func TestModel_Update_WindowSize(t *testing.T) {
	is := is.New(t)

	lib := &mockLibrary{
		schemas: []library.SchemaInfo{
			{ID: "test", DisplayName: "Test"},
		},
	}

	model, err := libselect.New(lib)
	is.NoErr(err)

	// Send window size message
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	newModel, cmd := model.Update(msg)
	is.Equal(cmd, nil)             // Window size should not return a cmd
	is.True(newModel.View() != "") // Model should update
}

func TestModel_Update_QuitKeys(t *testing.T) {
	is := is.New(t)

	lib := &mockLibrary{
		schemas: []library.SchemaInfo{
			{ID: "test", DisplayName: "Test"},
		},
	}

	model, err := libselect.New(lib)
	is.NoErr(err)

	quitKeys := []string{"ctrl+c", "ctrl+d"}
	for _, key := range quitKeys {
		t.Run(key, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tea.KeyCtrlC}
			if key == "ctrl+d" {
				msg = tea.KeyMsg{Type: tea.KeyCtrlD}
			}
			_, cmd := model.Update(msg)
			is.True(cmd != nil) // Quit key should return quit cmd
		})
	}
}

func TestModel_View_EmptyLibrary(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	lib := &mockLibrary{
		schemas: []library.SchemaInfo{},
	}

	model, err := libselect.New(lib)
	is.NoErr(err)

	// Set window size to trigger proper rendering
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	model, _ = model.Update(msg)

	assert.StringContains(model.View(), "No schemas in library")
}

func TestModel_View_WithSchemas(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	updatedAt := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	lib := &mockLibrary{
		schemas: []library.SchemaInfo{
			{ID: "schema1", DisplayName: "Schema One", UpdatedAt: updatedAt},
			{ID: "schema2", DisplayName: "Schema Two", UpdatedAt: updatedAt},
		},
	}

	model, err := libselect.New(lib)
	is.NoErr(err)

	// Set window size to trigger proper rendering
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	model, _ = model.Update(msg)

	assert.StringContains(testx.NormalizeView(model.View()), testx.NormalizeView(`
		│ Schema One (id: schema1)
		│ last updated: 2024-03-15 10:30

		  Schema Two (id: schema2)
		  last updated: 2024-03-15 10:30
	`))
}

func TestModel_View_WithSchemas_NoName(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	lib := &mockLibrary{
		schemas: []library.SchemaInfo{
			{ID: "schema1"},
		},
	}

	model, err := libselect.New(lib)
	is.NoErr(err)

	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	model, _ = model.Update(msg)

	assert.StringContains(testx.NormalizeView(model.View()), testx.NormalizeView(`
		│ schema1
		│ last updated: unknown
	`))
}

func TestNew_Handles_Error(t *testing.T) {
	is := is.New(t)

	lib := &mockLibrary{
		listErr: fmt.Errorf("failed to list schemas"),
	}

	_, err := libselect.New(lib)
	is.True(err != nil) // Should return error from library
}

func TestModel_Update_SetDefaultKey(t *testing.T) {
	is := is.New(t)

	lib := &mockLibrary{
		schemas: []library.SchemaInfo{
			{ID: "schema1", DisplayName: "Schema One"},
		},
	}

	model, err := libselect.New(lib)
	is.NoErr(err)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	_, cmd := model.Update(msg)
	is.True(cmd != nil) // Should return a cmd

	// Execute the cmd and verify it calls SetDefaultSchema and returns the message
	result := cmd()
	defaultMsg, ok := result.(libselect.DefaultSchemaSetMsg)
	is.True(ok) // Should return DefaultSchemaSetMsg
	is.Equal(defaultMsg.SchemaID, "schema1")
	is.Equal(lib.setDefaultID, "schema1") // Library should be called with the correct ID
}

func TestModel_Update_DefaultSchemaSetMsg(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	lib := &mockLibrary{
		schemas: []library.SchemaInfo{
			{ID: "schema1", DisplayName: "Schema One"},
			{ID: "schema2", DisplayName: "Schema Two"},
		},
	}

	model, err := libselect.New(lib)
	is.NoErr(err)

	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	model, _ = model.Update(msg)

	// Send DefaultSchemaSetMsg to mark schema2 as default
	model, _ = model.Update(libselect.DefaultSchemaSetMsg{SchemaID: "schema2"})

	assert.StringContains(testx.NormalizeView(model.View()), "Schema Two (Default; id: schema2)")
}
