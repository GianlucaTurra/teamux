package main

import (
	"fmt"
	"os"

	"github.com/GianlucaTurra/teamux/internal/components"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	const defaultWidth = 20
	p := tea.NewProgram(components.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}
