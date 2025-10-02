package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	tea "github.com/charmbracelet/bubbletea"
)

// Panel inside the overlay must be inset by padding, margin, and a 1-char border on all sides.
var overlayPanelMargin = 2 * (overlayMargin + overlayPadding + 1)

// overlayModel manages overlay display and message interception
type overlayModel struct {
	active    bool
	viewport  viewport.Model
	renderer  *glamour.TermRenderer
	content   string // original markdown content
	lastWidth int    // track last rendered width to avoid unnecessary re-renders
	rendered  string // cache rendered content
	keymap    overlayKeymap
	help      help.Model
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
	vp := viewport.New(0, 0)

	// Initialize glamour renderer once for the lifetime of the session
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)
	// If glamour fails, renderer will be nil and we'll use plain content

	model := overlayModel{
		active:    false,
		viewport:  vp,
		renderer:  nil,
		lastWidth: 0,
		help:      help.New(),
		keymap: overlayKeymap{
			Close: key.NewBinding(
				key.WithKeys(" "),
				key.WithHelp("space", "close overlay"),
			),
			Quit: quitKeyBinding,
		},
	}

	if err == nil {
		model.renderer = renderer
	}

	return model
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

	// Update viewport with the message
	var cmd tea.Cmd
	o.viewport, cmd = o.viewport.Update(msg)
	return o, cmd, true // intercepted
}

// Show activates the overlay with the given markdown content and size
func (o *overlayModel) Show(content string, width, height int) {
	o.content = content
	o.active = true

	// Set viewport size
	viewportWidth := width - overlayPanelMargin
	viewportHeight := height - overlayPanelMargin - helpHeight
	o.viewport.Width = viewportWidth
	o.viewport.Height = viewportHeight

	// Render markdown content using the shared glamour renderer
	if o.renderer != nil && viewportWidth > 0 {
		// Only recreate renderer if width has changed
		if viewportWidth != o.lastWidth {
			renderer, err := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(viewportWidth),
			)
			if err == nil {
				o.renderer = renderer
				o.lastWidth = viewportWidth
			}
		}

		// Render content with current renderer
		rendered, err := o.renderer.Render(content)
		if err == nil {
			o.rendered = rendered
			o.viewport.SetContent(rendered)
			return
		}
	}

	// Fallback to plain content if glamour fails or is unavailable
	o.viewport.SetContent(content)
}

// Hide deactivates the overlay
func (o *overlayModel) Hide() {
	o.active = false
}

// IsActive returns whether the overlay is currently shown
func (o overlayModel) IsActive() bool {
	return o.active
}

// View renders the overlay viewport content with help
func (o overlayModel) View() string {
	if !o.active {
		return ""
	}
	helpView := o.help.ShortHelpView(o.keymap.ShortHelp())
	return o.viewport.View() + "\n\n" + helpView
}
