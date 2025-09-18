package sessions

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionEditorKeyMap struct {
	NextField     key.Binding
	PreviousField key.Binding
	Save          key.Binding
	BackToBrowser key.Binding
}

func (k sessionEditorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k sessionEditorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.BackToBrowser, k.Save},
		{k.NextField, k.PreviousField},
	}
}

var sessionEditorKeys = sessionEditorKeyMap{
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

type SessionEditorHelpModel struct {
	keys     sessionEditorKeyMap
	Help     help.Model
	quitting bool
}

func NewSessionEditorHelpModel() SessionEditorHelpModel {
	help := help.New()
	help.ShowAll = true
	return SessionEditorHelpModel{
		keys: sessionEditorKeys,
		Help: help,
	}
}

func (m SessionEditorHelpModel) ViewHelp() string {
	return m.Help.View(m.keys)
}

func (m *SessionEditorHelpModel) ToggleHelp() {
	m.Help.ShowAll = !m.Help.ShowAll
}

func (m *SessionEditorHelpModel) HideHelp() {
	m.Help.ShowAll = false
}

func (m SessionEditorHelpModel) Update(msg tea.Msg) (SessionEditorHelpModel, tea.Cmd) {
	return m, nil
}

func (m SessionEditorHelpModel) View() string {
	if m.quitting {
		return ""
	}
	return m.Help.View(m.keys)
}

func (m SessionEditorHelpModel) Init() tea.Cmd {
	return nil
}
