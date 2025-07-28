package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	TmuxSessionsChanged struct{}
	TmuxErr             struct{}
	SelectMsg           struct{}
	DeleteMsg           struct{}
)

func Select() tea.Msg {
	return SelectMsg{}
}

func Delete() tea.Msg {
	return DeleteMsg{}
}

func OpenTmuxSession(script string) tea.Msg {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error checking home dir %v", err)
	}
	file := filepath.Join(home, script)
	cmd := exec.Command("/bin/sh", file)
	if _, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Error opening session %v", err)
		return TmuxErr{}
	}
	return TmuxSessionsChanged{}
}

func KillTmuxSession(sessionName string) tea.Msg {
	tmuxCmd := fmt.Sprintf("tmux kill-session -t \"%s\"", sessionName)
	cmd := exec.Command("sh", "-c", tmuxCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error deleting %s %v", sessionName, err)
	}
	return TmuxSessionsChanged{}
}

func CountTmuxSessions() string {
	cmd := exec.Command("sh", "-c", "tmux ls | wc -l")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error counting sessions %v", err)
		return "Error executing tmux ls"
	}
	return out.String()
}

func IsTmuxSessionOpen(sessionName string) bool {
	tmuxCmd := fmt.Sprintf("tmux ls | grep %s", sessionName)
	cmd := exec.Command("sh", "-c", tmuxCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error checking open sessions %v", err)
		return false
	}
	return len(out.String()) != 0
}

func SwitchTmuxSession(sessionName string) tea.Msg {
	tmuxCmd := fmt.Sprintf("tmux switch -t %s", sessionName)
	cmd := exec.Command("sh", "-c", tmuxCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error switching to sessions %s %v", sessionName, err)
	}
	return nil
}
