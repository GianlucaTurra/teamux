package tmux

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

/* SESSIONS */

func CreateSession(name string, workingDirectory string) error {
	newSessionCmd := fmt.Sprintf("tmux new-session -d -s \"%s\" -c %s", name, workingDirectory)
	cmd := exec.Command("sh", "-c", newSessionCmd)
	return cmd.Run()
}

func IsSessionOpen(name string) bool {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("tmux has-session -t \"%s\"", name))
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

/* WINDOWS */

func CreateWindow(name string, workingDirectory string) error {
	cmd := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("tmux neww -d -n \"%s\" -c %s", name, workingDirectory),
	)
	return cmd.Run()
}

func CreateWindowWithTarget(name string, workingDirectory string, target string) error {
	cmd := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("tmux neww -t %s -d -n \"%s\" -c %s", target, name, workingDirectory),
	)
	return cmd.Run()
}

func KillWindow(name string) error {
	cmd := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("tmux kill-window -t \"%s\"", name),
	)
	return cmd.Run()
}

/* PANES */

func CreatePane(target string, splitRatio int, workingDirectory string, horizontal bool) error {
	var tmuxCommand string
	if strings.TrimSpace(target) == "" {
		tmuxCommand = fmt.Sprintf(
			"tmux split-window -t \"%s\" -l %s -c \"%s\"",
			target,
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
