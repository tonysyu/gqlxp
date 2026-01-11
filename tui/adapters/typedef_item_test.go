package adapters

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/utils/testx"
	"github.com/tonysyu/gqlxp/utils/testx/assert"
)

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
