package terminal

import (
	"github.com/charmbracelet/lipgloss"
)

// ANSI color codes for terminal styling
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
