// Package sessions defines the UI components to perform tmux operations and
// CRUD operations with tmux sessions and saved session layout
package sessions

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/windows"
	"github.com/GianlucaTurra/teamux/database"
	tea "github.com/charmbracelet/bubbletea"
)

type SessionContainerModel struct {
	focusedModel tea.Model
	connector    database.Connector
}

func NewSessionContainerModel(connector database.Connector) SessionContainerModel {
	return SessionContainerModel{
		focusedModel: NewSessionBrowserModel(connector),
		connector:    connector,
	}
}

func (m SessionContainerModel) Init() tea.Cmd {
	return tea.Batch(m.focusedModel.Init())
}

func (m SessionContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case NewSessionMsg:
		m.focusedModel = NewSessionEditorModel(m.connector, nil)
		return m, nil
	case SessionCreatedMsg:
		m.focusedModel = NewSessionBrowserModel(m.connector)
		return m, common.Reaload
	case common.BrowseMsg:
		m.focusedModel = NewSessionBrowserModel(m.connector)
		return m, nil
	case EditSMsg:
		// FIXME: shouldn't the message pass down the session to the editor? AKA: no return?
		m.focusedModel = NewSessionEditorModel(m.connector, &msg.Session)
		return m, nil
		// TODO: handle new session shortcut
		// case tea.KeyMsg:
		// 	if msg.String() == "n" && m.focusedModel == sessionBrowser && m.sessionBrowser.State != common.Deleting {
		// 		return m, common.NewSession
		// 	}
	case AssociateWindowsMsg:
		m.focusedModel = NewSessionWindowsAssociationModel(m.connector, msg.Session)
		return m, common.LoadData
	case common.CreateWindowMsg:
		m.focusedModel = windows.NewWindowEditorModel(m.connector, nil)
		return m, nil
	case windows.EditWindowMsg:
		m.focusedModel = windows.NewWindowEditorModel(m.connector, &msg.Window)
		return m, nil
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
