package xplr

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
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

	schemaView, _ := adapters.ParseSchemaString(schemaString)
	model := New(schemaView)

	// Test initial state
	is.Equal(len(model.nav.Stack().All()), config.VisiblePanelCount)
	is.Equal(model.nav.Stack().Position(), 0)
	is.Equal(navTypeToGQLType(model.nav.CurrentType()), queryType)
	is.Equal(len(model.schema.GetQueryItems()), 2)    // getAllPosts, getPostById
	is.Equal(len(model.schema.GetMutationItems()), 1) // createPost

	// Test that first panel is properly initialized with Query fields
	firstPanel := model.nav.Stack().All()[0]
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

	model := New(adapters.SchemaView{})

	// Test initial stack position
	is.Equal(model.nav.Stack().Position(), 0)

	// Directly add panels to the stack for testing (simulating real navigation)
	// In real usage, panels are added via OpenPanel which truncates and appends
	stack := model.nav.Stack()
	allPanels := []*components.Panel{
		stack.All()[0],
		stack.All()[1],
		components.NewEmptyPanel("test3"),
		components.NewEmptyPanel("test4"),
	}
	stack.Replace(allPanels)
	// Now we have 4 panels total

	// Test next panel navigation (move forward in stack)
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updatedModel.(Model)
	is.Equal(model.nav.Stack().Position(), 1)

	// Test another forward navigation
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updatedModel.(Model)
	is.Equal(model.nav.Stack().Position(), 2)

	// Test another forward navigation
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updatedModel.(Model)
	is.Equal(model.nav.Stack().Position(), 3)

	// Test that we can't go beyond the last panel
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updatedModel.(Model)
	is.Equal(model.nav.Stack().Position(), 3) // Should stay at 3

	// Test previous panel navigation (move backward in stack)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	model = updatedModel.(Model)
	is.Equal(model.nav.Stack().Position(), 2)

	// Test another backward navigation
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	model = updatedModel.(Model)
	is.Equal(model.nav.Stack().Position(), 1)

	// Navigate to beginning
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	model = updatedModel.(Model)
	is.Equal(model.nav.Stack().Position(), 0)

	// Test that we can't go before the beginning
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	model = updatedModel.(Model)
	is.Equal(model.nav.Stack().Position(), 0) // Should stay at 0
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

	schemaView, _ := adapters.ParseSchemaString(schemaString)
	model := New(schemaView)

	// Test initial type
	is.Equal(model.nav.CurrentType(), navigation.QueryType)

	// Test forward cycling through types
	expectedTypes := []navigation.GQLType{
		navigation.MutationType, navigation.ObjectType, navigation.InputType,
		navigation.EnumType, navigation.ScalarType, navigation.InterfaceType,
		navigation.UnionType, navigation.DirectiveType, navigation.QueryType,
	}

	for _, expectedType := range expectedTypes {
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlT})
		model = updatedModel.(Model)
		is.Equal(model.nav.CurrentType(), expectedType)
		is.Equal(model.nav.Stack().Position(), 0) // Stack position should reset to 0
	}

	// Test reverse cycling
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlR})
	model = updatedModel.(Model)
	is.Equal(model.nav.CurrentType(), navigation.DirectiveType)
}

func TestModelWindowResize(t *testing.T) {
	is := is.New(t)

	model := New(adapters.SchemaView{})

	// Test window resize
	newWidth, newHeight := 120, 40
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: newWidth, Height: newHeight})
	model = updatedModel.(Model)

	is.Equal(model.width, newWidth)
	is.Equal(model.height, newHeight)
}

func TestModelWithEmptySchema(t *testing.T) {
	is := is.New(t)

	// Test with completely empty schema
	model := New(adapters.SchemaView{})

	// Model should still initialize properly
	is.Equal(len(model.nav.Stack().All()), config.VisiblePanelCount)
	is.Equal(model.nav.Stack().Position(), 0)
	is.Equal(model.nav.CurrentType(), navigation.QueryType)

	// Should be able to cycle through types even with empty schema
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlT})
	model = updatedModel.(Model)
	is.Equal(model.nav.CurrentType(), navigation.MutationType)
}

func TestModelKeyboardShortcuts(t *testing.T) {
	is := is.New(t)

	model := New(adapters.SchemaView{})

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
