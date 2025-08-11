package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionEditorKeyMap struct {
	Help          key.Binding
	Quit          key.Binding
	NextField     key.Binding
	PreviousField key.Binding
	Save          key.Binding
	BackToBrowser key.Binding
}

func (k sessionEditorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.BackToBrowser}
}

func (k sessionEditorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit, k.BackToBrowser},
		{k.NextField, k.PreviousField, k.Save},
	}
}

var sessionEditorKeys = sessionEditorKeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "show help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	NextField: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
	PreviousField: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "previous field"),
	),
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "save session"),
	),
	BackToBrowser: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to browser"),
	),
}

type sessionEditorHelpModel struct {
	keys       sessionEditorKeyMap
	help       help.Model
	inputStyle lipgloss.Style
	quitting   bool
}

func newSessionEditorHelpModel() sessionEditorHelpModel {
	return sessionEditorHelpModel{
		keys:       sessionEditorKeys,
		help:       help.New(),
		inputStyle: helpStyle,
	}
}

func (m sessionEditorHelpModel) Update(msg tea.Msg) (sessionEditorHelpModel, tea.Cmd) {
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

func (m sessionEditorHelpModel) View() string {
	if m.quitting {
		return ""
	}
	return m.help.View(m.keys)
}

func (m sessionEditorHelpModel) Init() tea.Cmd {
	return nil
}
