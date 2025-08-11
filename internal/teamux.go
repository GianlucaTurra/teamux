package internal

import (
	"bytes"
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	TmuxSessionsChanged struct{}
	TmuxErr             struct{}
	OpenMsg             struct{}
	DeleteMsg           struct{}
	SwitchMsg           struct{}
	NewMsg              struct{}
	KillMsg             struct{}
	InputErrMsg         struct{ Err error }
	ReloadMsg           struct{}
	SessionCreatedMsg   struct{}
)

func Open() tea.Msg    { return OpenMsg{} }
func Delete() tea.Msg  { return DeleteMsg{} }
func Switch() tea.Msg  { return SwitchMsg{} }
func New() tea.Msg     { return NewMsg{} }
func Kill() tea.Msg    { return KillMsg{} }
func Created() tea.Msg { return SessionCreatedMsg{} }
func Reaload() tea.Msg { return ReloadMsg{} }

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
