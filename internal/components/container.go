// Package components implements the UI for the application breaking it down
// into smaller pieces and assembling it in the container
package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	sessionList sessionListModel
	help        helpModel
}

func InitialModel() Model {
	return Model{newSessionListModel(), newHelpModel()}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.sessionList.Init(), m.help.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	newList, cmd := m.sessionList.Update(msg)
	m.sessionList = newList
	cmds = append(cmds, cmd)
	newHelp, cmd := m.help.Update(msg)
	m.help = newHelp
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.sessionList.View(), m.help.View())
}
