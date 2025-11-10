package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/utils/text"
)

// Panel inside the overlay must be inset by padding, margin, and a 1-char border on all sides.
var overlayPanelMargin = 2 * (config.OverlayMargin + config.OverlayPadding + 1)

// overlayModel manages overlay display and message interception
type overlayModel struct {
	active   bool
	viewport viewport.Model
	renderer *glamour.TermRenderer
	content  string // original markdown content
	rendered string // cache rendered content
	Styles   config.Styles

	width  int
	height int
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

func newOverlayModel(styles config.Styles) overlayModel {
	vp := viewport.New(0, 0)

	// Initialize glamour renderer once for the lifetime of the session
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)
	// If glamour fails, renderer will be nil and we'll use plain content
	model := overlayModel{
		active:   false,
		viewport: vp,
		renderer: nil,
		help:     help.New(),
		Styles:   styles,
		keymap: overlayKeymap{
			Close: key.NewBinding(
				key.WithKeys(" ", "q"),
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
	case tea.WindowSizeMsg:
		o.height = msg.Height
		o.width = msg.Width
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
	viewportHeight := height - overlayPanelMargin - config.HelpHeight
	o.viewport.Width = viewportWidth
	o.viewport.Height = viewportHeight

	// Render markdown content using the shared glamour renderer
	if o.renderer != nil && viewportWidth > 0 {
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
	content := text.JoinParagraphs(o.viewport.View(), helpView)

	overlay := o.Styles.Overlay.Render(content)

	// Center the overlay on screen
	overlayHeight := lipgloss.Height(overlay)
	overlayWidth := lipgloss.Width(overlay)

	verticalMargin := (o.height - overlayHeight) / 2
	horizontalMargin := (o.width - overlayWidth) / 2

	if verticalMargin < 0 {
		verticalMargin = 0
	}
	if horizontalMargin < 0 {
		horizontalMargin = 0
	}

	// Position the overlay over the main view
	positionedOverlay := lipgloss.NewStyle().
		MarginTop(verticalMargin).
		MarginLeft(horizontalMargin).
		Render(overlay)

	return lipgloss.Place(o.width, o.height, lipgloss.Center, lipgloss.Center, positionedOverlay)
}
