package xplr

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
	"github.com/tonysyu/gqlxp/tui/xplr/searchmodel"
)

// Key messages for simulating user input
var (
	keyNextPanel     = tea.KeyPressMsg{Code: ']'}
	keyPrevPanel     = tea.KeyPressMsg{Code: '['}
	keyNextType      = tea.KeyPressMsg{Code: '}'}
	keyPrevType      = tea.KeyPressMsg{Code: '{'}
	keyNextItem      = tea.KeyPressMsg{Code: tea.KeyDown}
	keyPrevItem      = tea.KeyPressMsg{Code: tea.KeyUp}
	keyToggleOverlay = tea.KeyPressMsg{Code: ' '}
	keyOpenLibSelect = tea.KeyPressMsg{Code: 'o', Mod: tea.ModCtrl}
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
	is.Equal(model.nav.CurrentType(), navigation.QueryType)
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
	is.True(model.keymap.NextGQLType.Enabled())
}

func TestModelPanelNavigation(t *testing.T) {
	is := is.New(t)

	model := New(adapters.SchemaView{})

	// Test initial stack position
	is.Equal(model.PanelPosition(), 0)

	// Build a 4-panel stack for navigation testing.
	// OpenPanel truncates-and-appends after current position, so navigate
	// forward between each push to accumulate without truncation.
	model.nav = model.nav.OpenPanel(components.NewEmptyPanel("test3")) // [p0, test3] at 0
	model.nav, _ = model.nav.NavigateForward()                         // pos 1
	model.nav = model.nav.OpenPanel(components.NewEmptyPanel("test4")) // [p0, test3, test4] at 1
	model.nav, _ = model.nav.NavigateForward()                         // pos 2
	model.nav = model.nav.OpenPanel(components.NewEmptyPanel("test5")) // [p0, test3, test4, test5] at 2
	model.nav, _ = model.nav.NavigateBackward()                        // pos 1
	model.nav, _ = model.nav.NavigateBackward()                        // pos 0
	// Now we have 4 panels at position 0

	// Test next panel navigation (move forward in stack)
	model, _ = model.Update(keyNextPanel)
	is.Equal(model.PanelPosition(), 1)

	// Test another forward navigation
	model, _ = model.Update(keyNextPanel)
	is.Equal(model.PanelPosition(), 2)

	// Test another forward navigation
	model, _ = model.Update(keyNextPanel)
	is.Equal(model.PanelPosition(), 3)

	// Test that we can't go beyond the last panel
	model, _ = model.Update(keyNextPanel)
	is.Equal(model.PanelPosition(), 3) // Should stay at 3

	// Test previous panel navigation (move backward in stack)
	model, _ = model.Update(keyPrevPanel)
	is.Equal(model.PanelPosition(), 2)

	// Test another backward navigation
	model, _ = model.Update(keyPrevPanel)
	is.Equal(model.PanelPosition(), 1)

	// Navigate to beginning
	model, _ = model.Update(keyPrevPanel)
	is.Equal(model.PanelPosition(), 0)

	// Test that we can't go before the beginning
	model, _ = model.Update(keyPrevPanel)
	is.Equal(model.PanelPosition(), 0) // Should stay at 0
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
		navigation.UnionType, navigation.DirectiveType, navigation.SearchType, navigation.QueryType,
	}

	for _, expectedType := range expectedTypes {
		model, _ = model.Update(keyNextType)
		is.Equal(model.nav.CurrentType(), expectedType)
		is.Equal(model.nav.Stack().Position(), 0) // Stack position should reset to 0
	}

	// Test reverse cycling
	model, _ = model.Update(keyPrevType)
	is.Equal(model.nav.CurrentType(), navigation.SearchType)
}

func TestModelWindowResize(t *testing.T) {
	is := is.New(t)

	model := New(adapters.SchemaView{})

	// Test window resize
	newWidth, newHeight := 120, 40
	model, _ = model.Update(tea.WindowSizeMsg{Width: newWidth, Height: newHeight})

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
	model, _ = model.Update(tea.KeyPressMsg{Code: '}'})
	is.Equal(model.nav.CurrentType(), navigation.MutationType)
}

