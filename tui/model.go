package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	intialPanels = 2
	maxPanes     = 6
	minPanes     = 1
	helpHeight   = 5
)

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
)

type keymap = struct {
	next, prev, add, remove, quit key.Binding
}

type mainModel struct {
	width  int
	height int
	keymap keymap
	help   help.Model
	panels []list.Model
	focus  int
}

func newListModel[T list.Item](choices []T) list.Model {
	items := make([]list.Item, len(choices))
	for i, choice := range choices {
		items[i] = choice
	}
	return list.New(items, list.NewDefaultDelegate(), 0, 0)
}

func NewModel[T list.Item](items []T) mainModel {
	m := mainModel{
		panels: make([]list.Model, intialPanels),
		help:   help.New(),
		keymap: keymap{
			next: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next"),
			),
			prev: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "prev"),
			),
			add: key.NewBinding(
				key.WithKeys("ctrl+n"),
				key.WithHelp("ctrl+n", "add an editor"),
			),
			remove: key.NewBinding(
				key.WithKeys("ctrl+w"),
				key.WithHelp("ctrl+w", "remove an editor"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
		},
	}
	// Initialize panels with empty list models
	for i := range intialPanels {
		m.panels[i] = newListModel([]list.Item{})
	}

	m.panels[0] = newListModel(items)
	m.updateKeybindings()
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
		case key.Matches(msg, m.keymap.add):
			fmt.Println("TODO")
		case key.Matches(msg, m.keymap.remove):
			m.panels = m.panels[:len(m.panels)-1]
			if m.focus > len(m.panels)-1 {
				m.focus = len(m.panels) - 1
			}
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	m.updateKeybindings()
	m.sizePanels()

	// Update all panels
	for i := range m.panels {
		newModel, cmd := m.panels[i].Update(msg)
		m.panels[i] = newModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *mainModel) sizePanels() {
	for i := range m.panels {
		m.panels[i].SetWidth(m.width / len(m.panels))
		m.panels[i].SetHeight(m.height - helpHeight)
	}
}

func (m *mainModel) updateKeybindings() {
	m.keymap.add.SetEnabled(len(m.panels) < maxPanes)
	m.keymap.remove.SetEnabled(len(m.panels) > minPanes)
}

func (m mainModel) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.next,
		m.keymap.prev,
		m.keymap.add,
		m.keymap.remove,
		m.keymap.quit,
	})

	var views []string
	for i := range m.panels {
		views = append(views, m.panels[i].View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n" + help
}
