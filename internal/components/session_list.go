package components

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/GianlucaTurra/teamux/internal"
	"github.com/GianlucaTurra/teamux/internal/data"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle               = lipgloss.NewStyle().MarginLeft(2)
	sessionStyle             = lipgloss.NewStyle().PaddingLeft(4)
	selectedSessionStyle     = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("140"))
	paginationStyle          = list.DefaultStyles().TitleBar.PaddingLeft(4)
	openSessionStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).PaddingLeft(2)
	selectedOpenSessionStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("200"))
)

type (
	item struct {
		title string
		desc  string
		open  bool
	}
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
	str := fmt.Sprintf("%d. %s", index+1, i.title)
	fn := sessionStyle.Render
	if i.open {
		fn = func(s ...string) string { return openSessionStyle.Render("* " + strings.Join(s, " ")) }
	}
	if index == m.Index() {
		if i.open {
			fn = func(s ...string) string {
				return selectedOpenSessionStyle.Render(">*" + strings.Join(s, " "))
			}
		} else {
			fn = func(s ...string) string {
				return selectedSessionStyle.Render("> " + strings.Join(s, " "))
			}
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
	data, layouts := loadData(db, logger)
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

func loadData(db *sql.DB, logger internal.Logger) (map[string]data.Session, []list.Item) {
	layouts := []list.Item{}
	sessions, err := data.ReadAllSessions(db)
	if err != nil {
		logger.Fatallogger.Fatalf("Failed to read sessions: %v", err)
	}
	data := make(map[string]data.Session)
	for _, s := range sessions {
		layouts = append(layouts, item{title: s.Name, open: s.IsOpen()})
		data[s.Name] = s
	}
	return data, layouts
}

func (m sessionListModel) Init() tea.Cmd {
	return nil
}

func (m sessionListModel) View() string {
	switch m.state {
	case quitting:
		return "Bye, have a nice day!"
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.list.View(),
		fmt.Sprintf("Open sessions: %s", m.openSessions),
	)
}

func (m sessionListModel) Update(msg tea.Msg) (sessionListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case internal.TmuxSessionsChanged:
		m.openSessions = internal.CountTmuxSessions()
		return m, nil
	case internal.OpenMsg:
		return m.openSelected()
	case internal.SwitchMsg:
		return m.switchToSelected()
	case internal.DeleteMsg:
		return m.deleteSelected()
	case internal.KillMsg:
		return m.killSelected()
	case internal.ReloadMsg:
		return newSessionListModel(m.db, m.logger), nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.state = quitting
			return m, tea.Quit
		case "enter", " ":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selected = i.title
			}
			return m, func() tea.Msg { return internal.OpenMsg{} }
		case "s":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.selected = i.title
			}
			return m, internal.Switch
		case "d":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selected = i.title
			}
			return m, internal.Delete
		case "K":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.selected = i.title
			}
			return m, internal.Kill
		case "n":
			return m, internal.New
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// switchToSelected Switch to the selected session opening it if necessary
func (m sessionListModel) switchToSelected() (sessionListModel, tea.Cmd) {
	s := m.data[m.selected]
	if s.IsOpen() {
		if err := s.Switch(); err != nil {
			m.logger.Errorlogger.Printf("Error switching to session %s: %v", m.selected, err)
			return m, func() tea.Msg { return internal.TmuxErr{} }
		}
	}
	if err := m.data[m.selected].Open(); err != nil {
		m.logger.Errorlogger.Printf("Error opening session %s: %v", m.selected, err)
		return m, func() tea.Msg { return internal.TmuxErr{} }
	}
	if err := s.Switch(); err != nil {
		m.logger.Errorlogger.Printf("Error switching to session %s: %v", m.selected, err)
		return m, func() tea.Msg { return internal.TmuxErr{} }
	}
	m.refreshItems()
	return m, func() tea.Msg { return internal.TmuxSessionsChanged{} }
}

// openSelected Opens the selected session. If it is already open nothing is
// done.
func (m sessionListModel) openSelected() (sessionListModel, tea.Cmd) {
	if s := m.data[m.selected]; s.IsOpen() {
		// TODO: does it make sense to return nil?
		return m, func() tea.Msg { return nil }
	}
	s := m.data[m.selected]
	if err := s.Open(); err != nil {
		m.logger.Errorlogger.Printf("Error opening session %s: %v", s.Name, err)
		return m, func() tea.Msg { return internal.TmuxErr{} }
	}
	m.refreshItems()
	return m, func() tea.Msg { return internal.TmuxSessionsChanged{} }
}

// deleteSelected kills the session if open and proceeds to delete it from the db
func (m sessionListModel) deleteSelected() (sessionListModel, tea.Cmd) {
	m.killSelected()
	s := m.data[m.selected]
	if err := s.Delete(); err != nil {
		m.logger.Errorlogger.Printf("Error deleting session %s: %v", m.selected, err)
		return m, func() tea.Msg { return internal.TmuxErr{} }
	}
	return m, func() tea.Msg { return internal.ReloadMsg{} }
}

// killSelected kills the selected session. If it is not open nothing is done.
func (m sessionListModel) killSelected() (sessionListModel, tea.Cmd) {
	s := m.data[m.selected]
	if !s.IsOpen() {
		return m, nil
	}
	if err := s.Close(); err != nil {
		m.logger.Errorlogger.Printf("Error killing session %s: %v", m.selected, err)
		return m, func() tea.Msg { return internal.TmuxErr{} }
	}
	m.refreshItems()
	return m, func() tea.Msg { return internal.TmuxSessionsChanged{} }
}

// refreshItems checks again if any item status has changed
func (m *sessionListModel) refreshItems() {
	var newList []list.Item
	for _, l := range m.list.Items() {
		i, ok := l.(item)
		if !ok {
			m.logger.Errorlogger.Printf("Failed to cast list item to item type: %v", l)
			continue
		}
		i.open = m.data[l.(item).title].IsOpen()
		newList = append(newList, i)
	}
	m.list.SetItems(newList)
}
