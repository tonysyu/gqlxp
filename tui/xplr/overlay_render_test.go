package xplr

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/utils/testx"
	"github.com/tonysyu/gqlxp/utils/testx/assert"
)

func renderOverlayFromMainModel(gqlType gqlType, schemaString string) string {
	tm := newTestModel(schemaString)
	// Override style to remove border and padding to simplify test comparison
	tm.Model.Overlay.Styles.Overlay = lipgloss.NewStyle()

	tm.Update(tea.WindowSizeMsg{Width: 80, Height: 60})
	tm.Update(SetGQLTypeMsg{GQLType: gqlType})
	tm.Update(tea.KeyMsg{Type: tea.KeySpace})

	return testx.NormalizeView(tm.View())
}

func TestOverlayIntegrationWithMainModel(t *testing.T) {
	assert := assert.New(t)

	t.Run("Render query details", func(t *testing.T) {
		view := renderOverlayFromMainModel(queryType, `
			type Query {
				"""
				Get user by ID
				"""
				getUser(id: ID!): User
			}
		`)

		assert.StringContains(view, testx.NormalizeView(`
			# getUser
            getUser(id: ID!): User
            Get user by ID
		`))
	})

	t.Run("Render mutation details", func(t *testing.T) {
		view := renderOverlayFromMainModel(mutationType, `
			type Mutation {
			    """
			    Create a new post
			    """
			    createPost(title: String!, content: String!, authorId: ID!): Post!
			}
		`)

		assert.StringContains(view, testx.NormalizeView(`
			# createPost
			createPost(title: String!, content: String!, authorId: ID!): Post!
			Create a new post
		`))
	})

	t.Run("Render object type details", func(t *testing.T) {
		view := renderOverlayFromMainModel(objectType, `
			interface Node {
				id: ID!
			}

			"""
			A user in the system
			"""
			type User implements Node {
				"""
				Unique identifier
				"""
				id: ID!
				"""
				User's full name
				"""
				name: String!
				email: String!
			}
		`)

		assert.StringContains(view, testx.NormalizeView(`
			# User

			A user in the system

			**Implements:** Node

			"""Unique identifier"""
			id: ID!

			"""User's full name"""
			name: String!

			email: String!
		`))
	})

	t.Run("Render input type details", func(t *testing.T) {
		view := renderOverlayFromMainModel(inputType, `
			"""
			Input for creating a user
			"""
			input CreateUserInput {
				"""
				User's full name
				"""
				name: String!
				"""
				User's email address
				"""
				email: String!
				age: Int
			}
		`)

		assert.StringContains(view, testx.NormalizeView(`
			# CreateUserInput
			Input for creating a user

			"""User's full name"""
			name: String!

			"""User's email address"""
			email: String!

			age: Int
		`))
	})

	t.Run("Render enum type details", func(t *testing.T) {
		view := renderOverlayFromMainModel(enumType, `
			"""
			User role in the system
			"""
			enum Role {
				"""
				Administrator role
				"""
				ADMIN
				"""
				Regular user role
				"""
				USER
				GUEST
			}
		`)

		assert.StringContains(view, testx.NormalizeView(`
			# Role
			User role in the system

			"""Administrator role"""
			ADMIN

			"""Regular user role"""
			USER

			GUEST
		`))
	})

	t.Run("Render scalar type details", func(t *testing.T) {
		view := renderOverlayFromMainModel(scalarType, `
			"""
			DateTime scalar type
			"""
			scalar DateTime
		`)

		assert.StringContains(view, testx.NormalizeView(`
			# DateTime
			DateTime scalar type
			*Scalar type*
		`))
	})

	t.Run("Render interface type details", func(t *testing.T) {
		view := renderOverlayFromMainModel(interfaceType, `
			"""
			Node interface for entities with ID
			"""
			interface Node {
				"""
				Unique identifier for the entity
				"""
				id: ID!
			}
		`)

		assert.StringContains(view, testx.NormalizeView(`
			# Node
			Node interface for entities with ID

			"""Unique identifier for the entity"""
			id: ID!
		`))
	})

	t.Run("Render union type details", func(t *testing.T) {
		view := renderOverlayFromMainModel(unionType, `
			"""
			Search result can be User or Post
			"""
			union SearchResult = User | Post
		`)

		assert.StringContains(view, testx.NormalizeView(`
			# SearchResult
			Search result can be User or Post

			**Union of:** User | Post
		`))
	})

	t.Run("Render directive type details", func(t *testing.T) {
		view := renderOverlayFromMainModel(directiveType, `
			"""
			Require authentication to access
			"""
			directive @auth(
				requires: String!
			) on FIELD_DEFINITION | OBJECT
		`)

		assert.StringContains(view, testx.NormalizeView(`
			# @auth
			@auth(requires: String!)
			Require authentication to access

			**Locations:**
			• FIELD_DEFINITION
			• OBJECT
		`))
	})
}
