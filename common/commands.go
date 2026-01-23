package common

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Open() tea.Msg                  { return OpenMsg{} }
func Delete() tea.Msg                { return DeleteMsg{} }
func Switch() tea.Msg                { return SwitchMsg{} }
func NewWindow() tea.Msg             { return NewWindowMsg{} }
func NewPane() tea.Msg               { return NewPaneMsg{} }
func Kill() tea.Msg                  { return KillMsg{} }
func Reaload() tea.Msg               { return ReloadMsg{} }
func Browse() tea.Msg                { return BrowseMsg{} }
func Quit() tea.Msg                  { return QuitMsg{} }
func NextTab() tea.Msg               { return NextTabMsg{} }
func PreviousTab() tea.Msg           { return PreviousTabMsg{} }
func UpDown() tea.Msg                { return UpDownMsg{} }
func ClearHelp(t FocusedTab) tea.Msg { return ClearHelpMsg{t} }
func LoadData() tea.Msg              { return LoadDataMsg{} }
func UpdateDetail() tea.Msg          { return UpdateDetailMsg{} }
