package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type windowBrowserKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Open  key.Binding
	Help  key.Binding
	Quit  key.Binding
	New   key.Binding
}

func (k windowBrowserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k windowBrowserKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.Help, k.Quit},
		{k.Open, k.New},
	}
}

var windowBrowserKeys = windowBrowserKeyMap{
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
		key.WithHelp("enter/space", "open"),
	),
	New: key.NewBinding(
		key.WithKeys("new", "n"),
		key.WithHelp("n", "create new window"),
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

type windowBrowserHelpModel struct {
	keys       windowBrowserKeyMap
	help       help.Model
	inputStyle lipgloss.Style
	quitting   bool
}

func newWindowBrowserHelpModel() windowBrowserHelpModel {
	return windowBrowserHelpModel{
		keys:       windowBrowserKeys,
		help:       help.New(),
		inputStyle: helpStyle,
	}
}

func (m windowBrowserHelpModel) Update(msg tea.Msg) (windowBrowserHelpModel, tea.Cmd) {
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

func (m windowBrowserHelpModel) View() string {
	if m.quitting {
		return ""
	}
	return m.help.View(m.keys)
}

func (m windowBrowserHelpModel) Init() tea.Cmd {
	return nil
}
