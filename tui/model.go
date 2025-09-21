package tui

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/tonysyu/gq/gql"
)

const (
	intialPanels = 2
	maxPanes     = 6
	minPanes     = 1
	helpHeight   = 5
	navbarHeight = 3
)

type FieldType string

const (
	QueryType    FieldType = "Query"
	MutationType FieldType = "Mutation"
)

// availableFieldTypes defines the ordered list of field types for navigation
var availableFieldTypes = []FieldType{QueryType, MutationType}

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())

	// Navbar styles
	navbarStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 0, 1, 0)

	activeTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("57")).
			Padding(0, 2).
			Bold(true)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("244")).
				Padding(0, 2)
)

type keymap = struct {
	next, prev, quit, toggle key.Binding
}

type mainModel struct {
	width     int
	height    int
	keymap    keymap
	help      help.Model
	panels    []Panel
	focus     int
	schema    gql.GraphQLSchema
	fieldType FieldType
}

func NewModel(schema gql.GraphQLSchema) mainModel {
	m := mainModel{
		panels:    make([]Panel, intialPanels),
		help:      help.New(),
		schema:    schema,
		fieldType: QueryType,
		keymap: keymap{
			next: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next"),
			),
			prev: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "prev"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
			toggle: key.NewBinding(
				key.WithKeys("ctrl+t"),
				key.WithHelp("ctrl+t", "toggle Query/Mutation"),
			),
		},
	}
	// Initialize panels with empty list models
	for i := range intialPanels {
		m.panels[i] = newListPanel([]list.Item{}, "")
	}

	// Load initial fields based on field type
	m.loadFieldsPanel()
	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.next):
			m.focus++
			if m.focus > len(m.panels)-1 {
				m.focus = 0
			}
		case key.Matches(msg, m.keymap.prev):
			m.focus--
			if m.focus < 0 {
				m.focus = len(m.panels) - 1
			}
		case key.Matches(msg, m.keymap.toggle):
			m.toggleFieldType()
		}
	case openPanelMsg:
		m.handleOpenPanel(msg.panel)
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	m.sizePanels()

	// Update all panels
	for i := range m.panels {
		newModel, cmd := m.panels[i].Update(msg)
		if panel, ok := newModel.(Panel); ok {
			m.panels[i] = panel
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *mainModel) sizePanels() {
	for i := range m.panels {
		m.panels[i].SetSize(m.width/len(m.panels), m.height-helpHeight-navbarHeight)
	}
}

// addPanel adds a new panel to the model
func (m *mainModel) addPanel(panel Panel) {
	if len(m.panels) < maxPanes {
		m.panels = append(m.panels, panel)
		m.sizePanels()
	}
}

// addStringPanel is a convenience method to add a string panel
func (m *mainModel) addStringPanel(content string) {
	m.addPanel(newStringPanel(content))
}

// addListPanel is a convenience method to add a list panel with list.Item interface
func (m *mainModel) addListPanel(items []list.Item) {
	m.addPanel(newListPanel(items, ""))
}

// handleOpenPanel handles when an item is opened
func (m *mainModel) handleOpenPanel(newPanel Panel) {
	nextPanelIndex := m.focus + 1

	// If there's a next panel, replace it
	if nextPanelIndex < len(m.panels) {
		m.panels[nextPanelIndex] = newPanel
	} else if len(m.panels) < maxPanes {
		// Add a new panel if we haven't reached the max
		m.addPanel(newPanel)
	}

	m.sizePanels()
}

// loadFieldsPanel loads the appropriate fields based on the current field type
func (m *mainModel) loadFieldsPanel() {
	var fields map[string]*ast.FieldDefinition

	switch m.fieldType {
	case QueryType:
		fields = m.schema.Query
	case MutationType:
		fields = m.schema.Mutation
	}

	title := fmt.Sprintf("%s Fields", string(m.fieldType))
	items := adaptGraphQLItems(fields)
	m.panels[0] = newListPanel(items, title)
}

// toggleFieldType cycles through available field types with wraparound
func (m *mainModel) toggleFieldType() {
	// Find current field type index
	currentIndex := slices.IndexFunc(availableFieldTypes, func(fieldType FieldType) bool {
		return m.fieldType == fieldType
	})

	// Cycle to next field type with wraparound
	nextIndex := (currentIndex + 1) % len(availableFieldTypes)
	m.fieldType = availableFieldTypes[nextIndex]

	m.loadFieldsPanel()
	m.sizePanels()
}

// renderFieldTypeNavbar creates the navbar showing field types
func (m *mainModel) renderFieldTypeNavbar() string {
	var tabs []string

	for _, fieldType := range availableFieldTypes {
		var style lipgloss.Style
		if m.fieldType == fieldType {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		tabs = append(tabs, style.Render(string(fieldType)))
	}

	navbar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	return navbarStyle.Render(navbar)
}

func (m mainModel) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.next,
		m.keymap.prev,
		m.keymap.toggle,
		m.keymap.quit,
	})

	var views []string
	for i := range m.panels {
		views = append(views, m.panels[i].View())
	}

	navbar := m.renderFieldTypeNavbar()
	panels := lipgloss.JoinHorizontal(lipgloss.Top, views...)

	return navbar + "\n" + panels + "\n\n" + help
}
