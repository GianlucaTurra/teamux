package panes

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type panesBrowserKeyMap struct {
	baseKeys common.BrowserKeyMap
}

func (k panesBrowserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k panesBrowserKeyMap) FullHelp() [][]key.Binding {
	var fullHelp [][]key.Binding
	fullHelp = append(fullHelp, k.baseKeys.FullHelp()...)
	return fullHelp
}

var panesBrowserKeys = panesBrowserKeyMap{
	baseKeys: common.GetBasicBrowserKeys(),
}

type PanesBrowserHelpModel struct {
	keys     panesBrowserKeyMap
	Help     help.Model
	quitting bool
}

func NewPanesBrowserHelpModel() common.HelpModel {
	return &PanesBrowserHelpModel{
		keys: panesBrowserKeys,
		Help: help.New(),
	}
}

func (m PanesBrowserHelpModel) ViewHelp() string {
	return m.Help.View(m.keys)
}

func (m *PanesBrowserHelpModel) ToggleHelp() {
	m.Help.ShowAll = !m.Help.ShowAll
}

func (m *PanesBrowserHelpModel) HideHelp() {
	m.Help.ShowAll = false
}

func (m PanesBrowserHelpModel) Update(msg tea.Msg) (common.HelpModel, tea.Cmd) {
	return &m, nil
}

func (m PanesBrowserHelpModel) View() string {
	if m.quitting {
		return ""
	}
	return m.Help.View(m.keys)
}

func (m PanesBrowserHelpModel) Init() tea.Cmd {
	return nil
}
