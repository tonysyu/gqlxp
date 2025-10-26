package gql_test

import (
	"testing"

	"github.com/matryer/is"
	. "github.com/tonysyu/gqlxp/gql"
)

// testSchema is a comprehensive schema used across multiple tests
const testSchema = `
	"""
	Represents a user status
	"""
	enum Status {
		"""Active user"""
		ACTIVE
		"""Inactive user"""
		INACTIVE
		PENDING
	}

	"""
	Custom date scalar
	"""
	scalar Date

	"""
	Input for creating a user
	"""
	input CreateUserInput {
		"""User's name"""
		name: String!
		email: String!
		status: Status = ACTIVE
	}

	"""
	Node interface for entities with IDs
	"""
	interface Node {
		"""Unique identifier"""
		id: ID!
	}

	"""
	Search result union
	"""
	union SearchResult = User | Post

	"""
	Marks field as deprecated
	"""
	directive @deprecated(reason: String = "No longer supported") on FIELD_DEFINITION | ENUM_VALUE

	"""
	User type
	"""
	type User implements Node {
		id: ID!
		name: String!
		email: String!
		status: Status!
		posts: [Post!]!
	}

	"""
	Post type
	"""
	type Post implements Node {
		id: ID!
		title: String!
		content: String
		author: User!
	}

	type Query {
		"""Get a user by ID"""
		getUser(id: ID!): User
		"""Search all content"""
		search(query: String!, limit: Int = 10): [SearchResult!]!
	}

	type Mutation {
		"""Create a new user"""
		createUser(input: CreateUserInput!): User!
	}
`

func TestField_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(testSchema))

	t.Run("Name returns field name", func(t *testing.T) {
		field := schema.Query["getUser"]
		is.Equal(field.Name(), "getUser")
	})

	t.Run("Description returns field description", func(t *testing.T) {
		field := schema.Query["getUser"]
		is.Equal(field.Description(), "Get a user by ID")
	})

	t.Run("TypeString returns correct type representation", func(t *testing.T) {
		testCases := []struct {
			fieldName    string
			expectedType string
		}{
			{"getUser", "User"},
			{"search", "[SearchResult!]!"},
		}

		for _, tc := range testCases {
			field := schema.Query[tc.fieldName]
			is.Equal(field.TypeString(), tc.expectedType)
		}
	})

	t.Run("TypeName returns unwrapped type name", func(t *testing.T) {
		field := schema.Query["search"]
		is.Equal(field.TypeName(), "SearchResult")
	})

	t.Run("Arguments returns field arguments", func(t *testing.T) {
		field := schema.Query["getUser"]
		args := field.Arguments()
		is.Equal(len(args), 1)
		is.Equal(args[0].Name(), "id")
		is.Equal(args[0].TypeString(), "ID!")
	})

	t.Run("Arguments with defaults", func(t *testing.T) {
		field := schema.Query["search"]
		args := field.Arguments()
		is.Equal(len(args), 2)
		is.Equal(args[0].Name(), "query")
		is.Equal(args[1].Name(), "limit")
	})

	t.Run("Signature includes arguments", func(t *testing.T) {
		field := schema.Query["getUser"]
		sig := field.Signature()
		is.True(sig != "")
		// Signature should contain the field name
		is.True(len(sig) > len("getUser"))
	})

	t.Run("ResolveResultType returns correct TypeDef for Object", func(t *testing.T) {
		field := schema.Query["getUser"]
		typeDef, err := field.ResolveResultType(&schema)
		is.NoErr(err)

		userObj, ok := typeDef.(*Object)
		is.True(ok)
		is.Equal(userObj.Name(), "User")
	})

	t.Run("ResolveResultType returns correct TypeDef for Union", func(t *testing.T) {
		field := schema.Query["search"]
		typeDef, err := field.ResolveResultType(&schema)
		is.NoErr(err)

		union, ok := typeDef.(*Union)
		is.True(ok)
		is.Equal(union.Name(), "SearchResult")
	})
}

