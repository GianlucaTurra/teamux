package components

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type browserKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Open   key.Binding
	Help   key.Binding
	Quit   key.Binding
	New    key.Binding
	Edit   key.Binding
	Delete key.Binding
}

func (k browserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k browserKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Open},
		{k.New, k.Edit, k.Delete},
	}
}

var browserKeys = browserKeyMap{
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
		key.WithHelp("enter/space", "open"),
	),
	New: key.NewBinding(
		key.WithKeys("new", "n"),
		key.WithHelp("n", "create new"),
	),
	Edit: key.NewBinding(
		key.WithKeys("edit", "e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("delete", "d"),
		key.WithHelp("d", "delete"),
	),
}

type BrowserHelpModel struct {
	keys     browserKeyMap
	Help     help.Model
	quitting bool
}

func NewBrowserHelpModel() common.HelpModel {
	return &BrowserHelpModel{
		keys: browserKeys,
		Help: help.New(),
	}
}

func (m BrowserHelpModel) ViewHelp() string {
	return m.Help.View(m.keys)
}

func (m *BrowserHelpModel) ToggleHelp() {
	m.Help.ShowAll = !m.Help.ShowAll
}

func (m *BrowserHelpModel) HideHelp() {
	m.Help.ShowAll = false
}

func (m BrowserHelpModel) Update(msg tea.Msg) (common.HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.Help.ShowAll = !m.Help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
		}
	}
	return &m, nil
}

func (m BrowserHelpModel) View() string {
	if m.quitting {
		return ""
	}
	return m.Help.View(m.keys)
}

func (m BrowserHelpModel) Init() tea.Cmd {
	return nil
}
