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
	tabs        []string
	mainModel   tea.Model
	detailModel tea.Model
	messageBox  MessageBoxModel
	helpModel   HelpModel
	focusedTab  int
	newPrefix   bool
	quitting    bool
	db          *sql.DB
	logger      common.Logger
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

func InitialModel(db *sql.DB, logger common.Logger) Model {
	return Model{
		tabs:        []string{SESSIONS, WINDOWS, PANES},
		mainModel:   sessions.NewSessionContainerModel(db, logger),
		detailModel: sessions.NewSessionTreeModel(db, logger, nil),
		messageBox:  NewMessageBoxModel(),
		helpModel:   NewHelpModel(),
		focusedTab:  0,
		newPrefix:   false,
		quitting:    false,
		db:          db,
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
		return m, common.ClearHelp
	case common.PreviousTabMsg:
		if m.focusedTab == 0 {
			m.focusedTab = len(m.tabs) - 1
		} else {
			m.focusedTab -= 1
		}
		return m, common.ClearHelp
	case common.NewSFocus:
		m.detailModel = sessions.NewSessionTreeModel(m.db, m.logger, &msg.Session)
		return m, nil
	case common.NewWFocus:
		m.detailModel = windows.NewWindowDetailModel(m.db, m.logger, &msg.Window)
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
	case common.ClearHelpMsg:
		switch m.focusedTab {
		case sessionsContainer:
			m.mainModel = sessions.NewSessionContainerModel(m.db, m.logger)
			m.detailModel = sessions.NewSessionTreeModel(m.db, m.logger, nil)
		case windwowBrowser:
			m.mainModel = windows.NewWindowContainerModel(m.db, m.logger)
			m.detailModel = windows.NewWindowDetailModel(m.db, m.logger, nil)
		case paneContainer:
			m.mainModel = panes.NewPaneContainerModel(m.db, m.logger)
			// TODO: add missing detail model
			// m.detailModel = panes.NewPaneDetailModel(m.db, m.logger, nil
		}
	}
	var cmds []tea.Cmd
	newMain, cmd := m.mainModel.Update(msg)
	m.mainModel = newMain
	cmds = append(cmds, cmd)
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
	mainView := m.mainModel.View()
	detailView := m.detailModel.View()
	left := lipgloss.JoinVertical(
		lipgloss.Left,
		common.TitleStyle.PaddingLeft(2).Render(tabHeader.String()),
		mainView,
	)
	right := lipgloss.NewStyle().Render(detailView)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, left, right),
		m.messageBox.View(),
		m.helpModel.View(),
	)
}
