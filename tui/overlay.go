package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// overlayModel manages overlay display and message interception
type overlayModel struct {
	active bool
	panel  Panel
	keymap overlayKeymap
}

type overlayKeymap struct {
	Close key.Binding
	Quit key.Binding
}

func newOverlayModel() overlayModel {
	return overlayModel{
		active: false,
		keymap: overlayKeymap{
			Close: key.NewBinding(
				key.WithKeys(" "),
				key.WithHelp("space", "close overlay"),
			),
			Quit: quitKeyBinding,
		},
	}
}

// Update processes messages and returns (model, cmd, intercepted)
// intercepted=true means the message was handled and should not be passed to main panels
func (o overlayModel) Update(msg tea.Msg) (overlayModel, tea.Cmd, bool) {
	if !o.active {
		return o, nil, false // pass through to main model
	}

	// Handle overlay-specific keys
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, o.keymap.Close):
			o.active = false
			return o, nil, true // intercepted
		case key.Matches(msg, o.keymap.Quit):
			return o, tea.Quit, true // intercepted
		}
	}

	// Update overlay panel with the message
	// TODO: Check if this panel update is actually necessary
	if o.panel != nil {
		newModel, cmd := o.panel.Update(msg)
		if panel, ok := newModel.(Panel); ok {
			o.panel = panel
		}
		return o, cmd, true // intercepted
	}

	return o, nil, true // intercepted (active but no panel)
}

// Show activates the overlay with the given panel and size
func (o *overlayModel) Show(panel Panel, width, height int) {
	o.panel = panel
	o.panel.SetSize(width-8, height-8) // Leave margin for border/padding
	o.active = true
}

// Hide deactivates the overlay
func (o *overlayModel) Hide() {
	o.active = false
}

// IsActive returns whether the overlay is currently shown
func (o overlayModel) IsActive() bool {
	return o.active
}

// View renders the overlay panel content
func (o overlayModel) View() string {
	if !o.active || o.panel == nil {
		return ""
	}
	return o.panel.View()
}

// GetCloseBinding returns the key binding for closing the overlay
func (o overlayModel) GetCloseBinding() key.Binding {
	return o.keymap.Close
}
