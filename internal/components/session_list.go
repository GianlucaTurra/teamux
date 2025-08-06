package components

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/GianlucaTurra/teamux/internal"
	"github.com/GianlucaTurra/teamux/internal/data"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle           = lipgloss.NewStyle().MarginLeft(2)
	sessionStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedSessionStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("140"))
	paginationStyle      = list.DefaultStyles().TitleBar.PaddingLeft(4)
)

type (
	item             string
	SessionState     int
	sessionListModel struct {
		list         list.Model
		selected     string
		openSessions string
		data         map[string]data.Session
		state        SessionState
		db           *sql.DB
		logger       internal.Logger
	}
)

type SessionDelegate struct{}

func (d SessionDelegate) Height() int                             { return 1 }
func (d SessionDelegate) Spacing() int                            { return 0 }
func (d SessionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d SessionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, i)
	fn := sessionStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedSessionStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}

const (
	browsing SessionState = iota
	creating
	editing
	quitting
)

func (s item) FilterValue() string { return "" }

func newSessionListModel(db *sql.DB, logger internal.Logger) sessionListModel {
	layouts := []list.Item{}
	sessions, err := data.ReadAllSessions(db)
	if err != nil {
		logger.Fatallogger.Fatalf("Failed to read sessions: %v", err)
	}
	data := make(map[string]data.Session)
	for _, s := range sessions {
		layouts = append(layouts, item(s.Name))
		data[s.Name] = s
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
	return sessionListModel{list: l, openSessions: openSessions, data: data, state: browsing, logger: logger, db: db}
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
		if s := m.data[m.selected]; s.IsOpen() {
			return m, func() tea.Msg { return internal.SwitchTmuxSession(m.selected) }
		}
		return m, func() tea.Msg {
			s := m.data[m.selected]
			if err := s.Open(); err != nil {
				m.logger.Errorlogger.Printf("Error opening session %s: %v", s.Name, err)
				return internal.TmuxErr{}
			}
			return internal.OpenMsg{}
		}
	case internal.SwitchMsg:
		if s := m.data[m.selected]; s.IsOpen() {
			if err := internal.SwitchTmuxSession(m.selected); err != nil {
				m.logger.Errorlogger.Printf("Error switching to session %s: %v", m.selected, err)
				return m, func() tea.Msg { return internal.TmuxErr{} }
			}
		}
		if err := m.data[m.selected].Open(); err != nil {
			m.logger.Errorlogger.Printf("Error opening session %s: %v", m.selected, err)
			return m, func() tea.Msg { return internal.TmuxErr{} }
		}
		return m, func() tea.Msg { return internal.TmuxSessionsChanged{} }
	case internal.DeleteMsg:
		if s := m.data[m.selected]; s.IsOpen() {
			if err := s.Delete(); err != nil {
				m.logger.Errorlogger.Printf("Error deleting session %s: %v", m.selected, err)
				return m, func() tea.Msg { return internal.TmuxErr{} }
			}
			return m, func() tea.Msg { return internal.TmuxSessionsChanged{} }
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.state = quitting
			return m, tea.Quit
		case "enter", " ":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selected = string(i)
			}
			return m, func() tea.Msg { return internal.OpenMsg{} }
		case "s":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.selected = string(i)
			}
			return m, internal.Switch
		case "d":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selected = string(i)
			}
			return m, internal.Delete
		case "n":
			return m, internal.New
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
