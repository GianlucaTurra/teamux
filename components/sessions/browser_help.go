package sessions

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionBrowserKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Open   key.Binding
	Kill   key.Binding
	Switch key.Binding
	Help   key.Binding
	Quit   key.Binding
	New    key.Binding
}

func (k sessionBrowserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k sessionBrowserKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.Help, k.Quit},
		{k.Open, k.Switch, k.Kill, k.New},
	}
}

var keys = sessionBrowserKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Open: key.NewBinding(
		key.WithKeys("space", "enter"),
		key.WithHelp("enter", "open or switch if open"),
	),
	Switch: key.NewBinding(
		key.WithKeys("s", "switch"),
		key.WithHelp("s", "switch to session"),
	),
	Kill: key.NewBinding(
		key.WithKeys("kill", "d"),
		key.WithHelp("d", "kill open session"),
	),
	New: key.NewBinding(
		key.WithKeys("new", "n"),
		key.WithHelp("n", "create new session"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type sessionBrowserHelpModel struct {
	keys       sessionBrowserKeyMap
	help       help.Model
	inputStyle lipgloss.Style
	quitting   bool
}

func newSessionBrowserHelpModel() sessionBrowserHelpModel {
	return sessionBrowserHelpModel{
		keys:       keys,
		help:       help.New(),
		inputStyle: common.HelpStyle,
	}
}

func (m sessionBrowserHelpModel) Update(msg tea.Msg) (sessionBrowserHelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
		}
	}
	return m, nil
}

func (m sessionBrowserHelpModel) View() string {
	if m.quitting {
		return ""
	}
	return m.help.View(m.keys)
}

func (m sessionBrowserHelpModel) Init() tea.Cmd {
	return nil
}
