// Package tmux contains the tools to interact with the tmux server
package tmux

import (
	"fmt"
	"os/exec"
	"strings"
)

func NewSession(name string, workingDirectory string) error {
	var newSessionCmd string
	if strings.TrimSpace(workingDirectory) == "" {
		newSessionCmd = fmt.Sprintf("tmux new-session -d -s \"%s\"", name)
	} else {
		newSessionCmd = fmt.Sprintf("tmux new-session -d -s \"%s\" -c %s", name, workingDirectory)
	}
	cmd := exec.Command("sh", "-c", newSessionCmd)
	return cmd.Run()
}

func HasSession(name string) bool {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tmux has-session -t \"%s\"", name))
	// TODO: needs testing
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func KillSession(name string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tmux kill-session -t \"%s\"", name))
	return cmd.Run()
}

func SwitchToSession(name string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tmux switch -t \"%s\"", name))
	return cmd.Run()
}