func TestField_NewField(t *testing.T) {
	is := is.New(t)

	t.Run("NewField returns nil for nil input", func(t *testing.T) {
		field := NewField(nil)
		is.True(field == nil)
	})
}

func TestInputValue_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(testSchema))

	t.Run("Name returns input value name", func(t *testing.T) {
		field := schema.Query["getUser"]
		args := field.Arguments()
		is.Equal(args[0].Name(), "id")
	})

	t.Run("Description returns input value description", func(t *testing.T) {
		input := schema.Input["CreateUserInput"]
		fields := input.Fields()
		is.Equal(fields[0].Description(), "User's name")
	})

	t.Run("TypeString returns correct type", func(t *testing.T) {
		field := schema.Query["getUser"]
		args := field.Arguments()
		is.Equal(args[0].TypeString(), "ID!")
	})

	t.Run("Signature for argument includes type", func(t *testing.T) {
		field := schema.Query["getUser"]
		args := field.Arguments()
		sig := args[0].Signature()
		is.True(sig != "")
		// Should contain both name and type
		is.True(len(sig) > len("id"))
	})

	t.Run("Signature for input field includes type", func(t *testing.T) {
		input := schema.Input["CreateUserInput"]
		fields := input.Fields()
		sig := fields[0].Signature()
		is.True(sig != "")
		is.True(len(sig) > len("name"))
	})
}

func TestInputValue_Constructors(t *testing.T) {
	is := is.New(t)

	t.Run("NewInputValue returns nil for nil input", func(t *testing.T) {
		inputValue := NewInputValue(nil)
		is.True(inputValue == nil)
	})

	t.Run("NewInputValueFromField returns nil for nil input", func(t *testing.T) {
		inputValue := NewInputValueFromField(nil)
		is.True(inputValue == nil)
	})
}

func TestObject_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(testSchema))

	t.Run("Name returns object name", func(t *testing.T) {
		user := schema.Object["User"]
		is.Equal(user.Name(), "User")
	})

	t.Run("Description returns object description", func(t *testing.T) {
		user := schema.Object["User"]
		is.Equal(user.Description(), "User type")
	})

	t.Run("Interfaces returns implemented interfaces", func(t *testing.T) {
		user := schema.Object["User"]
		interfaces := user.Interfaces()
		is.Equal(len(interfaces), 1)
		is.Equal(interfaces[0], "Node")
	})

	t.Run("Fields returns object fields", func(t *testing.T) {
		user := schema.Object["User"]
		fields := user.Fields()
		is.Equal(len(fields), 5) // id, name, email, status, posts

		// Check first field
		is.Equal(fields[0].Name(), "id")
		is.Equal(fields[0].TypeString(), "ID!")
	})

	t.Run("Object without interfaces", func(t *testing.T) {
		schemaStr := `
			type Simple {
				id: ID!
			}
			type Query {
				simple: Simple
			}
		`
		s, _ := ParseSchema([]byte(schemaStr))
		simple := s.Object["Simple"]
		is.Equal(len(simple.Interfaces()), 0)
	})
}

func TestObject_NewObject(t *testing.T) {
	is := is.New(t)

	t.Run("NewObject returns nil for nil input", func(t *testing.T) {
		obj := NewObject(nil)
		is.True(obj == nil)
	})
}

func TestInputObject_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(testSchema))

	t.Run("Name returns input object name", func(t *testing.T) {
		input := schema.Input["CreateUserInput"]
		is.Equal(input.Name(), "CreateUserInput")
	})

	t.Run("Description returns input object description", func(t *testing.T) {
		input := schema.Input["CreateUserInput"]
		is.Equal(input.Description(), "Input for creating a user")
	})

	t.Run("Fields returns input object fields", func(t *testing.T) {
		input := schema.Input["CreateUserInput"]
		fields := input.Fields()
		is.Equal(len(fields), 3) // name, email, status

		is.Equal(fields[0].Name(), "name")
		is.Equal(fields[0].TypeString(), "String!")
		is.Equal(fields[1].Name(), "email")
		is.Equal(fields[2].Name(), "status")
	})
}

