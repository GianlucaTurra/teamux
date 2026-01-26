package sessions

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionBrowserKeyMap struct {
	baseKeys           common.BrowserKeyMap
	AddWindows         key.Binding
	ShowRelatedWindows key.Binding
}

func (k sessionBrowserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k sessionBrowserKeyMap) FullHelp() [][]key.Binding {
	var fullHelp [][]key.Binding
	fullHelp = append(fullHelp, k.baseKeys.FullHelp()...)
	sessionBrowserKeys := [][]key.Binding{
		{k.AddWindows},
		{k.ShowRelatedWindows},
	}
	fullHelp = append(fullHelp, sessionBrowserKeys...)
	return fullHelp
}

var sessionBrowserKeys = sessionBrowserKeyMap{
	baseKeys: common.GetBasicBrowserKeys(),
	AddWindows: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add windows"),
	),
	ShowRelatedWindows: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "show rel windows"),
	),
}

type SessionBrowserHelpModel struct {
	keys     sessionBrowserKeyMap
	Help     help.Model
	quitting bool
}

func NewSessionBrowserHelpModel() common.HelpModel {
	return &SessionBrowserHelpModel{
		keys: sessionBrowserKeys,
		Help: help.New(),
	}
}

func (m SessionBrowserHelpModel) ViewHelp() string {
	return m.Help.View(m.keys)
}

func (m *SessionBrowserHelpModel) ToggleHelp() {
	m.Help.ShowAll = !m.Help.ShowAll
}

func (m *SessionBrowserHelpModel) HideHelp() {
	m.Help.ShowAll = false
}

func (m SessionBrowserHelpModel) Update(msg tea.Msg) (common.HelpModel, tea.Cmd) {
	return &m, nil
}

func (m SessionBrowserHelpModel) View() string {
	if m.quitting {
		return ""
	}
	return m.Help.View(m.keys)
}

func (m SessionBrowserHelpModel) Init() tea.Cmd {
	return nil
}
