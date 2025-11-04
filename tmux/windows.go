package tmux

import (
	"errors"
	"fmt"
	"strings"
)

// NewWindow creates a new window in the current session
func NewWindow(name string, workingDirectory string) error {
	baseCmd := fmt.Sprintf("tmux new-window -d -n \"%s\"", name)
	return newWindow(workingDirectory, baseCmd)
}

// NewWindowWithTarget creates a new window in the given target session
func NewWindowWithTarget(name string, workingDirectory string, target string) error {
	if strings.TrimSpace(target) == "" {
		return errors.New("missing window tartget")
	}
	baseCmd := fmt.Sprintf("tmux neww -t %s -d -n \"%s\"", target, name)
	return newWindow(workingDirectory, baseCmd)
}

func newWindow(workingDirectory string, cmd string) error {
	var nsdErr error
	if err := checkDirectory(workingDirectory); err != nil {
		nsdErr = NoSuchDirectoryError{"working directory doesn't exist"}
	}
	if strings.TrimSpace(workingDirectory) != "" {
		cmd += fmt.Sprintf(" -c %s", workingDirectory)
	}
	if err := executeCommand(cmd); err != nil {
		return err
	} else {
		return nsdErr
	}
}

// KillWindow FIXME: if the windows is in a detached session the name should be `Session:Window`
func KillWindow(name string) error {
	return executeCommand(fmt.Sprintf("tmux kill-window -t \"%s\"", name))
}
