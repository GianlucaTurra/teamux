package sessions

import (
	"errors"
	"fmt"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/database"
	"github.com/GianlucaTurra/teamux/tmux"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
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
		sessions     map[string]Session
		State        common.State
		connector    database.Connector
		logger       common.Logger
	}
)

func (s item) FilterValue() string { return "" }

func NewSessionBrowserModel(connector database.Connector, logger common.Logger) SessionBrowserModel {
	sessions, layouts := loadData(connector.DB, logger)
	l := list.New(layouts, sessionDelegate{}, 100, 10)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Styles.PaginationStyle = common.PaginationStyle
	return SessionBrowserModel{
		list:         l,
		openSessions: "0",
		sessions:     sessions,
		State:        common.Browsing,
		logger:       logger,
		connector:    connector,
	}
}

func loadData(db *gorm.DB, logger common.Logger) (map[string]Session, []list.Item) {
	layouts := []list.Item{}
	// TODO: doesn't need to be public
	sessions, err := ReadAllSessions(db)
	if err != nil {
		logger.Fatallogger.Fatalf("Failed to read sessions: %v", err)
	}
	data := make(map[string]Session)
	for _, s := range sessions {
		layouts = append(layouts, item{title: s.Name, open: s.IsOpen()})
		data[s.Name] = s
	}
	return data, layouts
}

func (m SessionBrowserModel) Init() tea.Cmd {
	return tmux.CountTmuxSessions()
}

func (m SessionBrowserModel) View() string {
	switch m.State {
	case common.Quitting:
		return ""
	case common.Deleting:
		return fmt.Sprintf("You are about to delete %s, press y to confirm", m.selected)
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.list.View(),
		fmt.Sprintf("Open sessions: %s", m.openSessions),
	)
}

func (m SessionBrowserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case TmuxSessionsChanged:
		return m, tmux.CountTmuxSessions()
	case tmux.NumberOfSessionsMsg:
		m.openSessions = msg.Number
		return m, nil
	case common.OpenMsg:
		return m, func() tea.Msg { return openSelected(m.logger, m.sessions[m.selected]) }
	case common.SwitchMsg:
		return m.switchToSelected()
	case common.DeleteMsg:
		return m, deleteSelected(m.logger, m.sessions[m.selected], m.connector)
	case common.KillMsg:
		return m, killSelected(m.logger, m.sessions[m.selected])
	case common.ReloadMsg:
		return NewSessionBrowserModel(m.connector, m.logger), nil
	case common.UpDownMsg:
		// TODO: refactor into a proper method for a clearer switch
		return m.selectUpDownItem()
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
			return m, common.Quit
		case "enter", " ":
			return m.open()
		case "e":
			return m.editSelected()
		case "s":
			return m.switchToSession()
		case "d":
			return m.delete()
		// case "x":
		// 	return m.kill()
		case "a":
			return m.addWindowsToSession()
		case "n":
			return m, NewSession
		case "j", "k", "up", "down":
			cmds = append(cmds, common.UpDown)
		case "?":
			return m, func() tea.Msg { return common.ShowFullHelpMsg{Component: common.SessionBrowser} }
		}
	}
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m SessionBrowserModel) addWindowsToSession() (tea.Model, tea.Cmd) {
	if i, ok := m.list.SelectedItem().(item); ok {
		session := m.sessions[i.title]
		return m, func() tea.Msg { return AssociateWindowsMsg{Session: session} }
	} else {
		err := errors.New("no available windows")
		return m, func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Info} }
	}
}

func (m SessionBrowserModel) kill() (tea.Model, tea.Cmd) {
	if i, ok := m.list.SelectedItem().(item); ok {
		m.selected = i.title
	}
	return m, common.Kill
}

func (m SessionBrowserModel) delete() (tea.Model, tea.Cmd) {
	i, ok := m.list.SelectedItem().(item)
	if ok {
		m.selected = i.title
	}
	m.State = common.Deleting
	return m, nil
}

func (m SessionBrowserModel) switchToSession() (tea.Model, tea.Cmd) {
	if i, ok := m.list.SelectedItem().(item); ok {
		m.selected = i.title
	}
	return m, common.Switch
}

func (m SessionBrowserModel) editSelected() (tea.Model, tea.Cmd) {
	if i, ok := m.list.SelectedItem().(item); ok {
		m.selected = i.title
	}
	return m, func() tea.Msg { return EditS(m.sessions[m.selected]) }
}

func (m SessionBrowserModel) open() (tea.Model, tea.Cmd) {
	i, ok := m.list.SelectedItem().(item)
	if ok {
		m.selected = i.title
	}
	return m, func() tea.Msg { return common.OpenMsg{} }
}

func (m SessionBrowserModel) selectUpDownItem() (tea.Model, tea.Cmd) {
	i, ok := m.list.SelectedItem().(item)
	if ok {
		m.selected = i.title
	}
	return m, func() tea.Msg { return NewSFocus{Session: m.sessions[m.selected]} }
}

// switchToSelected Switch to the selected session opening it if necessary
func (m SessionBrowserModel) switchToSelected() (SessionBrowserModel, tea.Cmd) {
	s := m.sessions[m.selected]
	if s.IsOpen() {
		if err := s.Switch(); err != nil {
			m.logger.Errorlogger.Printf("Error switching to session %s: %v", m.selected, err)
			return m, func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
		}
	}
	if err := m.sessions[m.selected].Open(true); err != nil {
		// TODO: should be a common func to handle error types
		switch err.(type) {
		case tmux.Warning:
			m.logger.Warninglogger.Printf("Error opening session %s: %v", m.selected, err)
			return m, func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Warning} }
		default:
			m.logger.Errorlogger.Printf("Error opening session %s: %v", m.selected, err)
			return m, func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
		}
	}
	if err := s.Switch(); err != nil {
		m.logger.Errorlogger.Printf("Error switching to session %s: %v", m.selected, err)
		return m, func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	m.refreshItems()
	return m, func() tea.Msg { return TmuxSessionsChanged{} }
}

// openSelected Opens the selected session. If it is already open nothing is
// done.
func openSelected(logger common.Logger, s Session) tea.Cmd {
	if s.IsOpen() {
		msg := fmt.Sprintf("session %s already open", s.Name)
		return func() tea.Msg {
			return common.OutputMsg{Err: errors.New(msg), Severity: common.Info}
		}
	}
	if err := s.Open(true); err != nil {
		// TODO: should be a common func to handle error types
		switch err.(type) {
		case tmux.Warning:
			logger.Warninglogger.Printf("Error opening session %s: %v", s.Name, err)
			return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Warning} }
		default:
			logger.Errorlogger.Printf("Error opening session %s: %v", s.Name, err)
			return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
		}
	}
	return func() tea.Msg { return TmuxSessionsChanged{} }
}

// deleteSelected kills the session if open and proceeds to delete it from the db
func deleteSelected(logger common.Logger, s Session, connector database.Connector) tea.Cmd {
	if _, err := s.Delete(connector); err != nil {
		logger.Errorlogger.Printf("Error deleting session %s: %v", s.Name, err)
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	return func() tea.Msg { return common.ReloadMsg{} }
}

// killSelected kills the selected session. If it is not open nothing is done.
func killSelected(logger common.Logger, s Session) tea.Cmd {
	if !s.IsOpen() {
		return nil
	}
	if err := s.Close(); err != nil {
		logger.Errorlogger.Printf("Error killing session %s: %v", s.Name, err)
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	return func() tea.Msg { return TmuxSessionsChanged{} }
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
