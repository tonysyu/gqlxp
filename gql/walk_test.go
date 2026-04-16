package gql_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/gql"
)

const walkTestSchema = `
	interface Node {
		id: ID!
	}

	type User implements Node {
		id: ID!
		name: String!
	}

	input CreateUserInput {
		name: String!
	}

	enum Role {
		ADMIN
		USER
	}

	union SearchResult = User

	type Query {
		user(id: ID!): User
	}

	type Mutation {
		createUser(input: CreateUserInput!): User
	}
`

func TestWalk_VisitsAllTopLevelKinds(t *testing.T) {
	is := is.New(t)

	schema, err := gql.ParseSchema([]byte(walkTestSchema))
	is.NoErr(err)

	var kinds []string
	schema.Walk(gql.SchemaVisitor{
		VisitField: func(ctx gql.VisitContext, name string, _ *gql.Field) {
			kinds = append(kinds, ctx.Kind+":"+name)
		},
		VisitObject: func(ctx gql.VisitContext, name string, _ *gql.Object) {
			kinds = append(kinds, "Object:"+name)
		},
		VisitInterface: func(ctx gql.VisitContext, name string, _ *gql.Interface) {
			kinds = append(kinds, "Interface:"+name)
		},
		VisitInput: func(ctx gql.VisitContext, name string, _ *gql.InputObject) {
			kinds = append(kinds, "Input:"+name)
		},
		VisitEnum: func(ctx gql.VisitContext, name string, _ *gql.Enum) {
			kinds = append(kinds, "Enum:"+name)
		},
		VisitUnion: func(ctx gql.VisitContext, name string, _ *gql.Union) {
			kinds = append(kinds, "Union:"+name)
		},
	})

	is.True(contains(kinds, "Query:user"))
	is.True(contains(kinds, "Mutation:createUser"))
	is.True(contains(kinds, "Object:User"))
	is.True(contains(kinds, "Interface:Node"))
	is.True(contains(kinds, "Input:CreateUserInput"))
	is.True(contains(kinds, "Enum:Role"))
	is.True(contains(kinds, "Union:SearchResult"))
}

func TestWalk_VisitsSubElements(t *testing.T) {
	is := is.New(t)

	schema, err := gql.ParseSchema([]byte(walkTestSchema))
	is.NoErr(err)

	var objectFields, interfaceFields, inputFields, enumValues []string

	schema.Walk(gql.SchemaVisitor{
		VisitObjectField: func(ctx gql.VisitContext, field *gql.Field) {
			objectFields = append(objectFields, ctx.ParentName+"."+field.Name())
		},
		VisitInterfaceField: func(ctx gql.VisitContext, field *gql.Field) {
			interfaceFields = append(interfaceFields, ctx.ParentName+"."+field.Name())
		},
		VisitInputField: func(ctx gql.VisitContext, field *gql.Field) {
			inputFields = append(inputFields, ctx.ParentName+"."+field.Name())
		},
		VisitEnumValue: func(ctx gql.VisitContext, value *gql.EnumValue) {
			enumValues = append(enumValues, ctx.ParentName+"."+value.Name())
		},
	})

	is.True(contains(objectFields, "User.id"))
	is.True(contains(objectFields, "User.name"))
	is.True(contains(interfaceFields, "Node.id"))
	is.True(contains(inputFields, "CreateUserInput.name"))
	is.True(contains(enumValues, "Role.ADMIN"))
	is.True(contains(enumValues, "Role.USER"))
}

func TestWalk_KindOrderIsStable(t *testing.T) {
	is := is.New(t)

	schema, err := gql.ParseSchema([]byte(walkTestSchema))
	is.NoErr(err)

	// Run Walk twice and verify we get the same order
	collect := func() []string {
		var names []string
		schema.Walk(gql.SchemaVisitor{
			VisitField: func(ctx gql.VisitContext, name string, _ *gql.Field) {
				names = append(names, ctx.Kind+":"+name)
			},
			VisitObject: func(ctx gql.VisitContext, name string, _ *gql.Object) {
				names = append(names, "Object:"+name)
			},
		})
		return names
	}

	first := collect()
	second := collect()
	is.Equal(first, second)
}

func TestWalk_NilCallbacksAreSkipped(t *testing.T) {
	is := is.New(t)

	schema, err := gql.ParseSchema([]byte(walkTestSchema))
	is.NoErr(err)

	// Walk with only VisitObject set — should not panic
	var names []string
	schema.Walk(gql.SchemaVisitor{
		VisitObject: func(ctx gql.VisitContext, name string, _ *gql.Object) {
			names = append(names, name)
		},
	})
	is.True(len(names) > 0)
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
