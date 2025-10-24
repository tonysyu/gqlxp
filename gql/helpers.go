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

var _ NamedTypeDef = (*Field)(nil)
var _ NamedTypeDef = (*Object)(nil)
var _ NamedTypeDef = (*InputObject)(nil)
var _ NamedTypeDef = (*Enum)(nil)
var _ NamedTypeDef = (*Scalar)(nil)
var _ NamedTypeDef = (*Interface)(nil)
var _ NamedTypeDef = (*Union)(nil)
var _ NamedTypeDef = (*Directive)(nil)

// namedType interface for types that have a Name field (both ast and wrapped types)
type namedType interface {
		*Field | *Object | *InputObject |
		*Enum | *Scalar | *Interface |
		*Union | *Directive
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
	case *Field:
		return n.Name()
	case *Object:
		return n.Name()
	case *InputObject:
		return n.Name()
	case *Enum:
		return n.Name()
	case *Scalar:
		return n.Name()
	case *Interface:
		return n.Name()
	case *Union:
		return n.Name()
	case *Directive:
		return n.Name()
	default:
		return ""
	}
}

func getInputValueString(inputValue *ast.InputValueDefinition) string {
	fieldName := inputValue.Name.Value
	fieldType := getTypeString(inputValue.Type)
	return fmt.Sprintf("%s: %s", fieldName, fieldType)
}

// Return string representing the `<field>: <type>` pair or signature of a field.
func getFieldString(field *ast.FieldDefinition) string {
	fieldName := field.Name.Value
	fieldType := getTypeString(field.Type)
	if len(field.Arguments) > 0 {
		var inputArgs []string
		for _, arg := range field.Arguments {
			inputArgs = append(inputArgs, getInputValueString(arg))
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
