// Package components implements the UI for the application breaking it down
// into smaller pieces and assembling it in the container
package components

import (
	"database/sql"
	"strings"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/panes"
	"github.com/GianlucaTurra/teamux/components/sessions"
	"github.com/GianlucaTurra/teamux/components/windows"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	tabs             []string
	sessionContainer sessions.SessionContainerModel
	windowContainer  windows.WindowContainerModel
	paneContainer    panes.PaneContainerModel
	sessionDetail    sessions.SessionDetailModel
	windowDetail     windows.WindowDetailModel
	messageBox       MessageBoxModel
	helpModel        HelpModel
	focusedTab       int
	newPrefix        bool
	quitting         bool
	db               *sql.DB
	logger           common.Logger
}

const (
	SESSIONS = "Sessions"
	WINDOWS  = "Windows"
	PANES    = "Panes"
)

const (
	sessionsContainer = iota
	windwowBrowser
	paneContainer
)

const (
	sessionDetail = iota
	windowDetail
	paneDetail
)

func InitialModel(db *sql.DB, logger common.Logger) Model {
	return Model{
		tabs:             []string{SESSIONS, WINDOWS, PANES},
		sessionContainer: sessions.NewSessionContainerModel(db, logger),
		windowContainer:  windows.NewWindowContainerModel(db, logger),
		paneContainer:    panes.NewPaneContainerModel(db, logger),
		sessionDetail:    sessions.NewSessionTreeModel(db, logger, nil),
		windowDetail:     windows.NewWindowDetailModel(db, logger, nil),
		messageBox:       NewMessageBoxModel(),
		helpModel:        NewHelpModel(),
		focusedTab:       0,
		newPrefix:        false,
		quitting:         false,
		db:               db,
		logger:           logger,
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
		return m, common.ClearHelp
	case common.PreviousTabMsg:
		if m.focusedTab == 0 {
			m.focusedTab = len(m.tabs) - 1
		} else {
			m.focusedTab -= 1
		}
		return m, common.ClearHelp
	case common.NewSFocus:
		m.sessionDetail = sessions.NewSessionTreeModel(m.db, m.logger, &msg.Session)
		return m, nil
	case common.NewWFocus:
		m.windowDetail = windows.NewWindowDetailModel(m.db, m.logger, &msg.Window)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "]" {
			return m, common.NextTab
		}
		if msg.String() == "[" {
			return m, common.PreviousTab
		}
	case common.OutputMsg:
		m.messageBox, _ = m.messageBox.Update(msg)
		return m, nil
	}
	var cmds []tea.Cmd
	switch m.focusedTab {
	case sessionsContainer:
		newInput, cmd := m.sessionContainer.Update(msg)
		m.sessionContainer = newInput
		cmds = append(cmds, cmd)
	case windwowBrowser:
		newList, cmd := m.windowContainer.Update(msg)
		m.windowContainer = newList
		cmds = append(cmds, cmd)
	case paneContainer:
		newPanes, cmd := m.paneContainer.Update(msg)
		m.paneContainer = newPanes
		cmds = append(cmds, cmd)
	}
	newHelp, cmd := m.helpModel.Update(msg)
	m.helpModel = newHelp
	cmds = append(cmds, cmd)
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
	var currentDetail string
	switch m.focusedTab {
	case sessionsContainer:
		focusedView = m.sessionContainer.View()
		currentDetail = m.sessionDetail.View()
	case windwowBrowser:
		focusedView = m.windowContainer.View()
		currentDetail = m.windowDetail.View()
	case paneContainer:
		focusedView = m.paneContainer.View()
	}
	left := lipgloss.JoinVertical(
		lipgloss.Left,
		common.TitleStyle.PaddingLeft(2).Render(tabHeader.String()),
		focusedView,
	)
	right := lipgloss.NewStyle().Render(currentDetail)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, left, right),
		m.messageBox.View(),
		m.helpModel.View(),
	)
}
