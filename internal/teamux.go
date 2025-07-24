package internal

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

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

func OpenTmuxSession(script string) tea.Msg {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	file := filepath.Join(home, script)
	cmd := exec.Command("/bin/sh", file)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		fmt.Println(err)
		// if err := cmd.Run(); err != nil {
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
