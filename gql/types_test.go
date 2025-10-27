package gql_test

import (
	"testing"

	"github.com/matryer/is"
	. "github.com/tonysyu/gqlxp/gql"
)

func TestField_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
		type User {
			id: ID!
			name: String!
		}

		type Post {
			title: String!
		}

		union SearchResult = User | Post

		type Query {
			"""Get a user by ID"""
			getUser(id: ID!): User
			search(query: String!, limit: Int = 10): [SearchResult!]!
		}
	`))

	getUser := schema.Query["getUser"]
	search := schema.Query["search"]

	t.Run("Name returns field name", func(t *testing.T) {
		is.Equal(getUser.Name(), "getUser")
		is.Equal(search.Name(), "search")
	})

	t.Run("Description returns field description", func(t *testing.T) {
		is.Equal(getUser.Description(), "Get a user by ID")
		is.Equal(search.Description(), "")
	})

	t.Run("TypeString returns correct type representation", func(t *testing.T) {
		is.Equal(getUser.TypeString(), "User")
		is.Equal(search.TypeString(), "[SearchResult!]!")
	})

	t.Run("TypeName returns unwrapped type name", func(t *testing.T) {
		is.Equal(getUser.TypeName(), "User")
		is.Equal(search.TypeName(), "SearchResult")
	})

	t.Run("Signature returns full call signature", func(t *testing.T) {
		is.Equal(getUser.Signature(), "getUser(id: ID!): User")
		is.Equal(search.Signature(), "search(query: String!, limit: Int = 10): [SearchResult!]!")
	})

	t.Run("Arguments returns field arguments", func(t *testing.T) {
		args := getUser.Arguments()
		is.Equal(len(args), 1)
		is.Equal(args[0].Name(), "id")
		is.Equal(args[0].TypeString(), "ID!")
	})

	t.Run("Arguments with defaults", func(t *testing.T) {
		args := search.Arguments()
		is.Equal(len(args), 2)
		is.Equal(args[0].Name(), "query")
		is.Equal(args[1].Name(), "limit")
	})

	t.Run("ResolveResultType returns correct TypeDef for Object", func(t *testing.T) {
		typeDef, err := getUser.ResolveResultType(&schema)
		is.NoErr(err)

		userObj, ok := typeDef.(*Object)
		is.True(ok)
		is.Equal(userObj.Name(), "User")
	})

	t.Run("ResolveResultType returns correct TypeDef for Union", func(t *testing.T) {
		typeDef, err := search.ResolveResultType(&schema)
		is.NoErr(err)

		union, ok := typeDef.(*Union)
		is.True(ok)
		is.Equal(union.Name(), "SearchResult")
	})
}

func TestArguments_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
		type User {
			id: ID!
		}

		type Query {
			getUser(id: ID!): User
		}
	`))

	userQuery := schema.Query["getUser"]

	t.Run("Name returns input value name", func(t *testing.T) {
		args := userQuery.Arguments()
		is.Equal(args[0].Name(), "id")
	})

	t.Run("TypeString returns correct type", func(t *testing.T) {
		args := userQuery.Arguments()
		is.Equal(args[0].TypeString(), "ID!")
	})

	t.Run("Signature for argument includes type", func(t *testing.T) {
		args := userQuery.Arguments()
		is.Equal(args[0].Signature(), "id: ID!")
	})
}

