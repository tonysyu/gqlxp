package tui

import (
	"github.com/charmbracelet/bubbles/list"
)

// list item that can be "opened" to provide additional information about the item.
// The opened data is represented as a Panel instance that can be rendered to users.
type ListItem interface {
	list.DefaultItem

	// Open Panel to show additional information.
	Open() Panel
}

// Ensure that all item types implements ListItem interface
var _ ListItem = (*item)(nil)
