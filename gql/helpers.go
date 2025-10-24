package gql

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
)

// getStringValue converts ast.StringValue to string representation
// Note that StringValue pointers are nullable, so this avoids
func getStringValue(s *ast.StringValue) string {
	if s == nil {
		return ""
	}
	return s.Value
}

type NamedTypeDef interface {
	Name() string
	Description() string
}

var _ NamedTypeDef = (*FieldDefinition)(nil)
var _ NamedTypeDef = (*ObjectDefinition)(nil)
var _ NamedTypeDef = (*InputObjectDefinition)(nil)
var _ NamedTypeDef = (*EnumDefinition)(nil)
var _ NamedTypeDef = (*ScalarDefinition)(nil)
var _ NamedTypeDef = (*InterfaceDefinition)(nil)
var _ NamedTypeDef = (*UnionDefinition)(nil)
var _ NamedTypeDef = (*DirectiveDefinition)(nil)

// namedType interface for types that have a Name field (both ast and wrapped types)
type namedType interface {
		*FieldDefinition | *ObjectDefinition | *InputObjectDefinition |
		*EnumDefinition | *ScalarDefinition | *InterfaceDefinition |
		*UnionDefinition | *DirectiveDefinition
}

// getTypeString converts ast.Type to string representation
// ast.Types are ast.Named types wrapped in arbitrary numbers of lists and non-nulls.
func getTypeString(t ast.Type) string {
	switch typ := t.(type) {
	case *ast.Named:
		return typ.Name.Value
	case *ast.List:
		return "[" + getTypeString(typ.Type) + "]"
	case *ast.NonNull:
		return getTypeString(typ.Type) + "!"
	default:
		return "Unknown"
	}
}

// getNamedFromType converts ast.Type to string representation
// ast.Types are ast.Named types wrapped in arbitrary numbers of lists and non-nulls.
func getNamedFromType(t ast.Type) *ast.Named {
	switch typ := t.(type) {
	case *ast.Named:
		return typ
	case *ast.List:
		return getNamedFromType(typ.Type)
	case *ast.NonNull:
		return getNamedFromType(typ.Type)
	default:
		return nil
	}
}

// getTypeName extracts the name from various AST node types and wrapped types.
func getTypeName[T namedType](node T) string {
	// All these GraphQL types have `Name` attributes, but this isn't exposed in any shared
	// interface, so we make due with this silly switch statement.
	switch n := any(node).(type) {
	case *FieldDefinition:
		return n.Name()
	case *ObjectDefinition:
		return n.Name()
	case *InputObjectDefinition:
		return n.Name()
	case *EnumDefinition:
		return n.Name()
	case *ScalarDefinition:
		return n.Name()
	case *InterfaceDefinition:
		return n.Name()
	case *UnionDefinition:
		return n.Name()
	case *DirectiveDefinition:
		return n.Name()
	default:
		return ""
	}
}

func getInputValueDefinitionString(inputValue *ast.InputValueDefinition) string {
	fieldName := inputValue.Name.Value
	fieldType := getTypeString(inputValue.Type)
	return fmt.Sprintf("%s: %s", fieldName, fieldType)
}

// Return string representing the `<field>: <type>` pair or signature of a field.
func getFieldDefinitionString(field *ast.FieldDefinition) string {
	fieldName := field.Name.Value
	fieldType := getTypeString(field.Type)
	if len(field.Arguments) > 0 {
		var inputArgs []string
		for _, arg := range field.Arguments {
			inputArgs = append(inputArgs, getInputValueDefinitionString(arg))
		}
		inputArgString := strings.Join(inputArgs, ", ")
		return fieldName + "(" + inputArgString + "): " + fieldType
	}
	return fieldName + ": " + fieldType
}

// CollectAndSortMapValues extracts values from a map, sorts them by name, and returns a slice.
// GraphQL types in GraphQLSchema are stored as `name` -> `ast.TypeDefinition`/`ast.FieldDefinition`
// maps. This helper extracts and sorts the values by name.
func CollectAndSortMapValues[T namedType](m map[string]T) []T {
	values := slices.Collect(maps.Values(m))
	sort.Slice(values, func(i, j int) bool {
		return getTypeName(values[i]) < getTypeName(values[j])
	})
	return values
}
