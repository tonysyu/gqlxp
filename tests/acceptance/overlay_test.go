package acceptance

import (
	"testing"

	"github.com/tonysyu/gqlxp/tui/xplr/navigation"
)

// TestOverlayRendering verifies that the overlay displays correct details
// for different GraphQL types when opened from the explorer
func TestOverlayRendering(t *testing.T) {
	t.Run("Render query details", func(t *testing.T) {
		schema := `
			type Query {
				"""
				Get user by ID
				"""
				getUser(id: ID!): User
			}
		`
		h := New(t, schema, WithWindowSize(80, 60), WithoutOverlayBorders())
		h.overlay.OpenForType(navigation.QueryType)

		h.assert.OverlayContains(`
			# getUser
            getUser(id: ID!): User
            Get user by ID
		`)
	})

	t.Run("Render mutation details", func(t *testing.T) {
		schema := `
			type Mutation {
			    """
			    Create a new post
			    """
			    createPost(title: String!, content: String!, authorId: ID!): Post!
			}
		`
		h := New(t, schema, WithWindowSize(80, 60), WithoutOverlayBorders())
		h.overlay.OpenForType(navigation.MutationType)

		h.assert.OverlayContains(`
			# createPost
			createPost(title: String!, content: String!, authorId: ID!): Post!
			Create a new post
		`)
	})

	t.Run("Render object type details", func(t *testing.T) {
		schema := `
			type Query { placeholder: String }

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
		`
		h := New(t, schema, WithWindowSize(80, 60), WithoutOverlayBorders())
		h.overlay.OpenForType(navigation.ObjectType)

		h.assert.OverlayContains(`
			# User

			A user in the system

			**Implements:** Node

			"""Unique identifier"""
			id: ID!

			"""User's full name"""
			name: String!

			email: String!
		`)
	})

	t.Run("Render input type details", func(t *testing.T) {
		schema := `
			type Query { placeholder: String }

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
		`
		h := New(t, schema, WithWindowSize(80, 60), WithoutOverlayBorders())
		h.overlay.OpenForType(navigation.InputType)

		h.assert.OverlayContains(`
			# CreateUserInput
			Input for creating a user

			"""User's full name"""
			name: String!

			"""User's email address"""
			email: String!

			age: Int
		`)
	})

	t.Run("Render enum type details", func(t *testing.T) {
		schema := `
			type Query { placeholder: String }

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
		`
		h := New(t, schema, WithWindowSize(80, 60), WithoutOverlayBorders())
		h.overlay.OpenForType(navigation.EnumType)

		h.assert.OverlayContains(`
			# Role
			User role in the system

			"""Administrator role"""
			ADMIN

			"""Regular user role"""
			USER

			GUEST
		`)
	})

	t.Run("Render scalar type details", func(t *testing.T) {
		schema := `
			type Query { placeholder: String }

			"""
			DateTime scalar type
			"""
			scalar DateTime
		`
		h := New(t, schema, WithWindowSize(80, 60), WithoutOverlayBorders())
		h.overlay.OpenForType(navigation.ScalarType)

		h.assert.OverlayContains(`
			# DateTime
			DateTime scalar type
			*Scalar type*
		`)
	})

	t.Run("Render interface type details", func(t *testing.T) {
		schema := `
			type Query { placeholder: String }

			"""
			Node interface for entities with ID
			"""
			interface Node {
				"""
				Unique identifier for the entity
				"""
				id: ID!
			}
		`
		h := New(t, schema, WithWindowSize(80, 60), WithoutOverlayBorders())
		h.overlay.OpenForType(navigation.InterfaceType)

		h.assert.OverlayContains(`
			# Node
			Node interface for entities with ID

			"""Unique identifier for the entity"""
			id: ID!
		`)
	})

	t.Run("Render union type details", func(t *testing.T) {
		schema := `
			type Query { placeholder: String }

			"""
			Search result can be User or Post
			"""
			union SearchResult = User | Post

			type User { id: ID! }
			type Post { id: ID! }
		`
		h := New(t, schema, WithWindowSize(80, 60), WithoutOverlayBorders())
		h.overlay.OpenForType(navigation.UnionType)

		h.assert.OverlayContains(`
			# SearchResult
			Search result can be User or Post

			**Union of:** User | Post
		`)
	})

	t.Run("Render directive type details", func(t *testing.T) {
		schema := `
			type Query { placeholder: String }

			"""
			Require authentication to access
			"""
			directive @auth(
				requires: String!
			) on FIELD_DEFINITION | OBJECT
		`
		h := New(t, schema, WithWindowSize(80, 60), WithoutOverlayBorders())
		h.overlay.OpenForType(navigation.DirectiveType)

		h.assert.OverlayContains(`
			# @auth
			@auth(requires: String!)
			Require authentication to access

			**Locations:**
			• FIELD_DEFINITION
			• OBJECT
		`)
	})
}
