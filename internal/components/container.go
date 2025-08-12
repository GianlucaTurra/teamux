// Package components implements the UI for the application breaking it down
// into smaller pieces and assembling it in the container
package components

import (
	"database/sql"

	"github.com/GianlucaTurra/teamux/internal"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	sessionList  sessionListModel
	sessionInput sessionInputModel
	focusedModel int
	newPrefix    bool
}

const (
	sessionList = iota
	sessionInput
)

func InitialModel(db *sql.DB, logger internal.Logger) Model {
	return Model{
		newSessionListModel(db, logger),
		newSessionInputModel(db, logger),
		0,
		false,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.sessionList.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case internal.NewSessionMsg:
		m.focusedModel = sessionInput
		return m, nil
	case internal.SessionCreatedMsg:
		m.focusedModel = sessionList
		return m, internal.Reaload
	case internal.BrowseMsg:
		m.focusedModel = sessionList
		return m, nil
	case internal.EditMsg:
		m.focusedModel = sessionInput
	case tea.KeyMsg:
		if msg.String() == "n" && m.focusedModel == sessionList && m.sessionList.state != deleting {
			m.newPrefix = true
			return m, nil
		}
		if m.newPrefix {
			m.newPrefix = false
			switch msg.String() {
			case "s":
				return m, internal.NewSession
			}
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
	default:
		return lipgloss.JoinVertical(lipgloss.Left, m.sessionInput.View())
	}
}
