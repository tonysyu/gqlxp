package config

import "github.com/charmbracelet/lipgloss"

// Colors
// See https://hexdocs.pm/color_palette/ansi_color_codes.html for color codes
const (
	ColorLightGray    lipgloss.Color = "244" // gray
	ColorMidGray      lipgloss.Color = "240" // davys_grey
	ColorDarkGray     lipgloss.Color = "238" // dark_charcoal (dark gray)
	ColorDimWhite     lipgloss.Color = "253" // alto (off-white)
	ColorCream        lipgloss.Color = "230" // cream, very_pale_yellow
	ColorDimIndigo    lipgloss.Color = "62"  // indigo, slate_blue
	ColorBrightIndigo lipgloss.Color = "57"  // electric_indigo
	ColorDimMagenta   lipgloss.Color = "170" // orchid (pink/purple)
)

// Layout dimensions
const (
	VisiblePanelCount  = 2
	HelpHeight         = 5
	NavbarHeight       = 3
	BreadcrumbsHeight  = 1
	PanelTitleHPadding = 1
	ItemLeftPadding    = 2
	OverlayPadding     = 1
	OverlayMargin      = 2
)

// Styles contains all lipgloss styles used in the TUI
type Styles struct {
	// Panel styles for panels displaying lists of types, fields, etc.
	FocusedPanel lipgloss.Style
	BlurredPanel lipgloss.Style
	PanelTitle   lipgloss.Style

	// Navigation styles for navbar dislplaying GQL Type selection
	Navbar            lipgloss.Style
	ActiveTab         lipgloss.Style
	InactiveTab       lipgloss.Style
	Breadcrumbs       lipgloss.Style
	CurrentBreadcrumb lipgloss.Style

	// Overlay style for view displaying Details of GQL Types
	Overlay lipgloss.Style

	// Section and item styles for virtual navigation items
	SectionLabel  lipgloss.Style
	FocusedItem   lipgloss.Style
	UnfocusedItem lipgloss.Style
	Divider       lipgloss.Style
}

// DefaultStyles returns the default style configuration
func DefaultStyles() Styles {
	return Styles{
		FocusedPanel: lipgloss.NewStyle().
			Foreground(ColorLightGray).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorDarkGray),

		BlurredPanel: lipgloss.NewStyle().
			Foreground(ColorLightGray).
			Border(lipgloss.HiddenBorder()),

		// Panel title styling copied from bubbles/list
		PanelTitle: lipgloss.NewStyle().
			Background(ColorDimIndigo).
			Foreground(ColorCream).
			Padding(0, PanelTitleHPadding),

		Navbar: lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 0, 1, 0),

		ActiveTab: lipgloss.NewStyle().
			Foreground(ColorCream).
			Background(ColorBrightIndigo).
			Padding(0, 2).
			Bold(true),

		InactiveTab: lipgloss.NewStyle().
			Foreground(ColorLightGray).
			Padding(0, 2),

		Breadcrumbs: lipgloss.NewStyle().
			Foreground(ColorLightGray).
			Padding(0, 1),

		CurrentBreadcrumb: lipgloss.NewStyle().
			Foreground(ColorDimIndigo),

		Overlay: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorDarkGray).
			Padding(OverlayPadding).
			Margin(OverlayMargin),

		SectionLabel: lipgloss.NewStyle().
			Foreground(ColorMidGray).
			Bold(true).
			Padding(0, 1),

		FocusedItem: lipgloss.NewStyle().
			// Use left-border to indicated selected/focused item.
			// (Adapted from bubbles/list/defaultitem)
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(ColorDimMagenta).
			Foreground(ColorDimMagenta).
			Padding(0, 0, 0, ItemLeftPadding-1), // Subtract 1 due to left-border

		UnfocusedItem: lipgloss.NewStyle().
			Foreground(ColorDimWhite).
			Padding(0, 0, 0, ItemLeftPadding),

		Divider: lipgloss.NewStyle().
			Foreground(ColorDarkGray),
	}
}
