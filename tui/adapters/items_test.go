package adapters

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/utils/testx"
	"github.com/tonysyu/gqlxp/utils/testx/assert"
)

// renderMinimalPanel drastically simplifies Panel rendering to create simpler tests
//
// In particular, this will create a panel with the following characteristics:
// 1. Empty lines removed (using NormalizeView)
// 2. Item only show title (no description)
// 3. No selection indicator
// 4. No "status bar" (item count)
// 4. No help
func renderMinimalPanel(panel *components.Panel) string {
	panel.ListModel.SetDelegate(minimalItemDelegate{})
	panel.ListModel.SetShowStatusBar(false)
	panel.ListModel.ShowHelp()
	content := testx.NormalizeView(panel.View())
	return content
}

func nextPanelTab(panel *components.Panel) *components.Panel {
	updatedModel, _ := panel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'L'}})
	return updatedModel.(*components.Panel)
}

type minimalItemDelegate struct{}

func (d minimalItemDelegate) Height() int                             { return 1 }
func (d minimalItemDelegate) Spacing() int                            { return 0 }
func (d minimalItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d minimalItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(list.DefaultItem)
	if !ok {
		return
	}
	fmt.Fprint(w, i.Title())
}

func TestQueryAndMutationItemOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		type Post {
		  id: ID!
		  title: String!
		  content: String!
		}

		type Query {
		  """
		  Return all posts
		  """
		  getAllPosts: [Post!]!
		  getPostById(id: ID!): Post
		}

		type Mutation {
		  """
		  Create a new post
		  """
		  createPost(title: String!, content: String!, authorId: ID!): Post!
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	t.Run("Query field with no arguments shows description and result type", func(t *testing.T) {
		field := schema.Query["getAllPosts"]
		item := newFieldItem(field, resolver)
		panel, _ := item.OpenPanel()

		// Set a reasonable size for testing
		panel.SetSize(80, 40)

		content := testx.NormalizeView(panel.View())

		expected := testx.NormalizeView(`
			  getAllPosts
			  Return all posts
			  Type
			  [Post!]!
		`)

		assert.StringContains(content, expected)
		is.True(!strings.Contains(content, "Inputs")) // Should not have arguments section
	})

	t.Run("Query field with arguments shows all sections", func(t *testing.T) {
		field := schema.Query["getPostById"]
		item := newFieldItem(field, resolver)
		panel, _ := item.OpenPanel()

		// Set a reasonable size for testing
		panel.SetSize(80, 40)

		// Verify Result Type tab (default active tab)
		content := renderMinimalPanel(panel)
		assert.StringContains(content, testx.NormalizeView(`
			Type    Inputs
			Post
		`))

		panel = nextPanelTab(panel) // Switch to Input Arguments tab
		content = renderMinimalPanel(panel)
		assert.StringContains(content, testx.NormalizeView(`
			Type    Inputs
			id: ID!
		`))
	})

	t.Run("Mutation field with multiple arguments shows all sections", func(t *testing.T) {
		field := schema.Mutation["createPost"]
		item := newFieldItem(field, resolver)
		panel, _ := item.OpenPanel()

		// Set a reasonable size for testing
		panel.SetSize(80, 40)

		// Verify Result Type tab (default active tab)
		content := renderMinimalPanel(panel)
		assert.StringContains(content, testx.NormalizeView(`
			createPost
			Create a new post

			Type    Inputs
			Post!
		`))

		panel = nextPanelTab(panel) // Switch to Input Arguments tab
		content = renderMinimalPanel(panel)
		assert.StringContains(content, testx.NormalizeView(`
			createPost
			Create a new post

			Type    Inputs
			title: String!
			content: String!
			authorId: ID!
		`))
	})
}

func TestObjectDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		type User {
		  id: ID!
		  name: String!
		  email: String
		  posts: [Post!]!
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	userObj := schema.Object["User"]
	item := newTypeDefItem(userObj, resolver)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	content := renderMinimalPanel(panel)

	assert.StringContains(content, testx.NormalizeView(`
		User
		Fields
		id: ID!
		name: String!
		email: String
		posts: [Post!]!
	`))
}

func TestObjectWithInterfacesOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		interface Node {
		  id: ID!
		}

		interface Named {
		  name: String!
		}

		type User implements Node & Named {
		  id: ID!
		  name: String!
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	userObj := schema.Object["User"]
	item := newTypeDefItem(userObj, resolver)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	// First tab (Fields) should be displayed by default
	content := renderMinimalPanel(panel)
	assert.StringContains(content, testx.NormalizeView(`
		User
		Fields    Interfaces
		id: ID!
		name: String!
	`))

	// Navigate to Interfaces tab
	panel = nextPanelTab(panel)
	content = renderMinimalPanel(panel)
	assert.StringContains(content, testx.NormalizeView(`
		User
		Fields    Interfaces
		Node
		Named
	`))
}

func TestNavigateFromObjectToInterface(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		interface Node {
		  id: ID!
		}

		type User implements Node {
		  id: ID!
		  name: String!
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	userObj := schema.Object["User"]
	item := newTypeDefItem(userObj, resolver)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	// Navigate to Interfaces tab
	panel = nextPanelTab(panel)

	// Get the first interface item from the Interfaces tab
	items := panel.ListModel.Items()
	is.True(len(items) > 0)

	interfaceItem, ok := items[0].(components.ListItem)
	is.True(ok)

	// Open panel for the interface
	interfacePanel, ok := interfaceItem.OpenPanel()
	is.True(ok)

	interfacePanel.SetSize(80, 40)
	content := renderMinimalPanel(interfacePanel)

	// Verify the interface panel shows its fields and usages
	assert.StringContains(content, testx.NormalizeView(`
		Node
		Fields    Usages
		id: ID!
	`))
}

func TestInputDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		input CreateUserInput {
		  name: String!
		  email: String!
		  age: Int = 18
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	inputObj := schema.Input["CreateUserInput"]
	item := newTypeDefItem(inputObj, resolver)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	content := renderMinimalPanel(panel)

	assert.StringContains(content, testx.NormalizeView(`
		CreateUserInput

		name: String!
		email: String!
		age: Int = 18
	`))
}

func TestEnumDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		enum Status {
		  ACTIVE
		  INACTIVE
		  PENDING
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	enumObj := schema.Enum["Status"]
	item := newTypeDefItem(enumObj, resolver)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	content := renderMinimalPanel(panel)
	assert.StringContains(content, testx.NormalizeView(`
		ACTIVE
		INACTIVE
		PENDING
	`))
}

func TestScalarDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)

	schemaString := "scalar Date"

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	scalarObj := schema.Scalar["Date"]
	item := newTypeDefItem(scalarObj, resolver)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	// Scalar types should have minimal content (just the name)
	content := panel.View()
	is.True(len(content) > 0)
}

func TestInterfaceDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		interface Node {
		  id: ID!
		  createdAt: String
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	interfaceObj := schema.Interface["Node"]
	item := newTypeDefItem(interfaceObj, resolver)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	content := renderMinimalPanel(panel)

	assert.StringContains(content, testx.NormalizeView(`
		Node
		Fields
		id: ID!
		createdAt: String
	`))
}

func TestInterfaceWithInterfacesOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		interface Node {
		  id: ID!
		}

		interface Timestamped {
		  createdAt: String
		  updatedAt: String
		}

		interface Resource implements Node & Timestamped {
		  id: ID!
		  createdAt: String
		  updatedAt: String
		  name: String!
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	interfaceObj := schema.Interface["Resource"]
	item := newTypeDefItem(interfaceObj, resolver)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	// First tab (Fields) should be displayed by default
	content := renderMinimalPanel(panel)
	assert.StringContains(content, testx.NormalizeView(`
		Resource
		Fields    Interfaces
		id: ID!
		createdAt: String
		updatedAt: String
		name: String!
	`))

	// Navigate to Interfaces tab
	panel = nextPanelTab(panel)
	content = renderMinimalPanel(panel)
	assert.StringContains(content, testx.NormalizeView(`
		Resource
		Fields    Interfaces
		Node
		Timestamped
	`))
}

