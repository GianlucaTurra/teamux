package panes

import (
	"database/sql"

	"github.com/GianlucaTurra/teamux/common"
	tea "github.com/charmbracelet/bubbletea"
)

type PaneContainerModel struct {
	model common.TeamuxModel
}

func NewPaneContainerModel(db *sql.DB, logger common.Logger) PaneContainerModel {
	return PaneContainerModel{
		model: NewPaneBrowserModel(db, logger),
	}
}

func (m PaneContainerModel) Init() tea.Cmd {
	return tea.Batch(m.model.Init())
}

func (m PaneContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case common.NewPaneMsg, common.EditPMsg:
		m.model = NewPaneEditorModel(m.model.GetDB(), m.model.GetLogger())
		return m, nil
	case common.PaneCreatedMsg, common.BrowseMsg:
		m.model = NewPaneBrowserModel(m.model.GetDB(), m.model.GetLogger())
		return m, common.Reaload
	}
	newModel, cmd := m.model.Update(msg)
	m.model = newModel
	return m, cmd
}

func (m PaneContainerModel) View() string {
	return m.model.View()
}
