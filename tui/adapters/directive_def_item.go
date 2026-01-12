package adapters

import (
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/gqlfmt"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
)

// Adapter/delegate for gql.DirectiveDef to support ListItem interface (for schema directives)
type directiveDefItem struct {
	gqlDirective  *gql.DirectiveDef
	resolver      gql.TypeResolver
	directiveName string
}

func newDirectiveDefItem(directive *gql.DirectiveDef, resolver gql.TypeResolver) components.ListItem {
	return directiveDefItem{
		gqlDirective:  directive,
		resolver:      resolver,
		directiveName: directive.Name(),
	}
}

func (i directiveDefItem) Title() string       { return i.gqlDirective.Signature() }
func (i directiveDefItem) FilterValue() string { return i.directiveName }
func (i directiveDefItem) TypeName() string    { return "@" + i.directiveName }
func (i directiveDefItem) RefName() string     { return i.directiveName }

func (i directiveDefItem) Description() string {
	return i.gqlDirective.Description()
}

func (i directiveDefItem) Details() string {
	return gqlfmt.GenerateDirectiveMarkdown(i.gqlDirective, i.resolver)
}

// OpenPanel displays arguments of directive (if any)
func (i directiveDefItem) OpenPanel() (*components.Panel, bool) {
	panel := components.NewPanel([]components.ListItem{}, "@"+i.directiveName)
	panel.SetDescription(i.Description())

	var tabs []components.Tab
	if len(i.gqlDirective.Arguments()) > 0 {
		tabs = append(tabs, newArgumentsTab(i.gqlDirective.Arguments(), i.resolver))
	}

	// Add Usages tab if the directive is used elsewhere
	if usages, _ := i.resolver.ResolveUsages(i.directiveName); len(usages) > 0 {
		tabs = append(tabs, newUsagesTab(usages, i.resolver))
	}

	if len(tabs) > 0 {
		panel.SetTabs(tabs)
	}

	return panel, true
}
