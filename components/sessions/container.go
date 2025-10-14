// Package sessions defines the UI components to perform tmux operations and
// CRUD operations with tmux sessions and saved session layout
package sessions

import (
	"database/sql"

	"github.com/GianlucaTurra/teamux/common"
	tea "github.com/charmbracelet/bubbletea"
)

type SessionContainerModel struct {
	focusedModel common.TeamuxModel
}

func NewSessionContainerModel(db *sql.DB, logger common.Logger) SessionContainerModel {
	return SessionContainerModel{
		focusedModel: NewSessionBrowserModel(db, logger),
	}
}

func (m SessionContainerModel) Init() tea.Cmd {
	return tea.Batch(m.focusedModel.Init())
}

func (m SessionContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case common.NewSessionMsg:
		m.focusedModel = NewSessionEditorModel(m.focusedModel.GetDB(), m.focusedModel.GetLogger())
		return m, nil
	case common.SessionCreatedMsg:
		m.focusedModel = NewSessionBrowserModel(m.focusedModel.GetDB(), m.focusedModel.GetLogger())
		return m, common.Reaload
	case common.BrowseMsg:
		m.focusedModel = NewSessionBrowserModel(m.focusedModel.GetDB(), m.focusedModel.GetLogger())
		return m, nil
	case common.EditSMsg:
		m.focusedModel = NewSessionEditorModel(m.focusedModel.GetDB(), m.focusedModel.GetLogger())
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
