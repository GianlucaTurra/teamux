package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	Session string
	Model   struct {
		list     list.Model
		selected string
	}
)

func (s Session) FilterValue() string { return "" }

func InitialModel() Model {
	layouts := []list.Item{}
	for _, layout := range ReadLayouts() {
		layouts = append(layouts, Session(layout))
	}
	l := list.New(layouts, SessionDelegate{}, 100, 10)
	l.Title = "Available session layouts"
	l.Styles.Title = titleStyle
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return Model{list: l}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return "\n" + m.list.View()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", " ":
			i, ok := m.list.SelectedItem().(Session)
			if !ok {
				m.selected = string(i)
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func ReadLayouts() []string {
	var layoutFiles []string
	f, err := os.Open("./")
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	files, err := f.ReadDir(0)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext != ".sh" {
			continue
		}
		if !strings.Contains(file.Name(), "teamux") {
			continue
		}
		layoutFiles = append(layoutFiles, file.Name())
	}
	return layoutFiles
}
