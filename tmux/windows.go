package tmux

import (
	"fmt"
	"strings"
)

// NewWindow creates a new window in the current session
func NewWindow(name string, workingDirectory string, target *string) error {
	var baseCmd string
	if target == nil || (strings.TrimSpace(*target) == "") {
		baseCmd = fmt.Sprintf("tmux new-window -d -n \"%s\"", name)
	} else {
		baseCmd = fmt.Sprintf("tmux neww -t %s -d -n \"%s\"", *target, name)
	}
	return commandWithWorkDir(workingDirectory, baseCmd)
}

// KillWindow kills the given window.
// To close a specific window in a given session the name should be in the form
// of `SessionName:WindowName`
func KillWindow(name string) error {
	return executeCommand(fmt.Sprintf("tmux kill-window -t \"%s\"", name))
}

func ReorderWindows(target string) error {
	return executeCommand(fmt.Sprintf("tmux movew -r -t \"%s\"", target))
}
