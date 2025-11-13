package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/tui/components"
)

func TestShouldPanelReceiveMessage(t *testing.T) {
	is := is.New(t)

	model := newModel(adapters.SchemaView{})

	tests := []struct {
		name          string
		displayOffset int
		msg           tea.Msg
		shouldReceive bool
	}{
		{
			name:          "left panel (offset 0) receives key message",
			displayOffset: 0,
			msg:           tea.KeyMsg{Type: tea.KeyEnter},
			shouldReceive: true,
		},
		{
			name:          "all panels receive window size message",
			displayOffset: 0,
			msg:           tea.WindowSizeMsg{Width: 100, Height: 50},
			shouldReceive: true,
		},
		{
			name:          "global navigation keys not sent to panels",
			displayOffset: 0,
			msg:           tea.KeyMsg{Type: tea.KeyTab},
			shouldReceive: false,
		},
		{
			name:          "OpenPanelMsg not sent to panels",
			displayOffset: 0,
			msg:           components.OpenPanelMsg{Panel: components.NewEmptyListPanel("test")},
			shouldReceive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.shouldFocusedPanelReceiveMessage(tt.msg)
			is.Equal(result, tt.shouldReceive)
		})
	}
}

func TestGlobalNavigationKeysNotSentToPanels(t *testing.T) {
	is := is.New(t)

	model := newModel(adapters.SchemaView{})

	// Test all global navigation keys
	globalKeys := []tea.KeyMsg{
		{Type: tea.KeyTab},      // next
		{Type: tea.KeyShiftTab}, // prev
		{Type: tea.KeyCtrlC},    // quit
		{Type: tea.KeyCtrlD},    // quit
		{Type: tea.KeyCtrlT},    // toggle
	}

	for _, keyMsg := range globalKeys {
		// Even the left panel (offset 0) should not receive global navigation keys
		shouldReceive := model.shouldFocusedPanelReceiveMessage(keyMsg)
		is.True(!shouldReceive)
	}
}
