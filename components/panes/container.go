package panes

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	tea "github.com/charmbracelet/bubbletea"
)

type PaneContainerModel struct {
	model     common.TeamuxModel
	connector data.Connector
}

func NewPaneContainerModel(connector data.Connector, logger common.Logger) PaneContainerModel {
	return PaneContainerModel{
		model: NewPaneBrowserModel(connector, logger),
	}
}

func (m PaneContainerModel) Init() tea.Cmd {
	return tea.Batch(m.model.Init())
}

func (m PaneContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case common.NewPaneMsg, common.EditPMsg:
		m.model = NewPaneEditorModel(m.connector, m.model.GetLogger())
		return m, nil
	case common.PanesEditedMsg, common.BrowseMsg:
		m.model = NewPaneBrowserModel(m.connector, m.model.GetLogger())
		return m, common.Reaload
	}
	newModel, cmd := m.model.Update(msg)
	m.model = newModel
	return m, cmd
}

func (m PaneContainerModel) View() string {
	return m.model.View()
}
