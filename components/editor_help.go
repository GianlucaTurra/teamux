package components

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type editorKeyMap struct {
	NextField     key.Binding
	PreviousField key.Binding
	Save          key.Binding
	BackToBrowser key.Binding
}

func (k editorKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k editorKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.BackToBrowser, k.Save},
		{k.NextField, k.PreviousField},
	}
}

var editorKeys = editorKeyMap{
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
		key.WithHelp("enter", "save"),
	),
	BackToBrowser: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to browser"),
	),
}

type EditorHelpModel struct {
	keys     editorKeyMap
	Help     help.Model
	quitting bool
}

func NewEditorHelpModel() common.HelpModel {
	return &EditorHelpModel{
		keys: editorKeys,
		Help: help.New(),
	}
}

func (m EditorHelpModel) ViewHelp() string {
	return m.Help.View(m.keys)
}

func (m *EditorHelpModel) ToggleHelp() {
	m.Help.ShowAll = !m.Help.ShowAll
}

func (m *EditorHelpModel) HideHelp() {
	m.Help.ShowAll = false
}

func (m EditorHelpModel) Update(msg tea.Msg) (common.HelpModel, tea.Cmd) {
	return &m, nil
}

func (m EditorHelpModel) View() string {
	if m.quitting {
		return ""
	}
	return m.Help.View(m.keys)
}

func (m EditorHelpModel) Init() tea.Cmd {
	return nil
}