func TestInputObject_NewInputObject(t *testing.T) {
	is := is.New(t)

	t.Run("NewInputObject returns nil for nil input", func(t *testing.T) {
		obj := NewInputObject(nil)
		is.True(obj == nil)
	})
}

func TestEnum_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(testSchema))

	t.Run("Name returns enum name", func(t *testing.T) {
		status := schema.Enum["Status"]
		is.Equal(status.Name(), "Status")
	})

	t.Run("Description returns enum description", func(t *testing.T) {
		status := schema.Enum["Status"]
		is.Equal(status.Description(), "Represents a user status")
	})

	t.Run("Values returns enum values", func(t *testing.T) {
		status := schema.Enum["Status"]
		values := status.Values()
		is.Equal(len(values), 3)

		is.Equal(values[0].Name(), "ACTIVE")
		is.Equal(values[1].Name(), "INACTIVE")
		is.Equal(values[2].Name(), "PENDING")
	})
}

func TestEnum_NewEnum(t *testing.T) {
	is := is.New(t)

	t.Run("NewEnum returns nil for nil input", func(t *testing.T) {
		enum := NewEnum(nil)
		is.True(enum == nil)
	})
}

func TestEnumValue_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(testSchema))

	t.Run("Name returns enum value name", func(t *testing.T) {
		status := schema.Enum["Status"]
		values := status.Values()
		is.Equal(values[0].Name(), "ACTIVE")
	})

	t.Run("Description returns enum value description", func(t *testing.T) {
		status := schema.Enum["Status"]
		values := status.Values()
		is.Equal(values[0].Description(), "Active user")
		is.Equal(values[1].Description(), "Inactive user")
	})

	t.Run("EnumValue without description", func(t *testing.T) {
		status := schema.Enum["Status"]
		values := status.Values()
		is.Equal(values[2].Description(), "")
	})
}

func TestEnumValue_NewEnumValue(t *testing.T) {
	is := is.New(t)

	t.Run("NewEnumValue returns nil for nil input", func(t *testing.T) {
		ev := NewEnumValue(nil)
		is.True(ev == nil)
	})
}

func TestScalar_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(testSchema))

	t.Run("Name returns scalar name", func(t *testing.T) {
		date := schema.Scalar["Date"]
		is.Equal(date.Name(), "Date")
	})

	t.Run("Description returns scalar description", func(t *testing.T) {
		date := schema.Scalar["Date"]
		is.Equal(date.Description(), "Custom date scalar")
	})

	t.Run("Scalar without description", func(t *testing.T) {
		schemaStr := `
			scalar CustomScalar
			type Query {
				test: String
			}
		`
		s, _ := ParseSchema([]byte(schemaStr))
		custom := s.Scalar["CustomScalar"]
		is.Equal(custom.Description(), "")
	})
}

func TestScalar_NewScalar(t *testing.T) {
	is := is.New(t)

	t.Run("NewScalar returns nil for nil input", func(t *testing.T) {
		scalar := NewScalar(nil)
		is.True(scalar == nil)
	})
}

