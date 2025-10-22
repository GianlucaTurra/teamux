package windows

import (
	"database/sql"

	"github.com/GianlucaTurra/teamux/common"
	tea "github.com/charmbracelet/bubbletea"
)

type WindowContainerModel struct {
	model common.TeamuxModel
}

func NewWindowContainerModel(db *sql.DB, logger common.Logger) WindowContainerModel {
	return WindowContainerModel{
		model: NewWindowBrowserModel(db, logger),
	}
}

func (m WindowContainerModel) Init() tea.Cmd {
	return tea.Batch(m.model.Init())
}

func (m WindowContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case common.NewWindowMsg, common.EditWMsg:
		m.model = NewWindowEditorModel(m.model.GetDB(), m.model.GetLogger())
		return m, nil
	case common.WindowCreatedMsg:
		m.model = NewWindowBrowserModel(m.model.GetDB(), m.model.GetLogger())
		return m, common.Reaload
	case common.BrowseMsg:
		m.model = NewWindowBrowserModel(m.model.GetDB(), m.model.GetLogger())
		return m, nil
	}
	newModel, cmd := m.model.Update(msg)
	m.model = newModel
	return m, cmd
}

func (m WindowContainerModel) View() string {
	return m.model.View()
}
