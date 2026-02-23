package overlay

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/tui/utils"
	"github.com/tonysyu/gqlxp/utils/terminal"
	"github.com/tonysyu/gqlxp/utils/text"
)

// Panel inside the overlay must be inset by padding, margin, and a 1-char border on all sides.
var overlayPanelMargin = 2 * (config.OverlayMargin + config.OverlayPadding + 1)

// ClosedMsg is sent when the overlay requests to be closed
type ClosedMsg struct{}

// Model manages overlay display and message interception
type Model struct {
	viewport viewport.Model
	renderer terminal.Renderer
	content  string // original markdown content
	rendered string // cache rendered content
	Styles   config.Styles

	width  int
	height int
	keymap config.OverlayKeymaps
	help   help.Model
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (m Model) ShortHelp() []key.Binding {
	return []key.Binding{
		m.keymap.Close,
		m.keymap.Quit,
	}
}

// New creates a new overlay model
func New(styles config.Styles) Model {
	vp := viewport.New(0, 0)

	// Initialize markdown renderer once for the lifetime of the session
	renderer, err := terminal.NewMarkdownRenderer()
	// If renderer fails, renderer will be nil and we'll use plain content
	model := Model{
		viewport: vp,
		renderer: nil,
		help:     help.New(),
		Styles:   styles,
		keymap:   config.NewOverlayKeymaps(),
	}

	if err == nil {
		model.renderer = renderer
	}

	return model
}

// Update processes messages and returns (model, cmd).
// xplr.Model is responsible for routing messages here only when overlay is active.
func (o Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// Handle overlay-specific keys
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, o.keymap.Close):
			return o, func() tea.Msg { return ClosedMsg{} }
		case key.Matches(msg, o.keymap.Quit):
			return o, tea.Quit
		}
	case tea.WindowSizeMsg:
		o.height = msg.Height
		o.width = msg.Width
	}

	// Update viewport with the message
	var cmd tea.Cmd
	o.viewport, cmd = o.viewport.Update(msg)
	return o, cmd
}

// Show configures the overlay with the given markdown content and size.
// xplr.Model is responsible for setting its state to xplrOverlayView when calling this.
func (o Model) Show(content string, width, height int) Model {
	o.content = content

	// Set viewport size
	viewportWidth := width - config.OverlayInsetMargin
	viewportHeight := height - config.OverlayInsetMargin - config.HelpHeight
	o.viewport.Width = viewportWidth
	o.viewport.Height = viewportHeight

	// Render markdown content using the shared glamour renderer
	if viewportWidth > 0 {
		rendered := terminal.RenderMarkdownOrPlain(o.renderer, content)
		o.rendered = rendered
		o.viewport.SetContent(rendered)
		return o
	}

	// Fallback to plain content if viewport not properly sized
	o.viewport.SetContent(content)
	return o
}

// View renders the overlay viewport content with help
func (o Model) View() string {
	helpView := o.help.ShortHelpView(o.ShortHelp())
	content := text.JoinParagraphs(o.viewport.View(), helpView)
	overlay := o.Styles.Overlay.Render(content)
	return utils.CenterOverlay(overlay, o.width, o.height)
}
