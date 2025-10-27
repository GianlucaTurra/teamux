package windows

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	tea "github.com/charmbracelet/bubbletea"
)

type WindowContainerModel struct {
	model     common.TeamuxModel
	connector data.Connector
}

func NewWindowContainerModel(connector data.Connector, logger common.Logger) WindowContainerModel {
	return WindowContainerModel{
		model:     NewWindowBrowserModel(connector, logger),
		connector: connector,
	}
}

func (m WindowContainerModel) Init() tea.Cmd {
	return tea.Batch(m.model.Init())
}

func (m WindowContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case common.NewWindowMsg, common.EditWMsg:
		m.model = NewWindowEditorModel(m.connector, m.model.GetLogger())
		return m, nil
	case common.WindowCreatedMsg:
		m.model = NewWindowBrowserModel(m.connector, m.model.GetLogger())
		return m, common.Reaload
	case common.BrowseMsg:
		m.model = NewWindowBrowserModel(m.connector, m.model.GetLogger())
		return m, nil
	}
	newModel, cmd := m.model.Update(msg)
	m.model = newModel
	return m, cmd
}

func (m WindowContainerModel) View() string {
	return m.model.View()
}
