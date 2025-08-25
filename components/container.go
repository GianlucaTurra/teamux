// Package components implements the UI for the application breaking it down
// into smaller pieces and assembling it in the container
package components

import (
	"database/sql"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/sessions"
	"github.com/GianlucaTurra/teamux/components/windows"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	sessionList   sessions.SessionBrowserModel
	sessionInput  sessions.SessionEditorModel
	windowBrowser windows.WindowBrowserModel
	focusedModel  int
	newPrefix     bool
}

const (
	sessionList = iota
	sessionInput
	windwowBrowser
)

func InitialModel(db *sql.DB, logger common.Logger) Model {
	return Model{
		sessions.NewSessionBrowserModel(db, logger),
		sessions.NewSessionEditorModel(db, logger),
		windows.NewWindowBrowserModel(db, logger),
		0,
		false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.sessionList.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.NewSessionMsg:
		m.focusedModel = sessionInput
		return m, nil
	case common.SessionCreatedMsg:
		m.focusedModel = sessionList
		return m, common.Reaload
	case common.BrowseMsg:
		m.focusedModel = sessionList
		return m, nil
	case common.EditMsg:
		m.focusedModel = sessionInput
	case tea.KeyMsg:
		if msg.String() == "n" && m.focusedModel == sessionList && m.sessionList.State != common.Deleting {
			m.newPrefix = true
			return m, nil
		}
		if m.newPrefix {
			m.newPrefix = false
			switch msg.String() {
			case "s":
				return m, common.NewSession
			}
		}
		if msg.String() == "b" && m.focusedModel == sessionList {
			m.focusedModel = windwowBrowser
			return m, nil
		}
	}
	var cmds []tea.Cmd
	switch m.focusedModel {
	case sessionInput:
		newInput, cmd := m.sessionInput.Update(msg)
		m.sessionInput = newInput
		cmds = append(cmds, cmd)
	case sessionList:
		newList, cmd := m.sessionList.Update(msg)
		m.sessionList = newList
		cmds = append(cmds, cmd)
	case windwowBrowser:
		newList, cmd := m.windowBrowser.Update(msg)
		m.windowBrowser = newList
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.focusedModel {
	case sessionList:
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.sessionList.View(),
		)
	case windwowBrowser:
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.windowBrowser.View(),
		)
	default:
		return lipgloss.JoinVertical(lipgloss.Left, m.sessionInput.View())
	}
}
