package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/tonysyu/igq/tui/adapters"
	"github.com/tonysyu/igq/tui/components"
)

func TestNewModel(t *testing.T) {
	is := is.New(t)

	// Create a basic schema for testing
	schemaString := `
		type Query {
			getAllPosts: [Post!]!
			getPostById(id: ID!): Post
		}

		type Mutation {
			createPost(title: String!, content: String!): Post!
		}

		type Post {
			id: ID!
			title: String!
		}
	`

	schemaView, _ := adapters.ParseSchema([]byte(schemaString))
	model := newModel(schemaView)

	// Test initial state
	is.Equal(len(model.panelStack), displayedPanels)
	is.Equal(model.stackPosition, 0)
	is.Equal(model.selectedGQLType, queryType)
	is.Equal(len(model.schema.GetQueryItems()), 2)    // getAllPosts, getPostById
	is.Equal(len(model.schema.GetMutationItems()), 1) // createPost

	// Test that first panel is properly initialized with Query fields
	firstPanel := model.panelStack[0]
	// The first panel should be a list panel with Query fields
	is.True(firstPanel != nil)

	// Test keybindings are properly set
	is.True(model.keymap.NextPanel.Enabled())
	is.True(model.keymap.PrevPanel.Enabled())
	is.True(model.keymap.Quit.Enabled())
	is.True(model.keymap.ToggleGQLType.Enabled())
}

func TestModelPanelNavigation(t *testing.T) {
	is := is.New(t)

	model := newModel(adapters.SchemaView{})

	// Test initial stack position
	is.Equal(model.stackPosition, 0)

	// Add some panels to the stack to enable navigation
	model.panelStack = append(model.panelStack, components.NewStringPanel("test3"))
	model.panelStack = append(model.panelStack, components.NewStringPanel("test4"))
	// Now we have 4 panels total

	// Test next panel navigation (move forward in stack)
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updatedModel.(mainModel)
	is.Equal(model.stackPosition, 1)

	// Test another forward navigation
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updatedModel.(mainModel)
	is.Equal(model.stackPosition, 2)

	// Test another forward navigation
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updatedModel.(mainModel)
	is.Equal(model.stackPosition, 3)

	// Test that we can't go beyond the last panel
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updatedModel.(mainModel)
	is.Equal(model.stackPosition, 3) // Should stay at 3

	// Test previous panel navigation (move backward in stack)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	model = updatedModel.(mainModel)
	is.Equal(model.stackPosition, 2)

	// Test another backward navigation
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	model = updatedModel.(mainModel)
	is.Equal(model.stackPosition, 1)

	// Navigate to beginning
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	model = updatedModel.(mainModel)
	is.Equal(model.stackPosition, 0)

	// Test that we can't go before the beginning
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	model = updatedModel.(mainModel)
	is.Equal(model.stackPosition, 0) // Should stay at 0
}

func TestModelGQLTypeSwitching(t *testing.T) {
	is := is.New(t)

	// Create schema with multiple types
	schemaString := `
		type Query {
			getUser: User
		}

		type Mutation {
			createUser: User
		}

		type User {
			id: ID!
			name: String!
		}

		input UserInput {
			name: String!
		}

		enum Status {
			ACTIVE
			INACTIVE
		}

		scalar Date

		interface Node {
			id: ID!
		}

		union SearchResult = User

		directive @deprecated on FIELD_DEFINITION
	`

	schemaView, _ := adapters.ParseSchema([]byte(schemaString))
	model := newModel(schemaView)

	// Test initial type
	is.Equal(model.selectedGQLType, queryType)

	// Test forward cycling through types
	expectedTypes := []gqlType{mutationType, objectType, inputType, enumType, scalarType, interfaceType, unionType, directiveType, queryType}

	for _, expectedType := range expectedTypes {
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlT})
		model = updatedModel.(mainModel)
		is.Equal(model.selectedGQLType, expectedType)
		is.Equal(model.stackPosition, 0) // Stack position should reset to 0
	}

	// Test reverse cycling
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlR})
	model = updatedModel.(mainModel)
	is.Equal(model.selectedGQLType, directiveType)
}

func TestModelWindowResize(t *testing.T) {
	is := is.New(t)

	model := newModel(adapters.SchemaView{})

	// Test window resize
	newWidth, newHeight := 120, 40
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: newWidth, Height: newHeight})
	model = updatedModel.(mainModel)

	is.Equal(model.width, newWidth)
	is.Equal(model.height, newHeight)
}

func TestModelWithEmptySchema(t *testing.T) {
	is := is.New(t)

	// Test with completely empty schema
	model := newModel(adapters.SchemaView{})

	// Model should still initialize properly
	is.Equal(len(model.panelStack), displayedPanels)
	is.Equal(model.stackPosition, 0)
	is.Equal(model.selectedGQLType, queryType)

	// Should be able to cycle through types even with empty schema
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlT})
	model = updatedModel.(mainModel)
	is.Equal(model.selectedGQLType, mutationType)
}

func TestModelKeyboardShortcuts(t *testing.T) {
	is := is.New(t)

	model := newModel(adapters.SchemaView{})

	// Test all keyboard shortcuts don't crash
	shortcuts := []tea.KeyMsg{
		{Type: tea.KeyTab},
		{Type: tea.KeyShiftTab},
		{Type: tea.KeyCtrlT},
		{Type: tea.KeyCtrlR},
		{Type: tea.KeyCtrlC},
		{Type: tea.KeyCtrlD},
	}

	for _, shortcut := range shortcuts {
		_, cmd := model.Update(shortcut)
		// Quit commands should return a quit command
		if shortcut.Type == tea.KeyCtrlC || shortcut.Type == tea.KeyCtrlD {
			is.True(cmd != nil) // Should return tea.Quit command
		}
	}
}
