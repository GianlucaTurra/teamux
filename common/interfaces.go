package common

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
)

type TeamuxModel interface {
	GetDB() *sql.DB
	GetLogger() Logger
	Init() tea.Cmd
	Update(tea.Msg) (TeamuxModel, tea.Cmd)
	View() string
}

type HelpModel interface {
	ViewHelp() string
	HideHelp()
	ToggleHelp()
	Init() tea.Cmd
	Update(tea.Msg) (HelpModel, tea.Cmd)
	View() string
}
