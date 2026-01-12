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
	resultTypeItem := newTypeDefItemFromField(i.gqlField, i.resolver)

	panel := components.NewPanel([]components.ListItem{}, i.fieldName)
	panel.SetDescription(i.Description())

	// Create tabs for Result Type and Input Arguments
	var tabs []components.Tab
	tabs = append(tabs, newTypeTab(resultTypeItem))
	if len(i.gqlField.Arguments()) > 0 {
		tabs = append(tabs, newInputsTab(i.gqlField.Arguments(), i.resolver))
	}
	if len(i.gqlField.Directives()) > 0 {
		tabs = append(tabs, newDirectivesTab(i.gqlField.Directives(), i.resolver))
	}
	panel.SetTabs(tabs)

	return panel, true
}
