package overlay

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/config"
)

func showDefaultOverlay() Model {
	overlay := New(config.DefaultStyles())
	overlay = overlay.Show("test content", 100, 50)
	return overlay
}

func TestOverlay(t *testing.T) {
	is := is.New(t)
	overlay := New(config.DefaultStyles())

	t.Run("NewOverlayModel initializes correctly", func(t *testing.T) {
		is.Equal(overlay.content, "")
		is.True(overlay.keymap.Close.Enabled())
		is.True(overlay.keymap.Quit.Enabled())
	})

	t.Run("Show sets content and viewport size", func(t *testing.T) {
		overlay := showDefaultOverlay()

		is.Equal(overlay.content, "test content")
	})

	t.Run("Close key returns ClosedMsg command", func(t *testing.T) {
		overlay := showDefaultOverlay()

		spaceKey := tea.KeyMsg{Type: tea.KeySpace}
		_, cmd := overlay.Update(spaceKey)

		is.True(cmd != nil)
		msg := cmd()
		_, ok := msg.(ClosedMsg)
		is.True(ok)
	})

	t.Run("Quit key returns quit command", func(t *testing.T) {
		overlay := showDefaultOverlay()

		quitKey := tea.KeyMsg{Type: tea.KeyCtrlC}
		_, cmd := overlay.Update(quitKey)

		is.True(cmd != nil)
	})

	t.Run("Viewport receives update messages", func(t *testing.T) {
		overlay := showDefaultOverlay()

		msg := tea.KeyMsg{Type: tea.KeyDown}
		updatedOverlay, _ := overlay.Update(msg)

		is.Equal(updatedOverlay.content, "test content")
	})

	t.Run("Show sets viewport size with margin", func(t *testing.T) {
		width, height := 200, 100
		overlay := New(config.DefaultStyles())
		overlay = overlay.Show("test content", width, height)

		is.Equal(overlay.viewport.Width, width-overlayPanelMargin)
		is.Equal(overlay.viewport.Height, height-overlayPanelMargin-config.HelpHeight)
	})
}
