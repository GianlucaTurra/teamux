package sessions

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	item struct {
		title string
		desc  string
		open  bool
	}
	SessionBrowserModel struct {
		list         list.Model
		selected     string
		openSessions string
		sessions     map[string]data.Session
		State        common.State
		db           *sql.DB
		logger       common.Logger
		help         sessionBrowserHelpModel
	}
	SessionDelegate struct{}
)

func (d SessionDelegate) Height() int                             { return 1 }
func (d SessionDelegate) Spacing() int                            { return 0 }
func (d SessionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d SessionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d. %s", index+1, i.title)
	fn := common.ItemStyle.Render
	if i.open {
		fn = func(s ...string) string { return common.OpenStyle.Render("* " + strings.Join(s, " ")) }
	}
	if index == m.Index() {
		if i.open {
			fn = func(s ...string) string {
				return common.SelectedOpenStyle.Render(">*" + strings.Join(s, " "))
			}
		} else {
			fn = func(s ...string) string {
				return common.SelectedStyle.Render("> " + strings.Join(s, " "))
			}
		}
	}
	fmt.Fprint(w, fn(str))
}

func (s item) FilterValue() string { return "" }

func NewSessionBrowserModel(db *sql.DB, logger common.Logger) SessionBrowserModel {
	sessions, layouts := loadData(db, logger)
	l := list.New(layouts, SessionDelegate{}, 100, 10)
	l.Title = "Available session layouts"
	l.Styles.Title = common.TitleStyle
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Styles.PaginationStyle = common.PaginationStyle
	l.Styles.HelpStyle = common.HelpStyle
	openSessions := data.CountTmuxSessions()
	return SessionBrowserModel{
		list:         l,
		openSessions: openSessions,
		sessions:     sessions,
		State:        common.Browsing,
		logger:       logger,
		db:           db,
		help:         newSessionBrowserHelpModel(),
	}
}

func loadData(db *sql.DB, logger common.Logger) (map[string]data.Session, []list.Item) {
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

func (m SessionBrowserModel) Init() tea.Cmd {
	return nil
}

func (m SessionBrowserModel) View() string {
	switch m.State {
	case common.Quitting:
		return "Bye, have a nice day!"
	case common.Deleting:
		return fmt.Sprintf("You are about to delete %s, press y to confirm", m.selected)
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.list.View(),
		fmt.Sprintf("Open sessions: %s", m.openSessions),
		m.help.View(),
	)
}

func (m SessionBrowserModel) Update(msg tea.Msg) (SessionBrowserModel, tea.Cmd) {
	switch msg := msg.(type) {
	case common.TmuxSessionsChanged:
		m.openSessions = data.CountTmuxSessions()
		return m, nil
	case common.OpenMsg:
		return m.openSelected()
	case common.SwitchMsg:
		return m.switchToSelected()
	case common.DeleteMsg:
		return m.deleteSelected()
	case common.KillMsg:
		return m.killSelected()
	case common.ReloadMsg:
		return NewSessionBrowserModel(m.db, m.logger), nil
	case tea.KeyMsg:
		if m.State == common.Deleting {
			switch msg.String() {
			case "y":
				m.State = common.Browsing
				return m, common.Delete
			default:
				m.State = common.Browsing
				return m, nil
			}
		}
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.State = common.Quitting
			return m, tea.Quit
		case "enter", " ":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selected = i.title
			}
			return m, func() tea.Msg { return common.OpenMsg{} }
		case "e":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.selected = i.title
			}
			return m, func() tea.Msg { return common.Edit(m.sessions[m.selected]) }
		case "s":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.selected = i.title
			}
			return m, common.Switch
		case "d":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selected = i.title
			}
			m.State = common.Deleting
			return m, nil
		case "K":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.selected = i.title
			}
			return m, common.Kill
		}
	}
	// handle sub-models updates
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	newHelp, cmd := m.help.Update(msg)
	m.help = newHelp
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// switchToSelected Switch to the selected session opening it if necessary
func (m SessionBrowserModel) switchToSelected() (SessionBrowserModel, tea.Cmd) {
	s := m.sessions[m.selected]
	if s.IsOpen() {
		if err := s.Switch(); err != nil {
			m.logger.Errorlogger.Printf("Error switching to session %s: %v", m.selected, err)
			return m, func() tea.Msg { return common.TmuxErr{} }
		}
	}
	if err := m.sessions[m.selected].Open(); err != nil {
		m.logger.Errorlogger.Printf("Error opening session %s: %v", m.selected, err)
		return m, func() tea.Msg { return common.TmuxErr{} }
	}
	if err := s.Switch(); err != nil {
		m.logger.Errorlogger.Printf("Error switching to session %s: %v", m.selected, err)
		return m, func() tea.Msg { return common.TmuxErr{} }
	}
	m.refreshItems()
	return m, func() tea.Msg { return common.TmuxSessionsChanged{} }
}

// openSelected Opens the selected session. If it is already open nothing is
// done.
func (m SessionBrowserModel) openSelected() (SessionBrowserModel, tea.Cmd) {
	if s := m.sessions[m.selected]; s.IsOpen() {
		// TODO: does it make sense to return nil?
		return m, func() tea.Msg { return nil }
	}
	s := m.sessions[m.selected]
	if err := s.Open(); err != nil {
		m.logger.Errorlogger.Printf("Error opening session %s: %v", s.Name, err)
		return m, func() tea.Msg { return common.TmuxErr{} }
	}
	m.refreshItems()
	return m, func() tea.Msg { return common.TmuxSessionsChanged{} }
}

// deleteSelected kills the session if open and proceeds to delete it from the db
func (m SessionBrowserModel) deleteSelected() (SessionBrowserModel, tea.Cmd) {
	m.killSelected()
	s := m.sessions[m.selected]
	if err := s.Delete(); err != nil {
		m.logger.Errorlogger.Printf("Error deleting session %s: %v", m.selected, err)
		return m, func() tea.Msg { return common.TmuxErr{} }
	}
	return m, func() tea.Msg { return common.ReloadMsg{} }
}

// killSelected kills the selected session. If it is not open nothing is done.
func (m SessionBrowserModel) killSelected() (SessionBrowserModel, tea.Cmd) {
	s := m.sessions[m.selected]
	if !s.IsOpen() {
		return m, nil
	}
	if err := s.Close(); err != nil {
		m.logger.Errorlogger.Printf("Error killing session %s: %v", m.selected, err)
		return m, func() tea.Msg { return common.TmuxErr{} }
	}
	m.refreshItems()
	return m, func() tea.Msg { return common.TmuxSessionsChanged{} }
}

// refreshItems checks again if any item status has changed
func (m *SessionBrowserModel) refreshItems() {
	var newList []list.Item
	for _, l := range m.list.Items() {
		i, ok := l.(item)
		if !ok {
			m.logger.Errorlogger.Printf("Failed to cast list item to item type: %v", l)
			continue
		}
		i.open = m.sessions[l.(item).title].IsOpen()
		newList = append(newList, i)
	}
	m.list.SetItems(newList)
}
