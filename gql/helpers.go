package gql

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

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

// NamedTypeDef is a custom ast.TypeDef with GetName() method.
//
// For some reason graphql-go defines ast.TypeDefinition without GetName() but all
// implementers should have this method.
type NamedTypeDef interface {
	ast.TypeDefinition
	GetName() *ast.Name
}

// NamedType interface for types that have a Name field (both ast and wrapped types)
type NamedType interface {
	*ast.FieldDefinition | *ast.ObjectDefinition | *ast.InputObjectDefinition |
		*ast.EnumDefinition | *ast.ScalarDefinition | *ast.InterfaceDefinition |
		*ast.UnionDefinition | *ast.DirectiveDefinition |
		*FieldDefinition | *ObjectDefinition | *InputObjectDefinition |
		*EnumDefinition | *ScalarDefinition | *InterfaceDefinition |
		*UnionDefinition | *DirectiveDefinition
}

// GetTypeString converts ast.Type to string representation
// ast.Types are ast.Named types wrapped in arbitrary numbers of lists and non-nulls.
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

// GetTypeString converts ast.Type to string representation
// ast.Types are ast.Named types wrapped in arbitrary numbers of lists and non-nulls.
func GetNamedFromType(t ast.Type) *ast.Named {
	switch typ := t.(type) {
	case *ast.Named:
		return typ
	case *ast.List:
		return GetNamedFromType(typ.Type)
	case *ast.NonNull:
		return GetNamedFromType(typ.Type)
	default:
		return nil
	}
}

// GetTypeName extracts the name from various AST node types and wrapped types.
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
	// Wrapped types all have Name() method
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

func GetInputValueDefinitionString(inputValue *ast.InputValueDefinition) string {
	fieldName := inputValue.Name.Value
	fieldType := GetTypeString(inputValue.Type)
	return fmt.Sprintf("%s: %s", fieldName, fieldType)
}

// Return string representing the `<field>: <type>` pair or signature of a field.
func GetFieldDefinitionString(field *ast.FieldDefinition) string {
	fieldName := field.Name.Value
	fieldType := GetTypeString(field.Type)
	if len(field.Arguments) > 0 {
		var inputArgs []string
		for _, arg := range field.Arguments {
			inputArgs = append(inputArgs, GetInputValueDefinitionString(arg))
		}
		inputArgString := strings.Join(inputArgs, ", ")
		return fieldName + "(" + inputArgString + "): " + fieldType
	}
	return fieldName + ": " + fieldType
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
