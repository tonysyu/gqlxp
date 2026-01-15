package utils

import "github.com/charmbracelet/lipgloss"

// CenterOverlay positions rendered content in the center of the screen.
func CenterOverlay(content string, width, height int) string {
	overlayHeight := lipgloss.Height(content)
	overlayWidth := lipgloss.Width(content)

	verticalMargin := max((height-overlayHeight)/2, 0)
	horizontalMargin := max((width-overlayWidth)/2, 0)

	positioned := lipgloss.NewStyle().
		MarginTop(verticalMargin).
		MarginLeft(horizontalMargin).
		Render(content)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, positioned)
}
