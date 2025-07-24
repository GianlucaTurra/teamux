package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GianlucaTurra/teamux/internal"
	"github.com/GianlucaTurra/teamux/internal/db"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	Session string
	Model   struct {
		list         list.Model
		selected     string
		openSessions string
		quitting     bool
		data         map[string]string
	}
)

func (s Session) FilterValue() string { return "" }

func InitialModel() Model {
	layouts := []list.Item{}
	sessionsInfo := db.ReadSeassions()
	for name := range sessionsInfo {
		layouts = append(layouts, Session(name))
	}
	l := list.New(layouts, SessionDelegate{}, 100, 10)
	l.Title = "Available session layouts"
	l.Styles.Title = titleStyle
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	openSessions := internal.CountTmuxSessions()
	return Model{list: l, openSessions: openSessions, quitting: false, data: sessionsInfo}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.quitting {
		return "\n See ya!"
	}
	return "\n" + m.list.View() + "\n" + fmt.Sprintf("Open sessions: %s", m.openSessions)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case internal.TmuxSessionOpened:
		m.openSessions = internal.CountTmuxSessions()
		return m, nil
	case internal.SelectMsg:
		return m, func() tea.Msg {
			return internal.OpenTmuxSession(m.data[m.selected])
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter", " ":
			i, ok := m.list.SelectedItem().(Session)
			if ok {
				m.selected = string(i)
			}
			return m, internal.Select
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
