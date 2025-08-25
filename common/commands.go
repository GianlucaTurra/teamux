package common

import (
	"github.com/GianlucaTurra/teamux/components/data"
	tea "github.com/charmbracelet/bubbletea"
)

func Open() tea.Msg               { return OpenMsg{} }
func Delete() tea.Msg             { return DeleteMsg{} }
func Switch() tea.Msg             { return SwitchMsg{} }
func NewSession() tea.Msg         { return NewSessionMsg{} }
func NewWinsow() tea.Msg          { return NewWindowMsg{} }
func Kill() tea.Msg               { return KillMsg{} }
func Created() tea.Msg            { return SessionCreatedMsg{} }
func Reaload() tea.Msg            { return ReloadMsg{} }
func Browse() tea.Msg             { return BrowseMsg{} }
func Edit(s data.Session) tea.Msg { return EditMsg{s} }