func TestInputFields_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
		input CreateUserInput {
			"""User's name"""
			name: String!
		}
	`))
	input := schema.Input["CreateUserInput"]
	fields := input.Fields()

	nameField := fields[0]

	t.Run("Description returns input value description", func(t *testing.T) {
		is.Equal(nameField.Name(), "name")
	})

	t.Run("Description returns input value description", func(t *testing.T) {
		is.Equal(nameField.Description(), "User's name")
	})

	t.Run("Signature for input field includes type", func(t *testing.T) {
		is.Equal(nameField.Signature(), "name: String!")
	})
}

func TestObject_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
		"""
		User type
		"""
		type User implements Node {
			id: ID!
			name: String!
			email: String!
		}
	`))
	user := schema.Object["User"]

	t.Run("Name returns object name", func(t *testing.T) {
		is.Equal(user.Name(), "User")
	})

	t.Run("Description returns object description", func(t *testing.T) {
		is.Equal(user.Description(), "User type")
	})

	t.Run("Interfaces returns implemented interfaces", func(t *testing.T) {
		interfaces := user.Interfaces()
		is.Equal(len(interfaces), 1)
		is.Equal(interfaces[0], "Node")
	})

	t.Run("Fields returns object fields", func(t *testing.T) {
		fields := user.Fields()
		is.Equal(len(fields), 3) // id, name, email

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

func TestInputObject_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
		enum Status {
			ACTIVE
			INACTIVE
			PENDING
		}

		"""
		Input for creating a user
		"""
		input CreateUserInput {
			"""User's name"""
			name: String!
			email: String!
			status: Status = ACTIVE
		}
	`))
	input := schema.Input["CreateUserInput"]

	t.Run("Name returns input object name", func(t *testing.T) {
		is.Equal(input.Name(), "CreateUserInput")
	})

	t.Run("Description returns input object description", func(t *testing.T) {
		is.Equal(input.Description(), "Input for creating a user")
	})

	t.Run("Fields returns input object fields", func(t *testing.T) {
		fields := input.Fields()
		is.Equal(len(fields), 3) // name, email, status

		is.Equal(fields[0].Name(), "name")
		is.Equal(fields[0].TypeString(), "String!")
		is.Equal(fields[1].Name(), "email")
		is.Equal(fields[2].Name(), "status")
	})
}

func TestEnum_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
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
	`))
	status := schema.Enum["Status"]

	t.Run("Name returns enum name", func(t *testing.T) {
		is.Equal(status.Name(), "Status")
	})

	t.Run("Description returns enum description", func(t *testing.T) {
		is.Equal(status.Description(), "Represents a user status")
	})

	t.Run("Values returns enum values", func(t *testing.T) {
		values := status.Values()
		is.Equal(len(values), 3)

		testCases := []struct {
			name        string
			description string
		}{
			{"ACTIVE", "Active user"},
			{"INACTIVE", "Inactive user"},
			{"PENDING", ""},
		}

		for i, tc := range testCases {
			is.Equal(values[i].Name(), tc.name)
			is.Equal(values[i].Description(), tc.description)
		}
	})
}

func TestEnumValue_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
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
	`))
	status := schema.Enum["Status"]

	t.Run("Name returns enum value name", func(t *testing.T) {
		values := status.Values()
		is.Equal(values[0].Name(), "ACTIVE")
	})

	t.Run("Description returns enum value description", func(t *testing.T) {
		values := status.Values()
		is.Equal(values[0].Description(), "Active user")
		is.Equal(values[1].Description(), "Inactive user")
	})

	t.Run("EnumValue without description", func(t *testing.T) {
		values := status.Values()
		is.Equal(values[2].Description(), "")
	})
}

func TestScalar_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
		"""
		Custom date scalar
		"""
		scalar Date
	`))
	date := schema.Scalar["Date"]

	t.Run("Name returns scalar name", func(t *testing.T) {
		is.Equal(date.Name(), "Date")
	})

	t.Run("Description returns scalar description", func(t *testing.T) {
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

func TestInterface_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
		"""
		Node interface for entities with IDs
		"""
		interface Node {
			"""Unique identifier"""
			id: ID!
		}
	`))
	node := schema.Interface["Node"]

	t.Run("Name returns interface name", func(t *testing.T) {
		is.Equal(node.Name(), "Node")
	})

	t.Run("Description returns interface description", func(t *testing.T) {
		is.Equal(node.Description(), "Node interface for entities with IDs")
	})

	t.Run("Fields returns interface fields", func(t *testing.T) {
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

func TestUnion_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
		type User { id: ID! }
		type Post { id: ID! }

		"""
		Search result union
		"""
		union SearchResult = User | Post
	`))
	searchResult := schema.Union["SearchResult"]

	t.Run("Name returns union name", func(t *testing.T) {
		is.Equal(searchResult.Name(), "SearchResult")
	})

	t.Run("Description returns union description", func(t *testing.T) {
		is.Equal(searchResult.Description(), "Search result union")
	})

	t.Run("Types returns union member types", func(t *testing.T) {
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

func TestDirective_Methods(t *testing.T) {
	is := is.New(t)
	schema, _ := ParseSchema([]byte(`
		"""
		Marks field as deprecated
		"""
		directive @deprecated(reason: String = "No longer supported") on FIELD_DEFINITION | ENUM_VALUE
	`))
	deprecated := schema.Directive["deprecated"]

	t.Run("Name returns directive name", func(t *testing.T) {
		is.Equal(deprecated.Name(), "deprecated")
	})

	t.Run("Description returns directive description", func(t *testing.T) {
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
