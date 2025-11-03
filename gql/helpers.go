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

// formatValue converts an AST value to its string representation
func formatValue(value *ast.Value) string {
	if value == nil {
		return "null"
	}

	switch value.Kind {
	case ast.IntValue, ast.FloatValue, ast.BooleanValue, ast.EnumValue:
		return value.Raw
	case ast.StringValue:
		return fmt.Sprintf(`"%s"`, value.Raw)
	case ast.NullValue:
		return "null"
	case ast.ListValue:
		if len(value.Children) == 0 {
			return "[]"
		}
		var values []string
		for _, child := range value.Children {
			values = append(values, formatValue(child.Value))
		}
		return "[" + strings.Join(values, ", ") + "]"
	case ast.ObjectValue:
		if len(value.Children) == 0 {
			return "{}"
		}
		var pairs []string
		for _, child := range value.Children {
			pairs = append(pairs, fmt.Sprintf("%s: %s", child.Name, formatValue(child.Value)))
		}
		return "{" + strings.Join(pairs, ", ") + "}"
	default:
		return value.Raw
	}
}

// getArgumentString returns a string representation of an argument definition
func getArgumentString(arg *ast.ArgumentDefinition) string {
	result := fmt.Sprintf("%s: %s", arg.Name, getTypeString(arg.Type))
	if arg.DefaultValue != nil {
		result += fmt.Sprintf(" = %s", formatValue(arg.DefaultValue))
	}
	return result
}

// getFieldString returns the signature of a field including arguments and default values
// Handles both output fields (with arguments) and input fields (with default values)
func getFieldString(field *ast.FieldDefinition) string {
	return formatFieldStringWithWidth(field, 0)
}

// formatFieldStringWithWidth returns the field signature with optional multiline formatting
// If maxWidth <= 0, always uses inline format
// If maxWidth > 0 and inline signature exceeds maxWidth, formats arguments on separate lines
func formatFieldStringWithWidth(field *ast.FieldDefinition, maxWidth int) string {
	var result string

	// Format with arguments if present
	if len(field.Arguments) > 0 {
		var inputArgs []string
		for _, arg := range field.Arguments {
			inputArgs = append(inputArgs, getArgumentString(arg))
		}

		// Try inline format first
		inputArgString := strings.Join(inputArgs, ", ")
		inlineResult := fmt.Sprintf("%s(%s): %s", field.Name, inputArgString, getTypeString(field.Type))

		// Check if we should use multiline format
		if maxWidth > 0 && len(inlineResult) > maxWidth {
			// Format with each argument on a new line, indented
			multilineArgs := strings.Join(inputArgs, ",\n  ")
			result = fmt.Sprintf("%s(\n  %s\n): %s", field.Name, multilineArgs, getTypeString(field.Type))
		} else {
			result = inlineResult
		}
	} else {
		result = fmt.Sprintf("%s: %s", field.Name, getTypeString(field.Type))
	}

	// Add default value if present (for input fields on InputObjects)
	if field.DefaultValue != nil {
		result += fmt.Sprintf(" = %s", formatValue(field.DefaultValue))
	}

	return result
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
