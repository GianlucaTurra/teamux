// Package components implements the UI for the application breaking it down
// into smaller pieces and assembling it in the container
package components

import (
	"strings"
	"time"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/panes"
	"github.com/GianlucaTurra/teamux/components/sessions"
	"github.com/GianlucaTurra/teamux/components/windows"
	"github.com/GianlucaTurra/teamux/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	tabs        []string
	mainModel   tea.Model
	detailModel tea.Model
	messageBox  tea.Model
	helpModel   HelpModel
	focusedTab  int
	newPrefix   bool
	quitting    bool
	connector   database.Connector
	logger      common.Logger
}

const (
	SESSIONS = "Sessions"
	WINDOWS  = "Windows"
	PANES    = "Panes"
)

func InitialModel(connector database.Connector, logger common.Logger) Model {
	return Model{
		tabs:        []string{SESSIONS, WINDOWS, PANES},
		mainModel:   sessions.NewSessionContainerModel(connector, logger),
		detailModel: sessions.NewSessionTreeModel(connector, logger, nil),
		messageBox:  NewMessageBoxModel(),
		helpModel:   NewHelpModel(),
		focusedTab:  0,
		newPrefix:   false,
		quitting:    false,
		connector:   connector,
		logger:      logger,
	}
}

func (m Model) Init() tea.Cmd {
	return m.mainModel.Init()
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
		return m, func() tea.Msg { return common.ClearHelp(common.FocusedTab(m.focusedTab)) }
	case common.PreviousTabMsg:
		if m.focusedTab == 0 {
			m.focusedTab = len(m.tabs) - 1
		} else {
			m.focusedTab -= 1
		}
		return m, func() tea.Msg { return common.ClearHelp(common.FocusedTab(m.focusedTab)) }
	case sessions.NewSFocus:
		m.detailModel = sessions.NewSessionTreeModel(m.connector, m.logger, &msg.Session)
		return m, nil
	case windows.NewWFocus:
		m.detailModel = windows.NewWindowDetailModel(m.connector, m.logger, &msg.Window)
		return m, nil
	case panes.NewPFocusMsg:
		m.detailModel = panes.NewPaneDetailModel(m.connector, m.logger, &msg.Pane)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "]" {
			return m, common.NextTab
		}
		if msg.String() == "[" {
			return m, common.PreviousTab
		}
	case common.SetOutputMsgTimerMsg:
		// FIXME: doesn't actually work as expected
		return m, func() tea.Msg {
			time.Sleep(2 * time.Second)
			return common.ResetOutputMsgMsg{}
		}
	case sessions.EditSMsg, windows.EditWMsg, panes.EditPMsg:
		m.helpModel.SetModel(NewEditorHelpModel())
	case common.ClearHelpMsg:
		switch m.focusedTab {
		case common.SessionsContainer:
			m.mainModel = sessions.NewSessionContainerModel(m.connector, m.logger)
			m.detailModel = sessions.NewSessionTreeModel(m.connector, m.logger, nil)
		case common.WindwowBrowser:
			m.mainModel = windows.NewWindowContainerModel(m.connector, m.logger)
			m.detailModel = windows.NewWindowDetailModel(m.connector, m.logger, nil)
		case common.PaneContainer:
			m.mainModel = panes.NewPaneContainerModel(m.connector, m.logger)
			m.detailModel = panes.NewPaneDetailModel(m.connector, m.logger, nil)
		}
	}
	var cmds []tea.Cmd
	newMain, cmd := m.mainModel.Update(msg)
	m.mainModel = newMain
	cmds = append(cmds, cmd)
	newHelp, cmd := m.helpModel.Update(msg)
	m.helpModel = newHelp
	cmds = append(cmds, cmd)
	newMsgBox, cmd := m.messageBox.Update(msg)
	m.messageBox = newMsgBox
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
	left := lipgloss.JoinVertical(
		lipgloss.Left,
		common.TitleStyle.PaddingLeft(2).Render(tabHeader.String()),
		m.mainModel.View(),
	)
	right := lipgloss.NewStyle().Render(m.detailModel.View())
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, left, right),
		m.messageBox.View(),
		m.helpModel.View(),
	)
}
