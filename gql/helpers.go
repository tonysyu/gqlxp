package gql

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// getTypeString converts ast.Type to string representation
// In gqlparser, Type is a struct with NamedType, Elem (for lists), and NonNull fields.
func getTypeString(t *ast.Type) string {
	if t == nil {
		return "Unknown"
	}

	// Build the type string recursively
	var result string
	if t.NamedType != "" {
		result = t.NamedType
	} else if t.Elem != nil {
		result = "[" + getTypeString(t.Elem) + "]"
	} else {
		return "Unknown"
	}

	if t.NonNull {
		result += "!"
	}

	return result
}

// getNamedTypeName extracts the base named type from an ast.Type
// unwrapping any List or NonNull wrappers.
func getNamedTypeName(t *ast.Type) string {
	if t == nil {
		return ""
	}

	if t.NamedType != "" {
		return t.NamedType
	}

	if t.Elem != nil {
		return getNamedTypeName(t.Elem)
	}

	return ""
}

// getTypeName extracts the name from various AST node types and wrapped types.
func getTypeName[T gqlType](node T) string {
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

// getArgumentString returns a string representation of an argument definition
func getArgumentString(arg *ast.ArgumentDefinition) string {
	return fmt.Sprintf("%s: %s", arg.Name, getTypeString(arg.Type))
}

// getInputFieldString returns a string representation of an input field
func getInputFieldString(field *ast.FieldDefinition) string {
	return fmt.Sprintf("%s: %s", field.Name, getTypeString(field.Type))
}

// getFieldString returns the signature of a field including arguments
func getFieldString(field *ast.FieldDefinition) string {
	if len(field.Arguments) > 0 {
		var inputArgs []string
		for _, arg := range field.Arguments {
			inputArgs = append(inputArgs, getArgumentString(arg))
		}
		inputArgString := strings.Join(inputArgs, ", ")
		return fmt.Sprintf("%s(%s): %s", field.Name, inputArgString, getTypeString(field.Type))
	}
	return fmt.Sprintf("%s: %s", field.Name, getTypeString(field.Type))
}

// CollectAndSortMapValues extracts values from a map, sorts them by name, and returns a slice.
// GraphQL types in GraphQLSchema are stored as `name` -> `ast.TypeDefinition`/`ast.FieldDefinition`
// maps. This helper extracts and sorts the values by name.
func CollectAndSortMapValues[T gqlType](m map[string]T) []T {
	values := slices.Collect(maps.Values(m))
	sort.Slice(values, func(i, j int) bool {
		return getTypeName(values[i]) < getTypeName(values[j])
	})
	return values
}
