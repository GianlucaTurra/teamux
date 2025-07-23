package components

import (
	"fmt"

	"github.com/GianlucaTurra/teamux/internal/layouts"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	sessionLayout []string
	cursor        int
	selected      map[int]struct{}
}

func InitialModel() Model {
	return Model{
		sessionLayout: layouts.ReadLayouts(),
		selected:      make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	s := "Available session layouts\n\n"
	for i, layout := range m.sessionLayout {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, layout)
	}
	s += "Press 'q' to quit\n"
	return s
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "j", "down":
			if m.cursor < len(m.sessionLayout)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}
