package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/config"
)

func showDefaultOverlay() overlayModel {
	overlay := newOverlayModel(config.DefaultStyles())
	overlay.Show("test content", 100, 50)
	return overlay
}

func TestInactiveOverlay(t *testing.T) {
	is := is.New(t)
	overlay := newOverlayModel(config.DefaultStyles())

	t.Run("NewOverlayModel initializes correctly", func(t *testing.T) {
		is.Equal(overlay.active, false)
		is.Equal(overlay.content, "")
		is.True(overlay.keymap.Close.Enabled())
		is.True(overlay.keymap.Quit.Enabled())
	})

	t.Run("Inactive overlay passes through messages", func(t *testing.T) {
		// When inactive, should not intercept messages
		updatedOverlay, cmd, intercepted := overlay.Update(tea.KeyMsg{Type: tea.KeyEnter})

		is.Equal(intercepted, false)
		is.True(cmd == nil)
		is.Equal(updatedOverlay.active, false)
	})

	t.Run("Show and Hide toggle overlay state", func(t *testing.T) {
		overlay := showDefaultOverlay()

		is.Equal(overlay.active, true)
		is.Equal(overlay.content, "test content")
		is.Equal(overlay.IsActive(), true)

		// Test Hide
		overlay.Hide()
		is.Equal(overlay.active, false)
		is.Equal(overlay.IsActive(), false)
	})

	t.Run("Close key binding deactivates overlay", func(t *testing.T) {
		overlay := showDefaultOverlay()

		is.Equal(overlay.active, true)

		spaceKey := tea.KeyMsg{Type: tea.KeySpace}
		updatedOverlay, cmd, intercepted := overlay.Update(spaceKey)

		is.Equal(intercepted, true)
		is.True(cmd == nil)
		is.Equal(updatedOverlay.active, false)
	})

	t.Run("Quit key binding returns quit command", func(t *testing.T) {
		overlay := showDefaultOverlay()

		// Send quit key (ctrl+c)
		quitKey := tea.KeyMsg{Type: tea.KeyCtrlC}
		updatedOverlay, cmd, intercepted := overlay.Update(quitKey)

		is.Equal(intercepted, true)
		is.True(cmd != nil)                   // Should return quit command
		is.Equal(updatedOverlay.active, true) // Active state unchanged
	})

	t.Run("Active overlay intercepts messages", func(t *testing.T) {
		overlay := showDefaultOverlay()

		// Any key message should be intercepted when active
		testKey := tea.KeyMsg{Type: tea.KeyEnter}
		_, _, intercepted := overlay.Update(testKey)

		is.Equal(intercepted, true)
	})

	t.Run("Viewport receives update messages", func(t *testing.T) {
		overlay := showDefaultOverlay()

		// Send a message to update viewport
		msg := tea.KeyMsg{Type: tea.KeyDown}
		updatedOverlay, _, intercepted := overlay.Update(msg)

		is.Equal(intercepted, true)
		is.Equal(updatedOverlay.active, true)
	})

	t.Run("Show sets viewport size with margin", func(t *testing.T) {
		width, height := 200, 100
		overlay := newOverlayModel(config.DefaultStyles())
		overlay.Show("test content", width, height)

		is.Equal(overlay.viewport.Width, width-overlayPanelMargin)
		is.Equal(overlay.viewport.Height, height-overlayPanelMargin-config.HelpHeight)
	})
}
