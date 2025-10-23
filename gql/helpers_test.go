package gql_test

import (
	"testing"

	"github.com/matryer/is"
	. "github.com/tonysyu/gqlxp/gql"
)

// Test GetTypeString handles types wrapped in non-nulls and lists
func TestGetTypeString(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Query {
		  getString: String
		  getRequiredString: String!
		  getStringList: [String]
		  getRequiredStringList: [String]!
		  getListOfRequiredStrings: [String!]
		  getRequiredListOfRequiredStrings: [String!]!
		  getDeeplyNested: [[[String!]!]!]!
		  getMatrix: [[Int]]
		  getComplexOptional: [String!]
		  getRequiredListOfOptional: [String]!
		}
	`

	schema, _ := ParseSchema([]byte(schemaString))

	testCases := []struct {
		fieldName    string
		expectedType string
	}{
		{"getString", "String"},
		{"getRequiredString", "String!"},
		{"getStringList", "[String]"},
		{"getRequiredStringList", "[String]!"},
		{"getListOfRequiredStrings", "[String!]"},
		{"getRequiredListOfRequiredStrings", "[String!]!"},
		{"getDeeplyNested", "[[[String!]!]!]!"},
		{"getMatrix", "[[Int]]"},
		{"getComplexOptional", "[String!]"},
		{"getRequiredListOfOptional", "[String]!"},
	}

	for _, tc := range testCases {
		t.Run(tc.fieldName, func(t *testing.T) {
			field, ok := schema.Query[tc.fieldName]
			is.True(ok)
			is.Equal(field.TypeString(), tc.expectedType)
		})
	}
}

// Test that empty map returns empty slice
func TestCollectAndSortMapValuesEmpty(t *testing.T) {
	is := is.New(t)

	emptyMap := make(map[string]*FieldDefinition)
	result := CollectAndSortMapValues(emptyMap)
	is.Equal(len(result), 0)
}

// Test that fields are sorted alphabetically by name
func TestCollectAndSortMapValuesSorting(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type Query {
		  zebra: String
		  alpha: String
		  beta: String
		  gamma: String
		}
	`

	schema, _ := ParseSchema([]byte(schemaString))
	sortedFields := CollectAndSortMapValues(schema.Query)

	is.Equal(len(sortedFields), 4)
	is.Equal(sortedFields[0].Name(), "alpha")
	is.Equal(sortedFields[1].Name(), "beta")
	is.Equal(sortedFields[2].Name(), "gamma")
	is.Equal(sortedFields[3].Name(), "zebra")
}

func TestGetTypeNameAllTypes(t *testing.T) {
	is := is.New(t)

	schemaString := `
		type TestObject {
		  id: ID!
		}

		input TestInput {
		  name: String!
		}

		enum TestEnum {
		  VALUE_A
		  VALUE_B
		}

		scalar TestScalar

		interface TestInterface {
		  id: ID!
		}

		union TestUnion = TestObject

		directive @testDirective on FIELD_DEFINITION

		type Query {
		  testField: String
		}
	`
	schema, _ := ParseSchema([]byte(schemaString))

	testField := schema.Query["testField"]
	is.Equal(GetTypeName(testField), "testField")

	testObject := schema.Object["TestObject"]
	is.Equal(GetTypeName(testObject), "TestObject")

	testInput := schema.Input["TestInput"]
	is.Equal(GetTypeName(testInput), "TestInput")

	testEnum := schema.Enum["TestEnum"]
	is.Equal(GetTypeName(testEnum), "TestEnum")

	testScalar := schema.Scalar["TestScalar"]
	is.Equal(GetTypeName(testScalar), "TestScalar")

	testInterface := schema.Interface["TestInterface"]
	is.Equal(GetTypeName(testInterface), "TestInterface")

	testUnion := schema.Union["TestUnion"]
	is.Equal(GetTypeName(testUnion), "TestUnion")

	testDirective := schema.Directive["testDirective"]
	is.Equal(GetTypeName(testDirective), "testDirective")
}
