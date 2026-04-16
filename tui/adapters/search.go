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
func (i searchResultItem) TypeName() string { return i.result.Kind }
func (i searchResultItem) RefName() string {
	// For field types, use the full path (Type.field) for breadcrumbs
	switch i.result.Kind {
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
	switch i.result.Kind {
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

// AdaptSearchResult converts a search.SearchResult to an appropriate ListItem
// based on the result type and path. Returns a fully-functional item with OpenPanel
// support when possible, falling back to SimpleItem if resolution fails.
func AdaptSearchResult(result search.SearchResult, schemaView *SchemaView) components.ListItem {
	switch result.Kind {
	case "Query", "Mutation", "ObjectField", "InputField", "InterfaceField":
		return adaptField(result, schemaView)
	case "Object", "Input", "Enum", "Scalar", "Interface", "Union":
		return adaptType(result, schemaView)
	case "Directive":
		return adaptDirective(result, schemaView)
	default:
		return createFallbackItem(result)
	}
}

// adaptField handles all field-bearing search result kinds.
// Path format: "ParentType.fieldName" (e.g. "Query.user", "User.email")
func adaptField(result search.SearchResult, schemaView *SchemaView) components.ListItem {
	parts := strings.SplitN(result.Path, ".", 2)
	if len(parts) != 2 {
		return createFallbackItem(result)
	}
	typeName, fieldName := parts[0], parts[1]

	field, err := schemaView.ResolveField(result.Kind, typeName, fieldName)
	if err != nil {
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
	description := result.Kind
	if result.Description != "" {
		description = result.Kind + ": " + result.Description
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
