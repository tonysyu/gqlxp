package components

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/tonysyu/gqlxp/utils/text"
)

// ListItem is a list item that can be "opened" to provide additional information about the item.
// The opened data is represented as a Panel instance that can be rendered to users.
type ListItem interface {
	list.DefaultItem

	// OpenPanel Panel to show child items, field types, and inputs.
	OpenPanel() (Panel, bool)

	// Details returns markdown-formatted details for the item.
	// The gqlxp TUI renders details in an overlay pane.
	Details() string

	// TypeName returns the name of the underlying GraphQL type.
	// This often matches the Title() but may differ for types wrapped in lists and non-nulls.
	// Title() is used when referencing this item in lists
	// In contrast, TypeName() is used as the title by Details()
	TypeName() string

	// RefName returns the reference name for the item.
	// For most types, this is equivalent to Title(), but for fields, this leaves off the type and
	// just has the field name (with no arguments)
	RefName() string
}

var _ ListItem = (*SimpleItem)(nil)

// SimpleItem is a ListItem implementation with arbitrary title and description and no-op Open() function.
type SimpleItem struct {
	title       string
	description string
	typename    string
}

// SimpleItemOption is a function that configures a SimpleItem.
type SimpleItemOption func(*SimpleItem)

// WithDescription sets the description for a SimpleItem.
func WithDescription(desc string) SimpleItemOption {
	return func(si *SimpleItem) {
		si.description = desc
	}
}

// WithTypeName sets the typename for a SimpleItem.
func WithTypeName(typename string) SimpleItemOption {
	return func(si *SimpleItem) {
		si.typename = typename
	}
}

// NewSimpleItem creates a new SimpleItem with the given title and optional configuration.
func NewSimpleItem(title string, opts ...SimpleItemOption) SimpleItem {
	si := SimpleItem{
		title:    title,
		typename: title,
	}
	for _, opt := range opts {
		opt(&si)
	}
	return si
}

func (si SimpleItem) Title() string       { return si.title }
func (si SimpleItem) Description() string { return si.description }
func (si SimpleItem) FilterValue() string { return si.Title() }
func (si SimpleItem) TypeName() string    { return si.typename }
func (si SimpleItem) RefName() string     { return si.typename }
func (si SimpleItem) Details() string {
	if si.Description() == "" {
		return ""
	}
	return text.JoinParagraphs(
		text.H1(si.TypeName()),
		si.Description(),
	)
}
func (si SimpleItem) OpenPanel() (Panel, bool) { return nil, false }
