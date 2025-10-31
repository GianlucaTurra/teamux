// Package sessions defines the UI components to perform tmux operations and
// CRUD operations with tmux sessions and saved session layout
package sessions

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	tea "github.com/charmbracelet/bubbletea"
)

type SessionContainerModel struct {
	focusedModel tea.Model
	connector    data.Connector
	logger       common.Logger
}

func NewSessionContainerModel(connector data.Connector, logger common.Logger) SessionContainerModel {
	return SessionContainerModel{
		focusedModel: NewSessionBrowserModel(connector, logger),
		connector:    connector,
		logger:       logger,
	}
}

func (m SessionContainerModel) Init() tea.Cmd {
	return tea.Batch(m.focusedModel.Init())
}

func (m SessionContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.NewSessionMsg:
		m.focusedModel = NewSessionEditorModel(m.connector, m.logger, nil)
		return m, nil
	case common.SessionCreatedMsg:
		m.focusedModel = NewSessionBrowserModel(m.connector, m.logger)
		return m, common.Reaload
	case common.BrowseMsg:
		m.focusedModel = NewSessionBrowserModel(m.connector, m.logger)
		return m, nil
	case common.EditSMsg:
		// FIXME: shouldn't the message pass down the session to the editor? AKA: no return?
		m.focusedModel = NewSessionEditorModel(m.connector, m.logger, &msg.Session)
		return m, nil
		// TODO: handle new session shortcut
		// case tea.KeyMsg:
		// 	if msg.String() == "n" && m.focusedModel == sessionBrowser && m.sessionBrowser.State != common.Deleting {
		// 		return m, common.NewSession
		// 	}
	}
	var cmds []tea.Cmd
	newModel, cmd := m.focusedModel.Update(msg)
	m.focusedModel = newModel
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m SessionContainerModel) View() string {
	return m.focusedModel.View()
}
