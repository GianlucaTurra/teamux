package common

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle        = lipgloss.NewStyle().MarginLeft(2)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedStyle     = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("140"))
	PaginationStyle   = list.DefaultStyles().TitleBar.PaddingLeft(4)
	OpenStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).PaddingLeft(2)
	SelectedOpenStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("200"))
	HelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)
