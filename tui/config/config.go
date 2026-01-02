package config

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tonysyu/gqlxp/utils/terminal"
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
	ActiveSubTab      lipgloss.Style
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
			Foreground(terminal.ColorLightGray).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(terminal.ColorDarkGray),

		BlurredPanel: lipgloss.NewStyle().
			Foreground(terminal.ColorLightGray).
			Border(lipgloss.HiddenBorder()),

		// Panel title styling copied from bubbles/list
		PanelTitle: lipgloss.NewStyle().
			Background(terminal.ColorDimIndigo).
			Foreground(terminal.ColorCream).
			Padding(0, PanelTitleHPadding),

		Navbar: lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 0, 1, 0),

		ActiveTab: lipgloss.NewStyle().
			Foreground(terminal.ColorCream).
			Background(terminal.ColorBrightIndigo).
			Padding(0, 2).
			Bold(true),

		ActiveSubTab: lipgloss.NewStyle().
			Foreground(terminal.ColorDarkGray).
			Background(terminal.ColorLightBlue).
			Padding(0, 2).
			Bold(true),

		InactiveTab: lipgloss.NewStyle().
			Foreground(terminal.ColorLightGray).
			Padding(0, 2),

		Breadcrumbs: lipgloss.NewStyle().
			Foreground(terminal.ColorLightGray).
			Padding(0, 1),

		CurrentBreadcrumb: lipgloss.NewStyle().
			Foreground(terminal.ColorDimIndigo),

		Overlay: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(terminal.ColorDarkGray).
			Padding(OverlayPadding).
			Margin(OverlayMargin),

		SectionLabel: lipgloss.NewStyle().
			Foreground(terminal.ColorMidGray).
			Bold(true).
			Padding(0, 1),

		FocusedItem: lipgloss.NewStyle().
			// Use left-border to indicated selected/focused item.
			// (Adapted from bubbles/list/defaultitem)
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(terminal.ColorDimMagenta).
			Foreground(terminal.ColorDimMagenta).
			Padding(0, 0, 0, ItemLeftPadding-1), // Subtract 1 due to left-border

		UnfocusedItem: lipgloss.NewStyle().
			Foreground(terminal.ColorDimWhite).
			Padding(0, 0, 0, ItemLeftPadding),

		Divider: lipgloss.NewStyle().
			Foreground(terminal.ColorDarkGray),
	}
}
