// Package components implements the UI for the application breaking it down
// into smaller pieces and assembling it in the container
package components

import (
	"github.com/GianlucaTurra/teamux/internal"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	sessionList  sessionListModel
	sessionInput sessionInputModel
	help         helpModel
	focusedModel int
}

const (
	sessionList = iota
	sessionInput
)

func InitialModel() Model {
	return Model{
		newSessionListModel(),
		newSessionInputModel(),
		newHelpModel(),
		0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.sessionList.Init(), m.help.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case internal.NewMsg:
		m.focusedModel = sessionInput
		return m, nil
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
	newHelp, cmd := m.help.Update(msg)
	m.help = newHelp
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.focusedModel {
	case sessionList:
		return lipgloss.JoinVertical(lipgloss.Left, m.sessionList.View(), m.help.View())
	default:
		return lipgloss.JoinVertical(lipgloss.Left, m.sessionInput.View(), m.help.View())
	}
}
