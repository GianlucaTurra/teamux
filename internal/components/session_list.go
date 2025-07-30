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
	"github.com/charmbracelet/lipgloss"
)

type (
	Session          string
	SessionState     int
	sessionListModel struct {
		list         list.Model
		selected     string
		openSessions string
		data         map[string]db.SessionInfo
		state        SessionState
	}
)

const (
	browsing SessionState = iota
	creating
	editing
	quitting
)

func (s Session) FilterValue() string { return "" }

func newSessionListModel() sessionListModel {
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
	l.SetShowHelp(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	openSessions := internal.CountTmuxSessions()
	return sessionListModel{list: l, openSessions: openSessions, data: sessionsInfo, state: browsing}
}

func (m sessionListModel) Init() tea.Cmd {
	return nil
}

func (m sessionListModel) View() string {
	switch m.state {
	case quitting:
		return "See ya!"
	}
	return lipgloss.JoinVertical(lipgloss.Top, m.list.View(), fmt.Sprintf("Open sessions: %s", m.openSessions))
}

func (m sessionListModel) Update(msg tea.Msg) (sessionListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case internal.TmuxSessionsChanged:
		m.openSessions = internal.CountTmuxSessions()
		return m, nil
	case internal.OpenMsg:
		selectedInfo := m.data[m.selected]
		if selectedInfo.IsOpen {
			return m, func() tea.Msg { return internal.SwitchTmuxSession(m.selected) }
		}
		return m, func() tea.Msg {
			mess := internal.OpenTmuxSession(selectedInfo.File)
			m.data[m.selected] = db.SessionInfo{File: selectedInfo.File, IsOpen: true}
			return mess
		}
	case internal.SwitchMsg:
		selectedInfo := m.data[m.selected]
		var mess tea.Msg = nil
		if !selectedInfo.IsOpen {
			mess = internal.OpenTmuxSession(selectedInfo.File)
			m.data[m.selected] = db.SessionInfo{File: selectedInfo.File, IsOpen: true}
		}
		internal.SwitchTmuxSession(m.selected)
		return m, func() tea.Msg { return mess }
	case internal.DeleteMsg:
		if m.data[m.selected].IsOpen {
			return m, func() tea.Msg {
				mess := internal.KillTmuxSession(m.selected)
				m.data[m.selected] = db.SessionInfo{File: m.data[m.selected].File, IsOpen: false}
				return mess
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.state = quitting
			return m, tea.Quit
		case "enter", " ":
			i, ok := m.list.SelectedItem().(Session)
			if ok {
				m.selected = string(i)
			}
			return m, internal.Open
		case "s":
			if i, ok := m.list.SelectedItem().(Session); ok {
				m.selected = string(i)
			}
			return m, internal.Switch
		case "d":
			i, ok := m.list.SelectedItem().(Session)
			if ok {
				m.selected = string(i)
			}
			return m, internal.Delete
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