func TestSearchResultsReadyClearsChildPanel(t *testing.T) {
	is := is.New(t)

	model := New(adapters.SchemaView{})
	model.nav = model.nav.SwitchType(navigation.SearchType)

	// Simulate a child panel opened from a previous search result
	model.nav = model.nav.OpenPanel(components.NewEmptyPanel("previous result"))
	is.Equal(model.nav.Stack().Next().Title(), "previous result")

	// New search returns no results
	model, _ = model.Update(searchmodel.ResultsReadyMsg{Items: nil})

	// Child panel should be cleared
	is.Equal(model.nav.Stack().Next().Title(), "")
}

func TestSearchResultsReadyClearsChildPanelWhenResultsExist(t *testing.T) {
	is := is.New(t)

	model := New(adapters.SchemaView{})
	model.nav = model.nav.SwitchType(navigation.SearchType)

	// Simulate a child panel opened from a previous search result
	model.nav = model.nav.OpenPanel(components.NewEmptyPanel("previous result"))
	is.Equal(model.nav.Stack().Next().Title(), "previous result")

	// New search returns results
	items := []components.ListItem{components.NewSimpleItem("new result")}
	model, _ = model.Update(searchmodel.ResultsReadyMsg{Items: items})

	// Child panel from the old search should still be cleared
	is.Equal(model.nav.Stack().Next().Title(), "")
}

func TestModelOpenLibSelect(t *testing.T) {
	is := is.New(t)

	model := New(adapters.SchemaView{})

	_, cmd := model.Update(keyOpenLibSelect)

	is.True(cmd != nil)
	msg := cmd()
	_, ok := msg.(OpenLibSelectMsg)
	is.True(ok)
}

func TestModelKeyboardShortcuts(t *testing.T) {
	schemaString := `
		type Query { getUser: User }
		type Mutation { createUser: User }
		type User { id: ID! }
	`
	schemaView, _ := adapters.ParseSchemaString(schemaString)

	tests := []struct {
		name   string
		key    tea.KeyPressMsg
		setup  func(m Model) Model
		verify func(t *testing.T, m Model)
	}{
		{
			name: "next panel moves stack forward",
			key:  keyNextPanel,
			setup: func(m Model) Model {
				m.nav = m.nav.OpenPanel(components.NewEmptyPanel("test"))
				return m
			},
			verify: func(t *testing.T, m Model) {
				is.New(t).Equal(m.PanelPosition(), 1)
			},
		},
		{
			name: "prev panel has no effect at position 0",
			key:  keyPrevPanel,
			verify: func(t *testing.T, m Model) {
				is.New(t).Equal(m.PanelPosition(), 0)
			},
		},
		{
			name: "next GQL type cycles forward",
			key:  keyNextType,
			verify: func(t *testing.T, m Model) {
				is.New(t).Equal(m.nav.CurrentType(), navigation.MutationType)
			},
		},
		{
			name: "prev GQL type wraps to last type",
			key:  keyPrevType,
			verify: func(t *testing.T, m Model) {
				is.New(t).Equal(m.nav.CurrentType(), navigation.SearchType)
			},
		},
		{
			name: "down key does not crash",
			key:  keyNextItem,
			verify: func(t *testing.T, m Model) {
				is.New(t).True(m.View() != "")
			},
		},
		{
			name: "up key does not crash",
			key:  keyPrevItem,
			verify: func(t *testing.T, m Model) {
				is.New(t).True(m.View() != "")
			},
		},
		{
			name: "space opens overlay",
			key:  keyToggleOverlay,
			verify: func(t *testing.T, m Model) {
				is.New(t).True(m.IsOverlayVisible())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New(schemaView)
			model, _ = model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			if tt.setup != nil {
				model = tt.setup(model)
			}
			model, _ = model.Update(tt.key)
			tt.verify(t, model)
		})
	}
}
