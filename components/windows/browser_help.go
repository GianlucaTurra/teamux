package windows

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type windowsBrowserKeyMap struct {
	baseKeys         common.BrowserKeyMap
	AddPanes         key.Binding
	ShowRelatedPanes key.Binding
}

func (k windowsBrowserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k windowsBrowserKeyMap) FullHelp() [][]key.Binding {
	var fullHelp [][]key.Binding
	fullHelp = append(fullHelp, k.baseKeys.FullHelp()...)
	windowsBrowserKeys := [][]key.Binding{
		{k.AddPanes},
		{k.ShowRelatedPanes},
	}
	fullHelp = append(fullHelp, windowsBrowserKeys...)
	return fullHelp
}

var windowsBrowserKeys = windowsBrowserKeyMap{
	baseKeys: common.GetBasicBrowserKeys(),
	AddPanes: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add panes"),
	),
	ShowRelatedPanes: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "show rel panes"),
	),
}

type WindowsBrowserHelpModel struct {
	keys     windowsBrowserKeyMap
	Help     help.Model
	quitting bool
}

func NewWindowsBrowserHelpModel() common.HelpModel {
	return &WindowsBrowserHelpModel{
		keys: windowsBrowserKeys,
		Help: help.New(),
	}
}

func (m WindowsBrowserHelpModel) ViewHelp() string {
	return m.Help.View(m.keys)
}

func (m *WindowsBrowserHelpModel) ToggleHelp() {
	m.Help.ShowAll = !m.Help.ShowAll
}

func (m *WindowsBrowserHelpModel) HideHelp() {
	m.Help.ShowAll = false
}

func (m WindowsBrowserHelpModel) Update(msg tea.Msg) (common.HelpModel, tea.Cmd) {
	return &m, nil
}

func (m WindowsBrowserHelpModel) View() string {
	if m.quitting {
		return ""
	}
	return m.Help.View(m.keys)
}

func (m WindowsBrowserHelpModel) Init() tea.Cmd {
	return nil
}
