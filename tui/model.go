package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tonysyu/gqlxp/library"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/libselect"
	"github.com/tonysyu/gqlxp/tui/xplr"
)

// sessionState tracks which submodel is active
type sessionState uint

const (
	libselectView sessionState = iota
	xplrView
)

// Model is the top-level TUI model that delegates to submodes
type Model struct {
	state     sessionState
	libselect libselect.Model
	xplr      xplr.Model
	width     int
	height    int
}

// newModelWithLibselect creates a model starting in library selection mode
func newModelWithLibselect() (Model, error) {
	lib := library.NewLibrary()
	libselectModel, err := libselect.New(lib)
	if err != nil {
		return Model{}, err
	}

	return Model{
		state:     libselectView,
		libselect: libselectModel,
		xplr:      xplr.NewEmpty(),
	}, nil
}

// newModelWithXplr creates a model starting in explorer mode
func newModelWithXplr(schema adapters.SchemaView) Model {
	return Model{
		state: xplrView,
		xplr:  xplr.New(schema),
	}
}

// newModelWithXplrAndLibrary creates a model starting in explorer mode with library data
func newModelWithXplrAndLibrary(schema adapters.SchemaView, schemaID string, metadata library.SchemaMetadata) Model {
	return Model{
		state: xplrView,
		xplr:  xplr.NewFromSchemaLibrary(schema, schemaID, metadata),
	}
}

func (m Model) Init() tea.Cmd {
	switch m.state {
	case libselectView:
		return m.libselect.Init()
	case xplrView:
		return m.xplr.Init()
	default:
		return nil
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Track window size at top level and forward to active submodel
	// Also forward to xplr if we're in libselect mode (so xplr has correct size on transition)
	if windowMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = windowMsg.Width
		m.height = windowMsg.Height

		// Always forward to the xplr model (even when not active)
		// This ensures it has correct dimensions when we transition to it
		m.xplr, _ = m.xplr.Update(windowMsg)

		// Also forward to libselect if it's the active view
		if m.state == libselectView {
			m.libselect, _ = m.libselect.Update(windowMsg)
		}
	}

	// Handle transitions between submodes
	switch msg := msg.(type) {
	case libselect.SchemaSelectedMsg:
		// Transition from libselect to xplr by sending schema to existing xplr model
		m.state = xplrView
		schemaLoadedMsg := xplr.SchemaLoadedMsg{
			Schema:         msg.Schema,
			SchemaID:       msg.SchemaID,
			Favorites:      msg.Metadata.Favorites,
			HasLibraryData: true,
		}
		m.xplr, cmd = m.xplr.Update(schemaLoadedMsg)
		return m, tea.Batch(cmd, m.xplr.Init())
	}

	// Delegate to active submodel
	switch m.state {
	case libselectView:
		m.libselect, cmd = m.libselect.Update(msg)
		return m, cmd
	case xplrView:
		m.xplr, cmd = m.xplr.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m Model) View() string {
	switch m.state {
	case libselectView:
		return m.libselect.View()
	case xplrView:
		return m.xplr.View()
	default:
		return ""
	}
}
