package components

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type State int

const (
	browsing State = iota
	deleting
	quitting
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedStyle     = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("140"))
	paginationStyle   = list.DefaultStyles().TitleBar.PaddingLeft(4)
	openStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).PaddingLeft(2)
	selectedOpenStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("200"))
)
