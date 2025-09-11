package data

import (
	"bytes"
	"fmt"
	"os/exec"
)

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
