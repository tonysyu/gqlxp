package adapters

import (
	"strings"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/search"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
)

// Ensure searchResultItem implements components.ListItem interface
var _ components.ListItem = (*searchResultItem)(nil)

// searchResultItem wraps a resolved ListItem to preserve search result context
// in the display while showing parent context with field highlighting for field types
type searchResultItem struct {
	result       search.SearchResult
	wrappedItem  components.ListItem
	displayTitle string
	resolver     gql.TypeResolver
	parentType   string // Parent type name (e.g., "User" for "User.email")
	fieldName    string // Field name (e.g., "email" for "User.email")
}

func newSearchResultItem(
	result search.SearchResult,
	wrappedItem components.ListItem,
	resolver gql.TypeResolver,
	parentType, fieldName string,
) components.ListItem {
	// Create display title that shows the path/type context
	displayTitle := result.Path
	if displayTitle == "" {
		displayTitle = result.Name
	}

	return searchResultItem{
		result:       result,
		wrappedItem:  wrappedItem,
		displayTitle: displayTitle,
		resolver:     resolver,
		parentType:   parentType,
		fieldName:    fieldName,
	}
}

func (i searchResultItem) Title() string { return i.displayTitle }
func (i searchResultItem) FilterValue() string {
	// Use wrapped item's filter value for better searching
	return i.wrappedItem.FilterValue()
}
func (i searchResultItem) TypeName() string { return i.result.Type }
func (i searchResultItem) RefName() string {
	// For field types, use the full path (Type.field) for breadcrumbs
	switch i.result.Type {
	case "Query", "Mutation", "ObjectField", "InputField", "InterfaceField":
		if i.result.Path != "" {
			return i.result.Path
		}
	}
	// For non-field types, delegate to wrapped item
	return i.wrappedItem.RefName()
}
func (i searchResultItem) Description() string { return i.wrappedItem.Description() }
func (i searchResultItem) Details() string     { return i.wrappedItem.Details() }

// OpenPanel opens the parent type's panel with the field highlighted for field types,
// or delegates to the wrapped item for non-field types
func (i searchResultItem) OpenPanel() (*components.Panel, bool) {
	// For field types, show the parent type's panel with the field highlighted
	switch i.result.Type {
	case "ObjectField", "InputField", "InterfaceField":
		if i.resolver != nil && i.parentType != "" && i.fieldName != "" {
			typeDef, err := i.resolver.ResolveType(i.parentType)
			if err == nil {
				typeItem := newTypeDefItem(typeDef, i.resolver)
				panel, ok := typeItem.OpenPanel()
				if ok {
					panel.SelectItemByName(i.fieldName)
				}
				return panel, ok
			}
		}
	}
	// For non-field types (Object, Enum, Directive, etc.) or fallback
	return i.wrappedItem.OpenPanel()
}

// findFieldByName searches for a field by name in a slice of fields
func findFieldByName(fields []*gql.Field, name string) *gql.Field {
	for _, field := range fields {
		if field.Name() == name {
			return field
		}
	}
	return nil
}

// AdaptSearchResult converts a search.SearchResult to an appropriate ListItem
// based on the result type and path. Returns a fully-functional item with OpenPanel
// support when possible, falling back to SimpleItem if resolution fails.
func AdaptSearchResult(result search.SearchResult, schemaView *SchemaView) components.ListItem {
	// Try to resolve the search result to a proper item type
	switch result.Type {
	case "Query", "Mutation":
		// Format: "Query.fieldName" or "Mutation.fieldName"
		return adaptQueryOrMutationField(result, schemaView)

	case "ObjectField":
		// Format: "TypeName.fieldName"
		return adaptObjectField(result, schemaView)

	case "InputField":
		// Format: "TypeName.fieldName"
		return adaptInputField(result, schemaView)

	case "InterfaceField":
		// Format: "TypeName.fieldName"
		return adaptInterfaceField(result, schemaView)

	case "Object", "Input", "Enum", "Scalar", "Interface", "Union":
		// Format: just the type name
		return adaptType(result, schemaView)

	case "Directive":
		// Format: "@directiveName"
		return adaptDirective(result, schemaView)

	default:
		// Unknown type - fallback to SimpleItem
		return createFallbackItem(result)
	}
}

