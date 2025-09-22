package gql

import (
	"maps"
	"slices"
	"sort"

	"github.com/graphql-go/graphql/language/ast"
)

// GetStringValue converts ast.StringValue to string representation
// Note that StringValue pointers are nullable, so this avoids
func GetStringValue(s *ast.StringValue) string {
	if s == nil {
		return ""
	}
	return s.Value
}

// NamedType interface for types that have a Name field
type NamedType interface {
	*ast.FieldDefinition | *ast.ObjectDefinition | *ast.InputObjectDefinition |
		*ast.EnumDefinition | *ast.ScalarDefinition | *ast.InterfaceDefinition |
		*ast.UnionDefinition | *ast.DirectiveDefinition
}

// GetTypeString converts ast.Type to string representation
func GetTypeString(t ast.Type) string {
	switch typ := t.(type) {
	case *ast.Named:
		return typ.Name.Value
	case *ast.List:
		return "[" + GetTypeString(typ.Type) + "]"
	case *ast.NonNull:
		return GetTypeString(typ.Type) + "!"
	default:
		return "Unknown"
	}
}

// GetTypeName extracts the name from various AST node types.
func GetTypeName[T NamedType](node T) string {
	// All these GraphQL types have `Name` attributes, but this isn't exposed in any shared
	// interface, so we make due with this silly switch statement.
	switch n := any(node).(type) {
	case *ast.FieldDefinition:
		return n.Name.Value
	case *ast.ObjectDefinition:
		return n.Name.Value
	case *ast.InputObjectDefinition:
		return n.Name.Value
	case *ast.EnumDefinition:
		return n.Name.Value
	case *ast.ScalarDefinition:
		return n.Name.Value
	case *ast.InterfaceDefinition:
		return n.Name.Value
	case *ast.UnionDefinition:
		return n.Name.Value
	case *ast.DirectiveDefinition:
		return n.Name.Value
	default:
		return ""
	}
}

// CollectAndSortMapValues extracts values from a map, sorts them by name, and returns a slice.
// GraphQL types in GraphQLSchema are stored as `name` -> `ast.TypeDefinition`/`ast.FieldDefinition`
// maps. This helper extracts and sorts the values by name.
func CollectAndSortMapValues[T NamedType](m map[string]T) []T {
	values := slices.Collect(maps.Values(m))
	sort.Slice(values, func(i, j int) bool {
		return GetTypeName(values[i]) < GetTypeName(values[j])
	})
	return values
}
