package tui

import (
	"github.com/charmbracelet/bubbles/list"
)

// List item that can be "expanded" to provide additional information about the item.
// The expanded data is represented as a Panel instance that can be rendered to users.
type ExpandableListItem interface {
	list.DefaultItem

	Expand() Panel
}