// adaptQueryOrMutationField handles Query and Mutation field results
func adaptQueryOrMutationField(result search.SearchResult, schemaView *SchemaView) components.ListItem {
	// Parse "Query.fieldName" or "Mutation.fieldName"
	parts := strings.SplitN(result.Path, ".", 2)
	if len(parts) != 2 {
		return createFallbackItem(result)
	}

	parentType := parts[0]
	fieldName := parts[1]
	var field *gql.Field

	if result.Type == "Query" {
		field = schemaView.schema.Query[fieldName]
	} else {
		field = schemaView.schema.Mutation[fieldName]
	}

	if field == nil {
		return createFallbackItem(result)
	}

	wrappedItem := newFieldItem(field, schemaView.resolver)
	return newSearchResultItem(result, wrappedItem, schemaView.resolver, parentType, fieldName)
}

// adaptObjectField handles ObjectField results
func adaptObjectField(result search.SearchResult, schemaView *SchemaView) components.ListItem {
	// Parse "TypeName.fieldName"
	parts := strings.SplitN(result.Path, ".", 2)
	if len(parts) != 2 {
		return createFallbackItem(result)
	}

	typeName := parts[0]
	fieldName := parts[1]

	obj := schemaView.schema.Object[typeName]
	if obj == nil {
		return createFallbackItem(result)
	}

	field := findFieldByName(obj.Fields(), fieldName)
	if field == nil {
		return createFallbackItem(result)
	}

	wrappedItem := newFieldItem(field, schemaView.resolver)
	return newSearchResultItem(result, wrappedItem, schemaView.resolver, typeName, fieldName)
}

// adaptInputField handles InputField results
func adaptInputField(result search.SearchResult, schemaView *SchemaView) components.ListItem {
	// Parse "TypeName.fieldName"
	parts := strings.SplitN(result.Path, ".", 2)
	if len(parts) != 2 {
		return createFallbackItem(result)
	}

	typeName := parts[0]
	fieldName := parts[1]

	input := schemaView.schema.Input[typeName]
	if input == nil {
		return createFallbackItem(result)
	}

	field := findFieldByName(input.Fields(), fieldName)
	if field == nil {
		return createFallbackItem(result)
	}

	wrappedItem := newFieldItem(field, schemaView.resolver)
	return newSearchResultItem(result, wrappedItem, schemaView.resolver, typeName, fieldName)
}

// adaptInterfaceField handles InterfaceField results
func adaptInterfaceField(result search.SearchResult, schemaView *SchemaView) components.ListItem {
	// Parse "TypeName.fieldName"
	parts := strings.SplitN(result.Path, ".", 2)
	if len(parts) != 2 {
		return createFallbackItem(result)
	}

	typeName := parts[0]
	fieldName := parts[1]

	iface := schemaView.schema.Interface[typeName]
	if iface == nil {
		return createFallbackItem(result)
	}

	field := findFieldByName(iface.Fields(), fieldName)
	if field == nil {
		return createFallbackItem(result)
	}

	wrappedItem := newFieldItem(field, schemaView.resolver)
	return newSearchResultItem(result, wrappedItem, schemaView.resolver, typeName, fieldName)
}

// adaptType handles type definition results (Object, Input, Enum, etc.)
func adaptType(result search.SearchResult, schemaView *SchemaView) components.ListItem {
	typeName := result.Path // For types, path is just the name

	// Try to resolve the type using the resolver
	typeDef, err := schemaView.resolver.ResolveType(typeName)
	if err != nil {
		return createFallbackItem(result)
	}

	wrappedItem := newTypeDefItem(typeDef, schemaView.resolver)
	return newSearchResultItem(result, wrappedItem, schemaView.resolver, "", "")
}

// adaptDirective handles Directive results
func adaptDirective(result search.SearchResult, schemaView *SchemaView) components.ListItem {
	// Parse "@directiveName"
	directiveName := strings.TrimPrefix(result.Path, "@")

	directive := schemaView.schema.Directive[directiveName]
	if directive == nil {
		return createFallbackItem(result)
	}

	wrappedItem := newDirectiveDefItem(directive, schemaView.resolver)
	return newSearchResultItem(result, wrappedItem, schemaView.resolver, "", "")
}

// createFallbackItem creates a SimpleItem when resolution fails
func createFallbackItem(result search.SearchResult) components.ListItem {
	title := result.Path
	if title == "" {
		title = result.Name
	}
	description := result.Type
	if result.Description != "" {
		description = result.Type + ": " + result.Description
	}

	return components.NewSimpleItem(
		title,
		components.WithDescription(description),
		components.WithTypeName(result.Name),
	)
}

// AdaptSearchResults converts a slice of search results to ListItems
func AdaptSearchResults(results []search.SearchResult, schemaView *SchemaView) []components.ListItem {
	items := make([]components.ListItem, 0, len(results))
	for _, result := range results {
		items = append(items, AdaptSearchResult(result, schemaView))
	}
	return items
}
