package common

import (
	"github.com/charmbracelet/bubbles/key"
)

type BrowserKeyMap struct {
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

func (k BrowserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k BrowserKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Open},
		{k.New, k.Edit, k.Delete},
	}
}

var browserKeys = BrowserKeyMap{
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
		key.WithKeys("n"),
		key.WithHelp("n", "create new"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
}

func GetBasicBrowserKeys() BrowserKeyMap {
	return browserKeys
}