func TestInterface_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(testSchema))

	t.Run("Name returns interface name", func(t *testing.T) {
		node := schema.Interface["Node"]
		is.Equal(node.Name(), "Node")
	})

	t.Run("Description returns interface description", func(t *testing.T) {
		node := schema.Interface["Node"]
		is.Equal(node.Description(), "Node interface for entities with IDs")
	})

	t.Run("Fields returns interface fields", func(t *testing.T) {
		node := schema.Interface["Node"]
		fields := node.Fields()
		is.Equal(len(fields), 1)

		is.Equal(fields[0].Name(), "id")
		is.Equal(fields[0].TypeString(), "ID!")
		is.Equal(fields[0].Description(), "Unique identifier")
	})

	t.Run("Interface with multiple fields", func(t *testing.T) {
		schemaStr := `
			interface MultiField {
				id: ID!
				name: String!
				count: Int
			}
			type Query {
				test: String
			}
		`
		s, _ := ParseSchema([]byte(schemaStr))
		multi := s.Interface["MultiField"]
		fields := multi.Fields()
		is.Equal(len(fields), 3)
	})
}

func TestInterface_NewInterface(t *testing.T) {
	is := is.New(t)

	t.Run("NewInterface returns nil for nil input", func(t *testing.T) {
		iface := NewInterface(nil)
		is.True(iface == nil)
	})
}

func TestUnion_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(testSchema))

	t.Run("Name returns union name", func(t *testing.T) {
		searchResult := schema.Union["SearchResult"]
		is.Equal(searchResult.Name(), "SearchResult")
	})

	t.Run("Description returns union description", func(t *testing.T) {
		searchResult := schema.Union["SearchResult"]
		is.Equal(searchResult.Description(), "Search result union")
	})

	t.Run("Types returns union member types", func(t *testing.T) {
		searchResult := schema.Union["SearchResult"]
		types := searchResult.Types()
		is.Equal(len(types), 2)
		is.Equal(types[0], "User")
		is.Equal(types[1], "Post")
	})

	t.Run("Union with multiple types", func(t *testing.T) {
		schemaStr := `
			type A { id: ID! }
			type B { id: ID! }
			type C { id: ID! }
			union Multi = A | B | C
			type Query {
				test: Multi
			}
		`
		s, _ := ParseSchema([]byte(schemaStr))
		multi := s.Union["Multi"]
		types := multi.Types()
		is.Equal(len(types), 3)
		is.Equal(types[0], "A")
		is.Equal(types[1], "B")
		is.Equal(types[2], "C")
	})
}

func TestUnion_NewUnion(t *testing.T) {
	is := is.New(t)

	t.Run("NewUnion returns nil for nil input", func(t *testing.T) {
		union := NewUnion(nil)
		is.True(union == nil)
	})
}

func TestDirective_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(testSchema))

	t.Run("Name returns directive name", func(t *testing.T) {
		deprecated := schema.Directive["deprecated"]
		is.Equal(deprecated.Name(), "deprecated")
	})

	t.Run("Description returns directive description", func(t *testing.T) {
		deprecated := schema.Directive["deprecated"]
		is.Equal(deprecated.Description(), "Marks field as deprecated")
	})

	t.Run("Directive without description", func(t *testing.T) {
		schemaStr := `
			directive @custom on FIELD_DEFINITION
			type Query {
				test: String
			}
		`
		s, _ := ParseSchema([]byte(schemaStr))
		custom := s.Directive["custom"]
		is.Equal(custom.Description(), "")
	})
}

func TestDirective_NewDirective(t *testing.T) {
	is := is.New(t)

	t.Run("NewDirective returns nil for nil input", func(t *testing.T) {
		directive := NewDirective(nil)
		is.True(directive == nil)
	})
}

func TestField_EmptyDescription(t *testing.T) {
	is := is.New(t)

	schemaStr := `
		type Query {
			fieldWithoutDescription: String
		}
	`
	schema, _ := ParseSchema([]byte(schemaStr))
	field := schema.Query["fieldWithoutDescription"]
	is.Equal(field.Description(), "")
}

func TestObject_EmptyDescription(t *testing.T) {
	is := is.New(t)

	schemaStr := `
		type NoDesc {
			id: ID!
		}
		type Query {
			test: NoDesc
		}
	`
	schema, _ := ParseSchema([]byte(schemaStr))
	obj := schema.Object["NoDesc"]
	is.Equal(obj.Description(), "")
}
