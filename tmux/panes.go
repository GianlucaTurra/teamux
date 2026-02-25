package tmux

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// SplitWindow creates a pane in the current window
func SplitWindow(splitRatio int, workingDirectory string, splitDirection string, shellCmd string) error {
	baseCmd := fmt.Sprintf("split-window -l %s%s -%s", strconv.Itoa(splitRatio), "%", splitDirection)
	return commandWithWorkDir(workingDirectory, baseCmd, "")
}

// SplitWindowWithTargetWindow creates a pane in the target window
// FIXME: if the window is in a detached session the target should be `Session:Window`
func SplitWindowWithTargetWindow(
	targetWindow string,
	splitRatio int,
	workingDirectory string,
	splitDirection string,
	shellCmd string,
) error {
	if strings.TrimSpace(targetWindow) == "" {
		return errors.New("missing target")
	}
	if splitRatio == 0 {
		return errors.New("missing splitRatio")
	}
	baseCmd := fmt.Sprintf(
		"split-window -t %s -l %s%% -%s",
		targetWindow,
		strconv.Itoa(splitRatio),
		splitDirection,
	)
	return commandWithWorkDir(workingDirectory, baseCmd, shellCmd)
}
