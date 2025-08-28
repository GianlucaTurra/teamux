package common

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle        = lipgloss.NewStyle().Bold(true).MarginBottom(1).MarginTop(1)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	SelectedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	PaginationStyle   = list.DefaultStyles().TitleBar.PaddingLeft(4)
	OpenStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).PaddingLeft(2)
	SelectedOpenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	HelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	FocusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	BlurredStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
)
