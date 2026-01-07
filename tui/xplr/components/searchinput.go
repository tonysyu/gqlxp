package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// SearchInput is a Bubble Tea component for search input
type SearchInput struct {
	textInput textinput.Model
}

// NewSearchInput creates a new search input component
func NewSearchInput() SearchInput {
	ti := textinput.New()
	ti.Placeholder = "Type to search schema..."
	ti.CharLimit = 100
	ti.Width = 50

	return SearchInput{
		textInput: ti,
	}
}

// Init initializes the search input component
func (si SearchInput) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the search input
func (si SearchInput) Update(msg tea.Msg) (SearchInput, tea.Cmd) {
	var cmd tea.Cmd
	si.textInput, cmd = si.textInput.Update(msg)
	return si, cmd
}

// View renders the search input
func (si SearchInput) View() string {
	return si.textInput.View()
}

// Focus focuses the input field
func (si *SearchInput) Focus() tea.Cmd {
	return si.textInput.Focus()
}

// Blur removes focus from the input field
func (si *SearchInput) Blur() {
	si.textInput.Blur()
}

// Value returns the current input value
func (si SearchInput) Value() string {
	return si.textInput.Value()
}

// SetValue sets the input value
func (si *SearchInput) SetValue(value string) {
	si.textInput.SetValue(value)
}

// Focused returns whether the input is focused
func (si SearchInput) Focused() bool {
	return si.textInput.Focused()
}
