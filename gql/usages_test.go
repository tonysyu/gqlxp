package gql_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
	"slices"
)

func TestBuildUsageIndex_FieldReturnTypes(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		type User {
			id: ID!
			name: String!
		}

		type Post {
			author: User!
			title: String!
		}

		type Query {
			user(id: ID!): User
			users: [User!]!
			post(id: ID!): Post
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	// Test User usages
	userUsages := schema.Usages["User"]
	is.True(userUsages != nil)
	is.Equal(len(userUsages), 3) // Query.user, Query.users, Post.author

	// Collect paths and verify all usages
	paths := make([]string, len(userUsages))
	for i, usage := range userUsages {
		paths[i] = usage.Path
	}

	// Verify all three paths exist
	is.True(slices.Contains(paths, "Query.user"))
	is.True(slices.Contains(paths, "Query.users"))
	is.True(slices.Contains(paths, "Post.author"))

	// Test Post usages
	postUsages := schema.Usages["Post"]
	is.True(postUsages != nil)
	is.Equal(len(postUsages), 1) // Query.post

	is.Equal(postUsages[0].ParentType, "Query")
	is.Equal(postUsages[0].FieldName, "post")
	is.Equal(postUsages[0].Path, "Query.post")
}

func TestBuildUsageIndex_MutationFields(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		type User {
			id: ID!
			name: String!
		}

		type Mutation {
			createUser(name: String!): User!
			updateUser(id: ID!, name: String!): User
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	userUsages := schema.Usages["User"]
	is.True(userUsages != nil)
	is.Equal(len(userUsages), 2) // Mutation.createUser, Mutation.updateUser

	// Collect paths (order not guaranteed due to map iteration)
	paths := make([]string, len(userUsages))
	for i, usage := range userUsages {
		paths[i] = usage.Path
		// Verify all have correct ParentType and ParentKind
		is.Equal(usage.ParentType, "Mutation")
		is.Equal(usage.ParentKind, "Mutation")
	}

	// Verify both paths exist
	is.True(slices.Contains(paths, "Mutation.createUser"))
	is.True(slices.Contains(paths, "Mutation.updateUser"))
}

func TestBuildUsageIndex_NestedFields(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		type User {
			id: ID!
			name: String!
			friends: [User!]!
		}

		type Query {
			user(id: ID!): User
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	userUsages := schema.Usages["User"]
	is.True(userUsages != nil)
	is.Equal(len(userUsages), 2) // Query.user, User.friends (self-reference)

	// Collect paths (order not guaranteed due to map iteration)
	paths := make([]string, len(userUsages))
	for i, usage := range userUsages {
		paths[i] = usage.Path
	}

	// Verify both paths exist
	is.True(slices.Contains(paths, "Query.user"))
	is.True(slices.Contains(paths, "User.friends"))
}

func TestBuildUsageIndex_WrappedTypes(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		type User {
			id: ID!
		}

		type Query {
			user: User
			users: [User]
			requiredUser: User!
			requiredUsers: [User!]!
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	userUsages := schema.Usages["User"]
	is.True(userUsages != nil)
	is.Equal(len(userUsages), 4) // All four variants should be tracked

	paths := make([]string, len(userUsages))
	for i, usage := range userUsages {
		paths[i] = usage.Path
	}

	// All should reference User despite different wrapping
	is.True(slices.Contains(paths, "Query.user"))
	is.True(slices.Contains(paths, "Query.users"))
	is.True(slices.Contains(paths, "Query.requiredUser"))
	is.True(slices.Contains(paths, "Query.requiredUsers"))
}

func TestBuildUsageIndex_NoUsages(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		type User {
			id: ID!
			name: String!
		}

		type Post {
			title: String!
		}

		type Query {
			post: Post
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	// User has no usages
	userUsages := schema.Usages["User"]
	is.Equal(len(userUsages), 0)

	// Post has one usage
	postUsages := schema.Usages["Post"]
	is.True(postUsages != nil)
	is.Equal(len(postUsages), 1)
}

func TestBuildUsageIndex_InterfaceFields(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		type User {
			id: ID!
		}

		interface Node {
			id: ID!
			owner: User
		}

		type Query {
			node(id: ID!): Node
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	// User should be used by Interface field
	userUsages := schema.Usages["User"]
	is.True(userUsages != nil)
	is.Equal(len(userUsages), 1)

	is.Equal(userUsages[0].ParentType, "Node")
	is.Equal(userUsages[0].ParentKind, "Interface")
	is.Equal(userUsages[0].FieldName, "owner")
	is.Equal(userUsages[0].Path, "Node.owner")
}

func TestBuildUsageIndex_ArgumentTypes(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		input UserInput {
			name: String!
			age: Int!
		}

		enum Role {
			ADMIN
			USER
		}

		scalar DateTime

		type User {
			id: ID!
		}

		type Query {
			user(id: ID!, role: Role): User
			createUser(input: UserInput!): User!
			search(createdAfter: DateTime): [User!]!
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	// Test UserInput usages (in arguments)
	inputUsages := schema.Usages["UserInput"]
	is.True(inputUsages != nil)
	is.Equal(len(inputUsages), 1)
	is.True(slices.Contains(getPaths(inputUsages), "Query.createUser(input: UserInput)"))

	// Test Role (enum) usages (in arguments)
	roleUsages := schema.Usages["Role"]
	is.True(roleUsages != nil)
	is.Equal(len(roleUsages), 1)
	is.True(slices.Contains(getPaths(roleUsages), "Query.user(role: Role)"))

	// Test DateTime (scalar) usages (in arguments)
	dateTimeUsages := schema.Usages["DateTime"]
	is.True(dateTimeUsages != nil)
	is.Equal(len(dateTimeUsages), 1)
	is.True(slices.Contains(getPaths(dateTimeUsages), "Query.search(createdAfter: DateTime)"))
}

func TestBuildUsageIndex_InputObjectFieldTypes(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		enum Status {
			ACTIVE
			INACTIVE
		}

		scalar Email

		input AddressInput {
			street: String!
			city: String!
		}

		input UserInput {
			name: String!
			email: Email!
			status: Status!
			address: AddressInput!
		}

		type User {
			id: ID!
		}

		type Mutation {
			createUser(input: UserInput!): User!
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	// Test Email (scalar) used in InputObject field
	emailUsages := schema.Usages["Email"]
	is.True(emailUsages != nil)
	is.Equal(len(emailUsages), 1)
	is.True(slices.Contains(getPaths(emailUsages), "UserInput.email"))

	// Test Status (enum) used in InputObject field
	statusUsages := schema.Usages["Status"]
	is.True(statusUsages != nil)
	is.Equal(len(statusUsages), 1)
	is.True(slices.Contains(getPaths(statusUsages), "UserInput.status"))

	// Test AddressInput used in InputObject field
	addressUsages := schema.Usages["AddressInput"]
	is.True(addressUsages != nil)
	is.Equal(len(addressUsages), 1)
	is.True(slices.Contains(getPaths(addressUsages), "UserInput.address"))
}

func TestBuildUsageIndex_InterfaceImplementations(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		interface Node {
			id: ID!
		}

		interface Entity {
			id: ID!
			createdAt: String!
		}

		interface Timestamped {
			createdAt: String!
		}

		type User implements Node & Entity {
			id: ID!
			createdAt: String!
			name: String!
		}

		type Post implements Node & Timestamped {
			id: ID!
			createdAt: String!
			title: String!
		}

		type Query {
			node(id: ID!): Node
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	// Test Node interface implementations
	nodeUsages := schema.Usages["Node"]
	is.True(nodeUsages != nil)
	// Should have: Query.node (field), User implements Node, Post implements Node
	is.Equal(len(nodeUsages), 3)
	paths := getPaths(nodeUsages)
	is.True(slices.Contains(paths, "Query.node"))
	is.True(slices.Contains(paths, "User"))
	is.True(slices.Contains(paths, "Post"))

	// Test Entity interface implementation
	entityUsages := schema.Usages["Entity"]
	is.True(entityUsages != nil)
	is.Equal(len(entityUsages), 1)
	is.True(slices.Contains(getPaths(entityUsages), "User"))

	// Test Timestamped interface implementation
	timestampedUsages := schema.Usages["Timestamped"]
	is.True(timestampedUsages != nil)
	is.Equal(len(timestampedUsages), 1)
	is.True(slices.Contains(getPaths(timestampedUsages), "Post"))
}

func TestBuildUsageIndex_InterfaceImplementsInterface(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		interface Node {
			id: ID!
		}

		interface Resource implements Node {
			id: ID!
			url: String!
		}

		type Image implements Resource & Node {
			id: ID!
			url: String!
			width: Int!
		}

		type Query {
			resource(id: ID!): Resource
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	// Test Node interface usages
	nodeUsages := schema.Usages["Node"]
	is.True(nodeUsages != nil)
	paths := getPaths(nodeUsages)
	// Should include: Resource implements Node, Image implements Node
	is.True(slices.Contains(paths, "Resource"))
	is.True(slices.Contains(paths, "Image"))
}

func TestBuildUsageIndex_UnionMemberTypes(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		type User {
			id: ID!
			name: String!
		}

		type Post {
			id: ID!
			title: String!
		}

		type Comment {
			id: ID!
			text: String!
		}

		union SearchResult = User | Post | Comment

		type Query {
			search(query: String!): [SearchResult!]!
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	// Test User as union member
	userUsages := schema.Usages["User"]
	is.True(userUsages != nil)
	paths := getPaths(userUsages)
	is.True(slices.Contains(paths, "SearchResult"))

	// Test Post as union member
	postUsages := schema.Usages["Post"]
	is.True(postUsages != nil)
	paths = getPaths(postUsages)
	is.True(slices.Contains(paths, "SearchResult"))

	// Test Comment as union member
	commentUsages := schema.Usages["Comment"]
	is.True(commentUsages != nil)
	paths = getPaths(commentUsages)
	is.True(slices.Contains(paths, "SearchResult"))
}

func TestBuildUsageIndex_DirectiveArgumentTypes(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		enum CacheControl {
			PUBLIC
			PRIVATE
		}

		input CacheConfig {
			maxAge: Int!
		}

		directive @cache(
			control: CacheControl!
			config: CacheConfig
		) on FIELD_DEFINITION

		type Query {
			user: String @cache(control: PUBLIC)
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	// Test CacheControl enum used in directive argument
	cacheControlUsages := schema.Usages["CacheControl"]
	is.True(cacheControlUsages != nil)
	is.Equal(len(cacheControlUsages), 1)
	is.True(slices.Contains(getPaths(cacheControlUsages), "cache(control: CacheControl)"))

	// Test CacheConfig input used in directive argument
	cacheConfigUsages := schema.Usages["CacheConfig"]
	is.True(cacheConfigUsages != nil)
	is.Equal(len(cacheConfigUsages), 1)
	is.True(slices.Contains(getPaths(cacheConfigUsages), "cache(config: CacheConfig)"))
}

func TestBuildUsageIndex_DirectiveUsages(t *testing.T) {
	is := is.New(t)

	schemaContent := []byte(`
		directive @deprecated(
			reason: String = "No longer supported"
		) on FIELD_DEFINITION | ENUM_VALUE

		directive @auth on OBJECT | FIELD_DEFINITION

		scalar DateTime @specifiedBy(url: "https://tools.ietf.org/html/rfc3339")

		enum Role @deprecated(reason: "Use UserRole instead") {
			ADMIN
			USER @deprecated(reason: "Use MEMBER instead")
		}

		type User @auth {
			id: ID!
			email: String! @deprecated(reason: "Use contactEmail")
		}

		interface Node {
			id: ID! @auth
		}

		union SearchResult @deprecated = User

		input UserInput @auth {
			name: String!
		}

		type Query {
			user(id: ID! @auth): User @deprecated
		}
	`)

	schema, err := gql.ParseSchema(schemaContent)
	is.NoErr(err)

	// Test @deprecated directive usages
	deprecatedUsages := schema.Usages["deprecated"]
	is.True(deprecatedUsages != nil)
	paths := getPaths(deprecatedUsages)
	is.True(slices.Contains(paths, "Role"))         // on Enum type
	is.True(slices.Contains(paths, "Role.USER"))    // on EnumValue
	is.True(slices.Contains(paths, "User.email"))   // on Object field
	is.True(slices.Contains(paths, "SearchResult")) // on Union type
	is.True(slices.Contains(paths, "Query.user"))   // on Query field

	// Test @auth directive usages
	authUsages := schema.Usages["auth"]
	is.True(authUsages != nil)
	paths = getPaths(authUsages)
	is.True(slices.Contains(paths, "User"))           // on Object type
	is.True(slices.Contains(paths, "Node.id"))        // on Interface field
	is.True(slices.Contains(paths, "UserInput"))      // on Input type
	is.True(slices.Contains(paths, "Query.user(id)")) // on argument

	// Test @specifiedBy directive usages
	specifiedByUsages := schema.Usages["specifiedBy"]
	is.True(specifiedByUsages != nil)
	is.Equal(len(specifiedByUsages), 1)
	is.True(slices.Contains(getPaths(specifiedByUsages), "DateTime")) // on Scalar type
}

// getPaths extracts all paths from a slice of usages
func getPaths(usages []*gql.Usage) []string {
	paths := make([]string, len(usages))
	for i, usage := range usages {
		paths[i] = usage.Path
	}
	return paths
}