func TestUnionDefinitionItemOpenPanel(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		type User {
		  id: ID!
		}

		type Post {
		  id: ID!
		}

		union SearchResult = User | Post
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	unionObj := schema.Union["SearchResult"]
	item := newTypeDefItem(unionObj, resolver)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	content := renderMinimalPanel(panel)

	assert.StringContains(content, testx.NormalizeView(`
		SearchResult
		User
		Post
	`))
}

func TestFieldDefinitionWithoutDescription(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		type Query {
		  simpleField: String
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	field := schema.Query["simpleField"]
	item := newFieldItem(field, resolver)

	is.Equal(item.Title(), "simpleField: String")
	is.Equal(item.Description(), "") // No description
	is.Equal(item.FilterValue(), "simpleField")

	panel, ok := item.OpenPanel()
	is.True(ok)
	panel.SetSize(80, 40)

	content := testx.NormalizeView(panel.View())

	expected := testx.NormalizeView(`
		Type
		String
	`)

	assert.StringContains(content, expected)
	is.True(!strings.Contains(content, "Inputs"))
}

func TestFieldDefinitionWithComplexArguments(t *testing.T) {
	is := is.New(t)
	assert := assert.New(t)

	schemaString := `
		input FilterInput {
		  search: String
		  limit: Int
		}

		type Query {
		  complexField(
		    id: ID!
		    filters: FilterInput
		    tags: [String!]!
		    metadata: [String]
		  ): [String!]!
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)

	field := schema.Query["complexField"]
	item := newFieldItem(field, resolver)
	panel, ok := item.OpenPanel()

	is.True(ok)
	panel.SetSize(80, 40)

	// Verify Result Type tab (default active tab)
	content := renderMinimalPanel(panel)
	assert.StringContains(content, testx.NormalizeView(`
		Type    Inputs
		[String!]!
	`))

	panel = nextPanelTab(panel) // Switch to Input Arguments tab
	content = renderMinimalPanel(panel)
	assert.StringContains(content, testx.NormalizeView(`
		Type    Inputs
		id: ID!
		filters: FilterInput
		tags: [String!]!
		metadata: [String]
	`))
}

func TestSimpleItemInterface(t *testing.T) {
	is := is.New(t)

	item := components.NewSimpleItem("Test Title", components.WithDescription("Test Description"))

	is.Equal(item.Title(), "Test Title")
	is.Equal(item.Description(), "Test Description")
	is.Equal(item.FilterValue(), "Test Title")

	// Simple items should not be openable
	panel, ok := item.OpenPanel()
	is.True(!ok)
	is.Equal(panel, nil)
}

func TestArgumentListCreation(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Query {
		  testField(arg1: String!, arg2: Int, arg3: [String]): String
		}
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)
	field := schema.Query["testField"]

	// Test argument items creation
	items := adaptArguments(field.Arguments(), resolver)
	is.Equal(len(items), 3)

	// Test first argument
	item1 := items[0]
	is.Equal(item1.Title(), "arg1: String!")

	// Test second argument
	item2 := items[1]
	is.Equal(item2.Title(), "arg2: Int")

	// Test third argument
	item3 := items[2]
	is.Equal(item3.Title(), "arg3: [String]")
}

func TestDirectiveDefinitionItemCreation(t *testing.T) {
	is := is.New(t)

	schemaString := `
		directive @deprecated(reason: String = "No longer supported") on FIELD_DEFINITION | ENUM_VALUE
	`

	schema, _ := gql.ParseSchema([]byte(schemaString))
	resolver := gql.NewSchemaResolver(&schema)
	directive := schema.Directive["deprecated"]

	item := newDirectiveDefItem(directive, resolver)
	is.Equal(item.Title(), "@deprecated(reason: String = \"No longer supported\")")
	is.Equal(item.Description(), "")

	// Directive items are now openable and show their arguments
	panel, ok := item.OpenPanel()
	is.True(ok)
	is.True(panel != nil)
}
