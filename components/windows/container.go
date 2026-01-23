package windows

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/database"
	tea "github.com/charmbracelet/bubbletea"
)

type WindowContainerModel struct {
	model     tea.Model
	connector database.Connector
	logger    common.Logger
}

func NewWindowContainerModel(connector database.Connector, logger common.Logger) WindowContainerModel {
	return WindowContainerModel{
		model:     NewWindowBrowserModel(connector, logger),
		connector: connector,
		logger:    logger,
	}
}

func (m WindowContainerModel) Init() tea.Cmd {
	return tea.Batch(m.model.Init())
}

func (m WindowContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.NewWindowMsg:
		m.model = NewWindowEditorModel(m.connector, m.logger, nil)
		return m, nil
	case EditWMsg:
		m.model = NewWindowEditorModel(m.connector, m.logger, &msg.Window)
	case WindowCreatedMsg:
		m.model = NewWindowBrowserModel(m.connector, m.logger)
		return m, common.Reaload
	case common.BrowseMsg:
		m.model = NewWindowBrowserModel(m.connector, m.logger)
		return m, nil
	case AssociatePanesMsg:
		m.model = NewWindowPanesAssociationModel(m.connector, m.logger, msg.Window)
		return m, common.LoadData
	}
	newModel, cmd := m.model.Update(msg)
	m.model = newModel
	return m, cmd
}

func (m WindowContainerModel) View() string {
	return m.model.View()
}
