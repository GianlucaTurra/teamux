package internal

import (
	"bytes"
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	TmuxSessionOpened struct{}
	TmuxErr           struct{}
	SelectMsg         struct{}
)

func Select() tea.Msg {
	return SelectMsg{}
}

func OpenTmuxSession(selectedLayout string) tea.Msg {
	cmd := exec.Command("/bin/sh", selectedLayout)
	if err := cmd.Run(); err != nil {
		return TmuxErr{}
	}
	return TmuxSessionOpened{}
}

func CountTmuxSessions() string {
	cmd := exec.Command("sh", "-c", "tmux ls | wc -l")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error %v", err)
		return "Error executing tmux ls"
	}
	return out.String()
}
