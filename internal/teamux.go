// Package internal
// Common things for all components
package internal

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/GianlucaTurra/teamux/internal/data"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	TmuxSessionsChanged struct{}
	TmuxErr             struct{}
	OpenMsg             struct{}
	DeleteMsg           struct{}
	SwitchMsg           struct{}
	NewSessionMsg       struct{}
	NewWindowMsg        struct{}
	KillMsg             struct{}
	InputErrMsg         struct{ Err error }
	ReloadMsg           struct{}
	SessionCreatedMsg   struct{}
	BrowseMsg           struct{}
	EditMsg             struct{ Session data.Session }
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

func CountTmuxSessions() string {
	cmd := exec.Command("sh", "-c", "tmux ls | wc -l")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error counting sessions %v", err)
		return "Error executing tmux ls"
	}
	return out.String()
}
