package searchmodel_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/config"
	"github.com/tonysyu/gqlxp/tui/xplr/searchmodel"
)

func newTestModel() searchmodel.Model {
	return searchmodel.New(config.NewMainKeymaps())
}

func TestNew_NotFocused(t *testing.T) {
	is := is.New(t)
	m := newTestModel()
	is.True(!m.IsFocused())
	is.True(m.Results() == nil)
}

func TestFocusBlur(t *testing.T) {
	is := is.New(t)
	m := newTestModel()

	m, _ = m.Focus()
	is.True(m.IsFocused())

	m = m.Blur()
	is.True(!m.IsFocused())
}

func TestStoreResults(t *testing.T) {
	is := is.New(t)
	m := newTestModel()
	is.True(m.Results() == nil)

	m = m.StoreResults(nil)
	is.True(m.Results() == nil)
}

func TestHandleMsg_SearchClear(t *testing.T) {
	is := is.New(t)
	m, _ := newTestModel().Focus()

	// SearchClear key (esc)
	clearMsg := tea.KeyPressMsg{Code: tea.KeyEscape}
	m, cmd := m.HandleMsg(clearMsg)
	is.True(m.IsFocused()) // should stay focused after clear
	is.True(cmd == nil)
}

func TestHandleMsg_SearchSubmitEmptyQuery(t *testing.T) {
	is := is.New(t)
	m, _ := newTestModel().Focus()

	// SearchSubmit with empty query — should stay focused, no cmd
	submitMsg := tea.KeyPressMsg{Code: tea.KeyEnter}
	m, cmd := m.HandleMsg(submitMsg)
	is.True(m.IsFocused())
	is.True(cmd == nil)
}

func TestHandleMsg_SearchSubmitEmptyDoesNotBlur(t *testing.T) {
	is := is.New(t)
	m, _ := newTestModel().Focus()

	// Submit with no text — should stay focused
	submitMsg := tea.KeyPressMsg{Code: tea.KeyEnter}
	m, cmd := m.HandleMsg(submitMsg)
	is.True(m.IsFocused())
	is.True(cmd == nil)
}

func TestHelpBindings(t *testing.T) {
	is := is.New(t)
	m := newTestModel()
	bindings := m.HelpBindings()
	is.True(len(bindings) > 0)
}

func TestSetContext(t *testing.T) {
	is := is.New(t)
	m := newTestModel()
	m = m.SetContext(nil, "my-schema-id")
	// Just verify it compiles and doesn't panic
	is.True(!m.IsFocused())
}

func TestSetBaseDir(t *testing.T) {
	is := is.New(t)
	m := newTestModel()
	m = m.SetBaseDir("/tmp/indexes")
	// Just verify it compiles and doesn't panic
	is.True(!m.IsFocused())
}
