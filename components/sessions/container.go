// Package sessions defines the UI components to perform tmux operations and
// CRUD operations with tmux sessions and saved session layout
package sessions

import (
	"database/sql"

	"github.com/GianlucaTurra/teamux/common"
	tea "github.com/charmbracelet/bubbletea"
)

type SessionContainerModel struct {
	sessionBrowser SessionBrowserModel
	sessionEditor  SessionEditorModel
	focusedModel   int
}

const (
	sessionBrowser = iota
	sessionEditor
)

func NewSessionContainerModel(db *sql.DB, logger common.Logger) SessionContainerModel {
	return SessionContainerModel{
		sessionBrowser: NewSessionBrowserModel(db, logger),
		sessionEditor:  NewSessionEditorModel(db, logger),
		focusedModel:   0,
	}
}

func (m SessionContainerModel) Init() tea.Cmd {
	return tea.Batch(m.sessionBrowser.Init())
}

func (m SessionContainerModel) Update(msg tea.Msg) (SessionContainerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.NewSessionMsg:
		m.focusedModel = sessionEditor
		return m, nil
	case common.SessionCreatedMsg:
		m.focusedModel = sessionBrowser
		return m, common.Reaload
	case common.BrowseMsg:
		m.focusedModel = sessionBrowser
		return m, nil
	case common.EditMsg:
		m.focusedModel = sessionEditor
	case tea.KeyMsg:
		if msg.String() == "n" && m.focusedModel == sessionBrowser && m.sessionBrowser.State != common.Deleting {
			return m, common.NewSession
		}
	}
	var cmds []tea.Cmd
	switch m.focusedModel {
	case sessionEditor:
		newEditor, cmd := m.sessionEditor.Update(msg)
		m.sessionEditor = newEditor
		cmds = append(cmds, cmd)
	case sessionBrowser:
		newBrowser, cmd := m.sessionBrowser.Update(msg)
		m.sessionBrowser = newBrowser
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m SessionContainerModel) View() string {
	switch m.focusedModel {
	case sessionBrowser:
		return m.sessionBrowser.View()
	case sessionEditor:
		return m.sessionEditor.View()
	default:
		return ""
	}
}
