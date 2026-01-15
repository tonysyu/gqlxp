package config

import "github.com/charmbracelet/bubbles/key"

// GlobalKeymaps contains keymaps used across all TUI models
type GlobalKeymaps struct {
	Quit           key.Binding
	CommandPalette key.Binding
}

// toggleCmdPaletteKey used to both open and close CmdPalette
const toggleCmdPaletteKey string = "ctrl+p"
// toggleOverlayKey used to both open and close Detail Overlay
const toggleOverlayKey string = " "
// closeOverlayKeys used to close CmdPalette and Detail Overlays
var closeOverlayKeys []string = []string{"q", "esc"}

// newGlobalKeymaps creates a new GlobalKeymaps with default bindings
func newGlobalKeymaps() GlobalKeymaps {
	return GlobalKeymaps{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "ctrl+d"),
			key.WithHelp("⌃+c/⌃+d", "quit"),
		),
		CommandPalette: key.NewBinding(
			key.WithKeys(toggleCmdPaletteKey),
			key.WithHelp("⌃+p", "command palette"),
		),
	}
}

// MainKeymaps contains keymaps for the main xplr model
type MainKeymaps struct {
	GlobalKeymaps
	NextPanel, PrevPanel, NextGQLType, PrevGQLType, ToggleOverlay key.Binding
	SearchFocus, SearchSubmit, SearchClear                        key.Binding
}

// NewMainKeymaps creates a new MainKeymaps with default bindings
func NewMainKeymaps() MainKeymaps {
	return MainKeymaps{
		GlobalKeymaps: newGlobalKeymaps(),
		NextPanel: key.NewBinding(
			key.WithKeys("]", "tab"),
			key.WithHelp("]/tab", "next"),
		),
		PrevPanel: key.NewBinding(
			key.WithKeys("[", "shift+tab"),
			key.WithHelp("[/⇧+tab", "prev"),
		),
		NextGQLType: key.NewBinding(
			key.WithKeys("}"),
			key.WithHelp("}", "next type"),
		),
		PrevGQLType: key.NewBinding(
			key.WithKeys("{"),
			key.WithHelp("{", "prev type"),
		),
		ToggleOverlay: key.NewBinding(
			key.WithKeys(toggleOverlayKey),
			key.WithHelp("space", "details"),
		),
		SearchFocus: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		SearchSubmit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit search"),
		),
		SearchClear: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear search"),
		),
	}
}

// OverlayKeymaps contains keymaps for the overlay model
type OverlayKeymaps struct {
	GlobalKeymaps
	Close key.Binding
}

// NewOverlayKeymaps creates a new OverlayKeymaps with default bindings
func NewOverlayKeymaps() OverlayKeymaps {
	closeKeys := append(closeOverlayKeys, toggleOverlayKey)
	return OverlayKeymaps{
		GlobalKeymaps: newGlobalKeymaps(),
		Close: key.NewBinding(
			key.WithKeys(closeKeys...),
			key.WithHelp("space/q/esc", "close overlay"),
		),
	}
}

// LibSelectKeymaps contains keymaps for the library selection model
type LibSelectKeymaps struct {
	GlobalKeymaps
	Select key.Binding
}

// NewLibSelectKeymaps creates a new LibSelectKeymaps with default bindings
func NewLibSelectKeymaps() LibSelectKeymaps {
	return LibSelectKeymaps{
		GlobalKeymaps: newGlobalKeymaps(),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

// PanelKeymaps contains keymaps for panel tab navigation
type PanelKeymaps struct {
	NextTab, PrevTab key.Binding
}

// NewPanelKeymaps creates a new PanelKeymaps with default bindings
func NewPanelKeymaps() PanelKeymaps {
	return PanelKeymaps{
		NextTab: key.NewBinding(
			key.WithKeys("L", "shift+right"),
			key.WithHelp("L", "next tab"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("H", "shift+left"),
			key.WithHelp("H", "prev tab"),
		),
	}
}

// CommandPaletteKeymaps contains keymaps for the command palette
type CommandPaletteKeymaps struct {
	GlobalKeymaps
	Close   key.Binding
	Execute key.Binding
}

// NewCommandPaletteKeymaps creates a new CommandPaletteKeymaps with default bindings
func NewCommandPaletteKeymaps() CommandPaletteKeymaps {
	closeKeys := append(closeOverlayKeys, toggleCmdPaletteKey)
	return CommandPaletteKeymaps{
		GlobalKeymaps: newGlobalKeymaps(),
		Close: key.NewBinding(
			key.WithKeys(closeKeys...),
			key.WithHelp("space/q/esc", "close"),
		),
		Execute: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "execute"),
		),
	}
}
