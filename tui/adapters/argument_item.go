package adapters

import (
	"github.com/tonysyu/gqlxp/gql"
	"github.com/tonysyu/gqlxp/tui/xplr/components"
	"github.com/tonysyu/gqlxp/utils/text"
)

// Adapter/delegate for gql.Argument to support ListItem interface
type argumentItem struct {
	gqlArgument *gql.Argument
	resolver    gql.TypeResolver
	argName     string
}

func newArgumentItem(gqlArgument *gql.Argument, resolver gql.TypeResolver) components.ListItem {
	return argumentItem{
		gqlArgument: gqlArgument,
		resolver:    resolver,
		argName:     gqlArgument.Name(),
	}
}

func (i argumentItem) Title() string       { return i.gqlArgument.Signature() }
func (i argumentItem) FilterValue() string { return i.argName }
func (i argumentItem) TypeName() string    { return i.gqlArgument.ObjectTypeName() }
func (i argumentItem) RefName() string     { return i.gqlArgument.Name() }

func (i argumentItem) Description() string {
	return i.gqlArgument.Description()
}

func (i argumentItem) Details() string {
	return text.JoinParagraphs(
		text.H1(i.argName),
		text.GqlCode(i.gqlArgument.FormatSignature(80)),
		i.Description(),
	)
}

// OpenPanel displays the argument's type definition
func (i argumentItem) OpenPanel() (*components.Panel, bool) {
	resultTypeItem := newTypeDefItemFromArgument(i.gqlArgument, i.resolver)

	panel := components.NewPanel([]components.ListItem{}, i.argName)
	panel.SetDescription(i.Description())

	// Create a single tab for Result Type
	panel.SetTabs([]components.Tab{
		{
			Label:   "Type",
			Content: []components.ListItem{resultTypeItem},
		},
	})

	return panel, true
}
