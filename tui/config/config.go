package config

import "github.com/charmbracelet/lipgloss"

// Layout dimensions
const (
	VisiblePanelCount = 2
	HelpHeight        = 5
	NavbarHeight      = 3
	OverlayPadding    = 1
	OverlayMargin     = 2
)

// Styles contains all lipgloss styles used in the TUI
type Styles struct {
	// Cursor styles
	Cursor     lipgloss.Style
	CursorLine lipgloss.Style

	// Placeholder styles
	Placeholder        lipgloss.Style
	FocusedPlaceholder lipgloss.Style
	EndOfBuffer        lipgloss.Style

	// Border styles
	FocusedBorder lipgloss.Style
	BlurredBorder lipgloss.Style

	// Navigation styles
	Navbar      lipgloss.Style
	ActiveTab   lipgloss.Style
	InactiveTab lipgloss.Style

	// Overlay style
	Overlay lipgloss.Style
}

// DefaultStyles returns the default style configuration
func DefaultStyles() Styles {
	return Styles{
		// Cursor styles
		Cursor: lipgloss.NewStyle().Foreground(lipgloss.Color("212")),

		CursorLine: lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230")),

		// Placeholder styles
		Placeholder: lipgloss.NewStyle().
			Foreground(lipgloss.Color("238")),

		FocusedPlaceholder: lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")),

		EndOfBuffer: lipgloss.NewStyle().
			Foreground(lipgloss.Color("235")),

		// Border styles
		FocusedBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")),

		BlurredBorder: lipgloss.NewStyle().
			Border(lipgloss.HiddenBorder()),

		// Navigation styles
		Navbar: lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 0, 1, 0),

		ActiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("57")).
			Padding(0, 2).
			Bold(true),

		InactiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Padding(0, 2),

		// Overlay style
		Overlay: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(OverlayPadding).
			Margin(OverlayMargin),
	}
}
