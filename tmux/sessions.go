// Package tmux contains the tools to interact with the tmux server
package tmux

import (
	"fmt"
	"os/exec"
	"strings"
)

func NewSession(name string, workingDirectory string) error {
	baseCmd := fmt.Sprintf("tmux new-session -d -s \"%s\"", name)
	return commandWithWorkDir(workingDirectory, baseCmd)
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
