package adapters

import (
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/gqlfmt"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
)

// Adapter/delegate for gql.FieldDefinition to support ListItem interface
type fieldItem struct {
	gqlField  *gql.Field
	resolver  gql.TypeResolver
	fieldName string
}

func newFieldItem(gqlField *gql.Field, resolver gql.TypeResolver) components.ListItem {
	return fieldItem{
		gqlField:  gqlField,
		resolver:  resolver,
		fieldName: gqlField.Name(),
	}
}

func (i fieldItem) Title() string       { return i.gqlField.Signature() }
func (i fieldItem) FilterValue() string { return i.fieldName }
func (i fieldItem) TypeName() string    { return i.gqlField.ObjectTypeName() }
func (i fieldItem) RefName() string     { return i.gqlField.Name() }

func (i fieldItem) Description() string {
	return i.gqlField.Description()
}

func (i fieldItem) Details() string {
	return gqlfmt.GenerateFieldMarkdown(i.gqlField, i.resolver)
}

// OpenPanel displays arguments of field (if any) and the field's ObjectType
func (i fieldItem) OpenPanel() (*components.Panel, bool) {
	argumentItems := adaptArguments(i.gqlField.Arguments(), i.resolver)
	resultTypeItem := newTypeDefItemFromField(i.gqlField, i.resolver)
	directiveItems := adaptAppliedDirectives(i.gqlField.Directives(), i.resolver)

	panel := components.NewPanel([]components.ListItem{}, i.fieldName)
	panel.SetDescription(i.Description())

	// Create tabs for Result Type and Input Arguments
	var tabs []components.Tab
	tabs = append(tabs, components.Tab{
		Label:   "Type",
		Content: []components.ListItem{resultTypeItem},
	})
	if len(argumentItems) > 0 {
		tabs = append(tabs, components.Tab{
			Label:   "Inputs",
			Content: argumentItems,
		})
	}
	if len(directiveItems) > 0 {
		tabs = append(tabs, components.Tab{
			Label:   "Directives",
			Content: directiveItems,
		})
	}
	panel.SetTabs(tabs)

	return panel, true
}
