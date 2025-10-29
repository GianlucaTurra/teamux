package common

import (
	tea "github.com/charmbracelet/bubbletea"
)

type HelpModel interface {
	ViewHelp() string
	HideHelp()
	ToggleHelp()
	Init() tea.Cmd
	Update(tea.Msg) (HelpModel, tea.Cmd)
	View() string
}
