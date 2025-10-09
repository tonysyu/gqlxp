package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/tonysyu/igq/tui/adapters"
	"github.com/tonysyu/igq/tui/components"
)

func TestShouldPanelReceiveMessage(t *testing.T) {
	is := is.New(t)

	model := newModel(adapters.SchemaView{})
	model.focus = 1 // Set focus to second panel

	tests := []struct {
		name          string
		panelIndex    int
		msg           tea.Msg
		shouldReceive bool
	}{
		{
			name:          "focused panel receives key message",
			panelIndex:    1,
			msg:           tea.KeyMsg{Type: tea.KeyEnter},
			shouldReceive: true,
		},
		{
			name:          "unfocused panel does not receive key message",
			panelIndex:    0,
			msg:           tea.KeyMsg{Type: tea.KeyEnter},
			shouldReceive: false,
		},
		{
			name:          "all panels receive window size message",
			panelIndex:    0,
			msg:           tea.WindowSizeMsg{Width: 100, Height: 50},
			shouldReceive: true,
		},
		{
			name:          "global navigation keys not sent to panels",
			panelIndex:    1,
			msg:           tea.KeyMsg{Type: tea.KeyTab},
			shouldReceive: false,
		},
		{
			name:          "OpenPanelMsg not sent to panels",
			panelIndex:    1,
			msg:           components.OpenPanelMsg{Panel: components.NewStringPanel("test")},
			shouldReceive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.shouldPanelReceiveMessage(tt.panelIndex, tt.msg)
			is.Equal(result, tt.shouldReceive)
		})
	}
}

func TestGlobalNavigationKeysNotSentToPanels(t *testing.T) {
	is := is.New(t)

	model := newModel(adapters.SchemaView{})
	model.focus = 0

	// Test all global navigation keys
	globalKeys := []tea.KeyMsg{
		{Type: tea.KeyTab},      // next
		{Type: tea.KeyShiftTab}, // prev
		{Type: tea.KeyCtrlC},    // quit
		{Type: tea.KeyCtrlD},    // quit
		{Type: tea.KeyCtrlT},    // toggle
	}

	for _, keyMsg := range globalKeys {
		// Even the focused panel should not receive global navigation keys
		shouldReceive := model.shouldPanelReceiveMessage(model.focus, keyMsg)
		is.True(!shouldReceive)
	}
}
