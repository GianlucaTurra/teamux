package sessions

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type SessionBrowserKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Open   key.Binding
	Kill   key.Binding
	Switch key.Binding
	New    key.Binding
	Edit   key.Binding
	Delete key.Binding
}

func (k SessionBrowserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k SessionBrowserKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.New, k.Open},
		{k.Switch, k.Kill, k.Edit, k.Delete},
	}
}

var keys = SessionBrowserKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
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
		key.WithKeys("kill", "K"),
		key.WithHelp("K", "kill open session"),
	),
	Edit: key.NewBinding(
		key.WithKeys("edit", "e"),
		key.WithHelp("e", "edit session"),
	),
	New: key.NewBinding(
		key.WithKeys("new", "n"),
		key.WithHelp("n", "create new session"),
	),
	Delete: key.NewBinding(
		key.WithKeys("delete", "d"),
		key.WithHelp("d", "delete session"),
	),
}

type SessionBrowserHelpModel struct {
	keys     SessionBrowserKeyMap
	Help     help.Model
	quitting bool
}

func NewSessionBrowserHelpModel() SessionBrowserHelpModel {
	return SessionBrowserHelpModel{
		keys: keys,
		Help: help.New(),
	}
}

func (m SessionBrowserHelpModel) ViewHelp() string {
	return m.Help.View(keys)
}

func (m SessionBrowserHelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
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
