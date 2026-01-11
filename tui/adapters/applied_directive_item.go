package adapters

import (
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/utils/text"
)

// Adapter/delegate for gql.AppliedDirective to support ListItem interface (for directives on fields/types)
type appliedDirectiveItem struct {
	gqlDirective  *gql.AppliedDirective
	resolver      gql.TypeResolver
	directiveName string
}

func newAppliedDirectiveItem(directive *gql.AppliedDirective, resolver gql.TypeResolver) components.ListItem {
	return appliedDirectiveItem{
		gqlDirective:  directive,
		resolver:      resolver,
		directiveName: directive.Name(),
	}
}

func (i appliedDirectiveItem) Title() string       { return i.gqlDirective.Signature() }
func (i appliedDirectiveItem) FilterValue() string { return i.directiveName }
func (i appliedDirectiveItem) TypeName() string    { return "@" + i.directiveName }
func (i appliedDirectiveItem) RefName() string     { return i.directiveName }

func (i appliedDirectiveItem) Description() string {
	return ""
}

func (i appliedDirectiveItem) Details() string {
	return text.JoinParagraphs(
		text.H1("@"+i.directiveName),
		text.GqlCode(i.gqlDirective.Signature()),
	)
}

// OpenPanel displays arguments from the directive definition
func (i appliedDirectiveItem) OpenPanel() (*components.Panel, bool) {
	// Resolve the directive name to its definition
	directiveDef, err := i.resolver.ResolveDirective(i.directiveName)
	if err != nil {
		// If we can't resolve the directive, return no panel
		return nil, false
	}

	// Get the arguments from the directive definition
	argumentItems := adaptArguments(directiveDef.Arguments(), i.resolver)

	panel := components.NewPanel(argumentItems, "@"+i.directiveName)
	panel.SetDescription(directiveDef.Description())

	return panel, true
}
