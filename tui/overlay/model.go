package overlay

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/utils/terminal"
	"github.com/tonysyu/gqlxp/utils/text"
)

var (
	quitKeyBinding = key.NewBinding(
		key.WithKeys("ctrl+c", "ctrl+d"),
		key.WithHelp("ctrl+c", "quit"),
	)
)

// Panel inside the overlay must be inset by padding, margin, and a 1-char border on all sides.
var overlayPanelMargin = 2 * (config.OverlayMargin + config.OverlayPadding + 1)

// Model manages overlay display and message interception
type Model struct {
	active   bool
	viewport viewport.Model
	renderer terminal.Renderer
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

// New creates a new overlay model
func New(styles config.Styles) Model {
	vp := viewport.New(0, 0)

	// Initialize markdown renderer once for the lifetime of the session
	renderer, err := terminal.NewMarkdownRenderer()
	// If renderer fails, renderer will be nil and we'll use plain content
	model := Model{
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
func (o Model) Update(msg tea.Msg) (Model, tea.Cmd, bool) {
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
func (o *Model) Show(content string, width, height int) {
	o.content = content
	o.active = true

	// Set viewport size
	viewportWidth := width - overlayPanelMargin
	viewportHeight := height - overlayPanelMargin - config.HelpHeight
	o.viewport.Width = viewportWidth
	o.viewport.Height = viewportHeight

	// Render markdown content using the shared glamour renderer
	if viewportWidth > 0 {
		rendered := terminal.RenderMarkdownOrPlain(o.renderer, content)
		o.rendered = rendered
		o.viewport.SetContent(rendered)
		return
	}

	// Fallback to plain content if viewport not properly sized
	o.viewport.SetContent(content)
}

// Hide deactivates the overlay
func (o *Model) Hide() {
	o.active = false
}

// IsActive returns whether the overlay is currently shown
func (o Model) IsActive() bool {
	return o.active
}

// View renders the overlay viewport content with help
func (o Model) View() string {
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
