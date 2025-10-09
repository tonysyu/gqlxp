package components

import (
	"github.com/charmbracelet/bubbles/list"
)

// ListItem is a list item that can be "opened" to provide additional information about the item.
// The opened data is represented as a Panel instance that can be rendered to users.
type ListItem interface {
	list.DefaultItem

	// Open Panel to show additional information.
	Open() (Panel, bool)

	// Details returns markdown-formatted details for the item.
	Details() string
}

// SimpleItem is a ListItem implementation with arbitrary title and description and no-op Open() function.
type SimpleItem struct {
	title       string
	description string
}

// NewSimpleItem creates a new SimpleItem with the given title and description
func NewSimpleItem(title, description string) SimpleItem {
	return SimpleItem{
		title:       title,
		description: description,
	}
}

func (si SimpleItem) Title() string       { return si.title }
func (si SimpleItem) Description() string { return si.description }
func (si SimpleItem) FilterValue() string { return si.title }
func (si SimpleItem) Details() string {
	if si.description != "" {
		return "# " + si.Title() + "\n\n" + si.Description()
	}
	return "# " + si.Title()
}
func (si SimpleItem) Open() (Panel, bool) { return nil, false }
