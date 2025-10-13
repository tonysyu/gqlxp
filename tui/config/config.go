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
	// Panel styles for panels displaying lists of types, fields, etc.
	FocusedPanel lipgloss.Style
	BlurredPanel lipgloss.Style
	PanelTitle lipgloss.Style

	// Navigation styles for navbar dislplaying GQL Type selection
	Navbar      lipgloss.Style
	ActiveTab   lipgloss.Style
	InactiveTab lipgloss.Style

	// Overlay style for view displaying Details of GQL Types
	Overlay lipgloss.Style
}

// DefaultStyles returns the default style configuration
// See https://hexdocs.pm/color_palette/ansi_color_codes.html for color codes
func DefaultStyles() Styles {
	return Styles{
		FocusedPanel: lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")). // 244 = gray
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")), // 238 = outer_space (dark gray)

		BlurredPanel: lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")). // 244 = gray
			Border(lipgloss.HiddenBorder()),

		// Panel title styling copied from bubbles/list
		PanelTitle: lipgloss.NewStyle().
			Background(lipgloss.Color("62")).  // 62 = indigo, slate_blue
			Foreground(lipgloss.Color("230")). // 230 = cream, very_pale_yellow
			Padding(0, 1),

		Navbar: lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 0, 1, 0),

		ActiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")). // 230 = cream, very_pale_yellow
			Background(lipgloss.Color("57")).  // 57 = electric_indigo
			Padding(0, 2).
			Bold(true),

		InactiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")). // 244 = gray
			Padding(0, 2),

		Overlay: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")). // 238 = outer_space (dark gray)
			Padding(OverlayPadding).
			Margin(OverlayMargin),
	}
}
