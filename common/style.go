package common

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	Black   = "0"
	Red     = "1"
	Green   = "2"
	Yellow  = "3"
	Blue    = "4"
	Magenta = "5"
	Cyan    = "6"
	White   = "7"
)

var (
	HeaderStyle       = lipgloss.NewStyle().Italic(true)
	TitleStyle        = lipgloss.NewStyle().MarginBottom(1).MarginTop(1)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	SelectedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(Yellow))
	PaginationStyle   = list.DefaultStyles().TitleBar.PaddingLeft(4)
	OpenStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color(Cyan))
	SelectedOpenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(Blue))
	HelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	FocusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(Magenta))
	BlurredStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(White))
)
