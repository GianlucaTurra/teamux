// Package components implements the UI for the application breaking it down
// into smaller pieces and assembling it in the container
package components

import (
	"database/sql"
	"strings"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/sessions"
	"github.com/GianlucaTurra/teamux/components/windows"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	tabs             []string
	sessionContainer sessions.SessionContainerModel
	windowBrowser    windows.WindowBrowserModel
	focusedTab       int
	newPrefix        bool
	quitting         bool
}

const (
	SESSIONS = "Sessions"
	WINDOWS  = "Windows"
)

const (
	sessionsContainer = iota
	windwowBrowser
)

func InitialModel(db *sql.DB, logger common.Logger) Model {
	return Model{
		[]string{SESSIONS, WINDOWS},
		sessions.NewSessionContainerModel(db, logger),
		windows.NewWindowBrowserModel(db, logger),
		0,
		false,
		false,
	}
}

func (m Model) Init() tea.Cmd {
	return m.sessionContainer.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.QuitMsg:
		m.quitting = true
		return m, tea.Quit
	case common.NextTabMsg:
		if m.focusedTab == len(m.tabs)-1 {
			m.focusedTab = 0
		} else {
			m.focusedTab += 1
		}
		return m, nil
	case common.PreviousTabMsg:
		if m.focusedTab == 0 {
			m.focusedTab = len(m.tabs) - 1
		} else {
			m.focusedTab -= 1
		}
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "]" {
			return m, common.NextTab
		}
		if msg.String() == "[" {
			return m, common.PreviousTab
		}
	}
	var cmds []tea.Cmd
	switch m.focusedTab {
	case sessionsContainer:
		newInput, cmd := m.sessionContainer.Update(msg)
		m.sessionContainer = newInput
		cmds = append(cmds, cmd)
	case windwowBrowser:
		newList, cmd := m.windowBrowser.Update(msg)
		m.windowBrowser = newList
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	tabHeader := strings.Builder{}
	separator := " "
	for i, t := range m.tabs {
		if i == m.focusedTab {
			tabHeader.WriteString(common.FocusedStyle.Render("*" + t))
		} else {
			tabHeader.WriteString(common.BlurredStyle.Render(t))
		}
		tabHeader.WriteString(separator)
	}
	var focusedView string
	switch m.focusedTab {
	case sessionsContainer:
		focusedView = m.sessionContainer.View()
	case windwowBrowser:
		focusedView = m.windowBrowser.View()
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		common.TitleStyle.PaddingLeft(2).Render(tabHeader.String()),
		focusedView,
	)
}
