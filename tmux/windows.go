package tmux

import (
	"errors"
	"fmt"
	"strings"
)

// NewWindow creates a new window in the current session
func NewWindow(name string, workingDirectory string) error {
	baseCmd := fmt.Sprintf("tmux new-window -d -n \"%s\"", name)
	return commandWithWorkDir(workingDirectory, baseCmd)
}

// NewWindowWithTarget creates a new window in the given target session
func NewWindowWithTarget(name string, workingDirectory string, target string) error {
	if strings.TrimSpace(target) == "" {
		return errors.New("missing window tartget")
	}
	baseCmd := fmt.Sprintf("tmux neww -t %s -d -n \"%s\"", target, name)
	return commandWithWorkDir(workingDirectory, baseCmd)
}

// KillWindow FIXME: if the windows is in a detached session the name should be `Session:Window`
func KillWindow(name string) error {
	return executeCommand(fmt.Sprintf("tmux kill-window -t \"%s\"", name))
}
