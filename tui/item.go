package tui

import (
	"github.com/charmbracelet/bubbles/list"
)

// Interactive list item that can be "opened" to provide additional information about the item.
// The opened data is represented as a Panel instance that can be rendered to users.
type InteractiveListItem interface {
	list.DefaultItem

	// Open Panel to show additional information.
	Open() Panel
}
