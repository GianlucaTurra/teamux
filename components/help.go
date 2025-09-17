package components

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/sessions"
	"github.com/GianlucaTurra/teamux/components/windows"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type basicHelpKeymap struct {
	help key.Binding
	quit key.Binding
}

func (k basicHelpKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.help, k.quit}
}

func (k basicHelpKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.help}, {k.quit},
	}
}

var basicHelpKeys = basicHelpKeymap{
	help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "show help")),
	quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
}

type HelpModel struct {
	basicKeys          basicHelpKeymap
	help               help.Model
	displayedHelp      common.ComponentWithHelp
	sessionBrowserHelp sessions.SessionBrowserHelpModel
	sessionEditorHelp  sessions.SessionEditorHelpModel
	windowBrowserHelp  windows.WindowBrowserHelpModel
}

func NewHelpModel() HelpModel {
	return HelpModel{
		basicKeys:          basicHelpKeys,
		help:               help.New(),
		displayedHelp:      common.SessionBrowser,
		sessionBrowserHelp: sessions.NewSessionBrowserHelpModel(),
		sessionEditorHelp:  sessions.NewSessionEditorHelpModel(),
		windowBrowserHelp:  windows.NewWindowBrowserHelpModel(),
	}
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.ShowFullHelpMsg:
		m.displayedHelp = msg.Component
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.basicKeys.help):
			switch m.displayedHelp {
			case common.SessionBrowser:
				newSessionBrowserHelp := m.sessionBrowserHelp
				newSessionBrowserHelp.Help.ShowAll = !newSessionBrowserHelp.Help.ShowAll
				m.sessionBrowserHelp = newSessionBrowserHelp
			// case common.SessionEditor:
			// 	newSessionEditorHelp := m.sessionEditorHelp
			// 	newSessionEditorHelp.Help.ShowAll = !newSessionEditorHelp.Help.ShowAll
			// 	m.sessionEditorHelp = newSessionEditorHelp
			case common.WindowBrowser:
				newWindowBrowserHelp := m.windowBrowserHelp
				newWindowBrowserHelp.Help.ShowAll = !newWindowBrowserHelp.Help.ShowAll
				m.windowBrowserHelp = newWindowBrowserHelp
			}
			return m, nil
		default:
			return m, nil
		}
	}
	return m, nil
}

func (m HelpModel) View() string {
	var fullHelp string
	switch m.displayedHelp {
	case common.SessionBrowser:
		fullHelp = m.sessionBrowserHelp.ViewHelp()
	case common.WindowBrowser:
		fullHelp = m.windowBrowserHelp.ViewHelp()
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.help.View(m.basicKeys),
		fullHelp,
	)
}
