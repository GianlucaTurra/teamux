package tmux

import (
	"fmt"
	"os/exec"
	"strconv"
)

/* PANES */

func CreatePane(target *string, splitRatio int, workingDirectory string, horizontal bool) error {
	var tmuxCommand string
	if target != nil {
		tmuxCommand = fmt.Sprintf(
			"tmux split-window -t \"%s\" -l %s -c \"%s\"",
			*target,
			strconv.Itoa(splitRatio)+"%",
			workingDirectory,
		)
	} else {
		tmuxCommand = fmt.Sprintf(
			"tmux split-window -l %s -c \"%s\"",
			strconv.Itoa(splitRatio)+"%",
			workingDirectory,
		)
	}
	if horizontal {
		tmuxCommand += " -h"
	} else {
		tmuxCommand += " -v"
	}
	cmd := exec.Command("sh", "-c", tmuxCommand)
	return cmd.Run()
}
