package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/tonysyu/gqlxp/tui/adapters"
	"github.com/tonysyu/gqlxp/utils/testx"
	"github.com/tonysyu/gqlxp/utils/text"
)

var (
	keyNextPanel = tea.KeyMsg{Type: tea.KeyTab}
	keyPrevPanel = tea.KeyMsg{Type: tea.KeyShiftTab}
	keyNextType  = tea.KeyMsg{Type: tea.KeyCtrlT}
	keyPrevType  = tea.KeyMsg{Type: tea.KeyCtrlR}
)

func createTestSchema() []byte {
	return []byte(`
		type Object1 {
			field1: String!
			field2: String!
			field3: String!
		}

		type Object2 {
			field1: String!
			field2: String!
			field3: String!
		}

		type Query {
			query1: Object1
			query2: Object2
		}

		type Mutation {
			mutation1: Object1
			mutation2: Object2
		}
	`)
}

// Wrapper for mainModel used to simplify testing
type testModel struct {
	model mainModel
}

func newTestModel() testModel {
	schemaView, _ := adapters.ParseSchema(createTestSchema())

	model := newModel(schemaView)
	model.width = 120
	model.height = 40

	return testModel{model: model}
}

func (tm *testModel) Update(msg tea.Msg) {
	model := tm.model
	updatedModel, cmd := model.Update(msg)
	tm.model = updatedModel.(mainModel)
	if cmd != nil {
		// If Update returns a cmd, execute it and recursively Update model
		tm.Update(cmd())
	}
}

// ViewBreadcrumbs returns line of text from rendered view that represents breadcrumbs
func (tm *testModel) ViewBreadcrumbs() string {
	view := tm.model.View()
	lines := text.SplitLines(view)
	// NOTE: This index for the breadcrumb line is hard-coded and will change if layout changes
	return testx.NormalizeView(lines[2])
}

func TestBreadcrumbs(t *testing.T) {
	is := is.New(t)

	model := newTestModel()

	t.Run("initial state has no breadcrumbs", func(t *testing.T) {
		is.Equal(model.ViewBreadcrumbs(), "")
	})

	t.Run("navigating forward adds query breadcrumb", func(t *testing.T) {
		model.Update(keyNextPanel)
		is.Equal(model.ViewBreadcrumbs(), "query1")
	})

	t.Run("navigating forward again adds result type breadcrumb", func(t *testing.T) {
		model.Update(keyNextPanel)
		is.Equal(model.ViewBreadcrumbs(), "query1 > Object1")
	})

	t.Run("navigating forward removes result type breadcrumb", func(t *testing.T) {
		model.Update(keyPrevPanel)
		is.Equal(model.ViewBreadcrumbs(), "query1")
	})

	t.Run("navigating backward again removes query breadcrumb", func(t *testing.T) {
		model.Update(keyPrevPanel)
		is.Equal(model.ViewBreadcrumbs(), "")
	})

	t.Run("switching GQL type resets breadcrumbs", func(t *testing.T) {
		model.Update(keyNextPanel)
		is.Equal(model.ViewBreadcrumbs(), "query1")
		model.Update(keyNextType)
		is.Equal(model.ViewBreadcrumbs(), "")
	})

	t.Run("reverse GQL type also resets breadcrumbs", func(t *testing.T) {
		model.Update(keyNextPanel)
		// Previous test puts us on the Mutation tab, and the first mutation is createPost
		is.Equal(model.ViewBreadcrumbs(), "mutation1")
		model.Update(keyNextType)
		is.Equal(model.ViewBreadcrumbs(), "")
	})
}
