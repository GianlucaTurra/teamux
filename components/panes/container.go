package panes

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/database"
	tea "github.com/charmbracelet/bubbletea"
)

type PaneContainerModel struct {
	model     tea.Model
	connector database.Connector
}

func NewPaneContainerModel(connector database.Connector) PaneContainerModel {
	return PaneContainerModel{
		model:     NewPaneBrowserModel(connector),
		connector: connector,
	}
}

func (m PaneContainerModel) Init() tea.Cmd {
	return tea.Batch(m.model.Init())
}

func (m PaneContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.NewPaneMsg:
		m.model = NewPaneEditorModel(m.connector, nil)
		return m, nil
	case EditPMsg:
		m.model = NewPaneEditorModel(m.connector, &msg.Pane)
	case PaneCreatedMsg, PanesEditedMsg, common.BrowseMsg:
		m.model = NewPaneBrowserModel(m.connector)
		return m, common.Reaload
	}
	newModel, cmd := m.model.Update(msg)
	m.model = newModel
	return m, cmd
}

func (m PaneContainerModel) View() string {
	return m.model.View()
}
