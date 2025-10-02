package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Panel inside the overlay must be inset by padding, margin, and a 1-char border on all sides.
var overlayPanelMargin = 2 * (overlayMargin + overlayPadding + 1)

// overlayModel manages overlay display and message interception
type overlayModel struct {
	active bool
	panel  Panel
	keymap overlayKeymap
	help   help.Model
}

type overlayKeymap struct {
	Close key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k overlayKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Close, k.Quit}
}

func newOverlayModel() overlayModel {
	return overlayModel{
		active: false,
		help:   help.New(),
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
	o.panel.SetSize(
		width-overlayPanelMargin,
		height-overlayPanelMargin-helpHeight,
	)
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

// View renders the overlay panel content with help
func (o overlayModel) View() string {
	if !o.active || o.panel == nil {
		return ""
	}
	helpView := o.help.ShortHelpView(o.keymap.ShortHelp())
	return o.panel.View() + "\n\n" + helpView
}
