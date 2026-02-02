// Package tmux contains the tools to interact with the tmux server
package tmux

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type NumberOfSessionsMsg struct{ Number string }

func NewSession(name string, workingDirectory string, detached bool) error {
	baseCmd := fmt.Sprintf("tmux new-session -d -s \"%s\"", name)
	err := commandWithWorkDir(workingDirectory, baseCmd, "")
	if err != nil {
		return err
	}
	if !detached {
		return SwitchToSession(name)
	}
	return nil
}

func HasSession(name string) bool {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tmux has-session -t \"%s\"", name))
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func KillSession(name string) error {
	return executeCommand(fmt.Sprintf("tmux kill-session -t \"%s\"", name))
}

func SwitchToSession(name string) error {
	return executeCommand(fmt.Sprintf("tmux switch -t \"%s\"", name))
}

// GetCurrentTmuxSessionName name of the session in which the app is launched
// should have no reason to fail unless tmux server is not running but that is
// checked at the top level
// TODO: display this in the browser
func GetCurrentTmuxSessionName() string {
	cmd := exec.Command("sh", "-c", "tmux display-message -p \"#S\"")
	out, _ := cmd.Output()
	return strings.TrimSpace(string(out))
}

func CountTmuxSessions() tea.Cmd {
	cmd := exec.Command("sh", "-c", "tmux ls | wc -l")
	var out bytes.Buffer
	cmd.Stdout = &out
	var ret string
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error counting sessions %v", err)
		ret = "Error executing tmux ls"
	} else {
		ret = out.String()
	}
	return func() tea.Msg { return NumberOfSessionsMsg{Number: ret} }
}
