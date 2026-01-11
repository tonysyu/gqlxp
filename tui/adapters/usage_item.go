package adapters

import (
	"fmt"

	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
)

type usageItem struct {
	usage    *gql.Usage
	resolver gql.TypeResolver
}

func newUsageItem(usage *gql.Usage, resolver gql.TypeResolver) components.ListItem {
	return &usageItem{usage: usage, resolver: resolver}
}

// Title returns the usage path (e.g., "Query.user")
func (i *usageItem) Title() string {
	return i.usage.Path
}

// Description returns context (e.g., "Query")
func (i *usageItem) Description() string {
	return i.usage.ParentKind
}

// FilterValue returns the value used for filtering
func (i *usageItem) FilterValue() string {
	return i.usage.Path
}

// TypeName returns the parent type name for navigation
func (i *usageItem) TypeName() string {
	return i.usage.ParentType
}

// RefName returns the reference name (clean name without type info)
func (i *usageItem) RefName() string {
	return i.usage.ParentType
}

// OpenPanel creates a panel showing the parent type with focus on the field
func (i *usageItem) OpenPanel() (*components.Panel, bool) {
	// Handle Query and Mutation types specially since they're not TypeDefs
	if i.usage.ParentType == "Query" || i.usage.ParentType == "Mutation" {
		field, err := i.resolver.ResolveQueryOrMutationField(i.usage.ParentType, i.usage.FieldName)
		if err != nil {
			// Fallback to simple info panel
			panel := components.NewPanel([]components.ListItem{}, i.usage.Path)
			panel.SetDescription(fmt.Sprintf("Used in %s.%s (%s)",
				i.usage.ParentType, i.usage.FieldName, i.usage.ParentKind))
			return panel, true
		}
		// Create a fieldItem and open its panel
		fieldItem := newFieldItem(field, i.resolver)
		return fieldItem.OpenPanel()
	}

	// For regular types (Object, Interface, etc.), show the type's panel
	typeDef, err := i.resolver.ResolveType(i.usage.ParentType)
	if err != nil {
		// Fallback to simple info panel
		panel := components.NewPanel([]components.ListItem{}, i.usage.Path)
		panel.SetDescription(fmt.Sprintf("Used in %s.%s (%s)",
			i.usage.ParentType, i.usage.FieldName, i.usage.ParentKind))
		return panel, true
	}

	// Create the type's panel and select the specific field
	typeItem := newTypeDefItem(typeDef, i.resolver)
	panel, ok := typeItem.OpenPanel()
	if ok {
		// Try to select the field by name
		panel.SelectItemByName(i.usage.FieldName)
	}
	return panel, ok
}

// Details returns markdown showing usage context
func (i *usageItem) Details() string {
	return fmt.Sprintf("# %s\n\n**Field:** %s in %s (%s)",
		i.usage.Path, i.usage.FieldName, i.usage.ParentType, i.usage.ParentKind)
}
