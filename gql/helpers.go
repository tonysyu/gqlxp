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
		return ""
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
func getArgumentString(arg *Argument) string {
	result := fmt.Sprintf("%s: %s", arg.Name(), arg.TypeString())
	if defaultValue := arg.DefaultValue(); defaultValue != "" {
		result = fmt.Sprintf("%s = %s", result, defaultValue)
	}
	return result
}

func getArgumentStringList(argList []*Argument) []string {
	var argStringList []string
	for _, arg := range argList {
		argStringList = append(argStringList, getArgumentString(arg))
	}
	return argStringList
}

type formatOpts struct {
	prefix   string
	suffix   string
	maxWidth int
}

// formatCallable returns the field signature with optional multiline formatting
// - If maxWidth <= 0, always uses inline format
// - If maxWidth > 0 and inline signature exceeds maxWidth, formats arguments on separate lines
//
// Suffix input is appended to the end of callable
// - For GQL fields, this will be ": <field-type>""
// - For everything else, this will be an empty string
func formatCallable(callable Callable, opts formatOpts) string {
	var result string

	name := opts.prefix + callable.Name()

	// If no arguments just return callable + suffix
	if len(callable.Arguments()) == 0 {
		return name + opts.suffix
	}

	inputArgs := getArgumentStringList(callable.Arguments())

	// Try inline format first
	inputArgString := strings.Join(inputArgs, ", ")
	inlineResult := fmt.Sprintf("%s(%s)%s", name, inputArgString, opts.suffix)

	// If inlineResult fits within maxWidth, just return inlineResult
	if opts.maxWidth <= 0 || len(inlineResult) <= opts.maxWidth {
		return inlineResult
	}

	// Format with each argument on a new line, indented
	multilineArgs := strings.Join(inputArgs, ",\n  ")
	result = fmt.Sprintf("%s(\n  %s\n)%s", name, multilineArgs, opts.suffix)
	return result
}

// getFieldString returns the signature of a field including arguments and default values
// Handles both output fields (with arguments) and input fields (with default values)
func getFieldString(field *Field) string {
	return formatFieldStringWithWidth(field, 0)
}

// formatFieldStringWithWidth returns the field signature with optional multiline formatting
// If maxWidth <= 0, always uses inline format
// If maxWidth > 0 and inline signature exceeds maxWidth, formats arguments on separate lines
func formatFieldStringWithWidth(field *Field, maxWidth int) string {
	result := formatCallable(
		field,
		formatOpts{
			suffix:   fmt.Sprintf(": %s", field.TypeString()),
			maxWidth: maxWidth,
		},
	)

	// Add default value if present (for input fields on InputObjects)
	if defaultValue := field.DefaultValue(); defaultValue != "" {
		result = fmt.Sprintf("%s = %s", result, defaultValue)
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
