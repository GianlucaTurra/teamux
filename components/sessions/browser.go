package sessions

import (
	"errors"
	"fmt"

	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
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
		sessions     map[string]data.Session
		State        common.State
		connector    data.Connector
		logger       common.Logger
	}
)

func (s item) FilterValue() string { return "" }

func NewSessionBrowserModel(connector data.Connector, logger common.Logger) SessionBrowserModel {
	sessions, layouts := loadData(connector.DB, logger)
	l := list.New(layouts, sessionDelegate{}, 100, 10)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.Styles.PaginationStyle = common.PaginationStyle
	openSessions := data.CountTmuxSessions()
	return SessionBrowserModel{
		list:         l,
		openSessions: openSessions,
		sessions:     sessions,
		State:        common.Browsing,
		logger:       logger,
		connector:    connector,
	}
}

func loadData(db *gorm.DB, logger common.Logger) (map[string]data.Session, []list.Item) {
	layouts := []list.Item{}
	sessions, err := data.ReadAllSessions(db)
	if err != nil {
		logger.Fatallogger.Fatalf("Failed to read sessions: %v", err)
	}
	// TODO: should the pwd be checked?
	/* for i := range sessions {
		s := &sessions[i]
		if err := s.GetPWD(); err != nil {
			logger.Errorlogger.Printf("Error reading session %s working directory.\n%v", s.Name, err)
		}
	} */
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
		return NewSessionBrowserModel(m.connector, m.logger), nil
	case common.UpDownMsg:
		// TODO: refactor into a proper method for a clearer switch
		i, ok := m.list.SelectedItem().(item)
		if ok {
			m.selected = i.title
		}
		return m, func() tea.Msg { return common.NewSFocus{Session: m.sessions[m.selected]} }
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
			// TODO: refactor into a proper method for a clearer switch
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selected = i.title
			}
			return m, func() tea.Msg { return common.OpenMsg{} }
		case "e":
			// TODO: refactor into a proper method for a clearer switch
			if i, ok := m.list.SelectedItem().(item); ok {
				m.selected = i.title
			}
			return m, tea.Batch(
				func() tea.Msg { return common.ShowFullHelpMsg{Component: common.SessionEditor} },
				func() tea.Msg { return common.EditS(m.sessions[m.selected]) },
			)
		case "s":
			// TODO: refactor into a proper method for a clearer switch
			if i, ok := m.list.SelectedItem().(item); ok {
				m.selected = i.title
			}
			return m, common.Switch
		case "d":
			// TODO: refactor into a proper method for a clearer switch
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.selected = i.title
			}
			m.State = common.Deleting
			return m, nil
		case "K":
			// TODO: refactor into a proper method for a clearer switch
			if i, ok := m.list.SelectedItem().(item); ok {
				m.selected = i.title
			}
			return m, common.Kill
		case "w":
			// TODO: refactor into a proper method for a clearer switch
			if i, ok := m.list.SelectedItem().(item); ok {
				session := m.sessions[i.title]
				return m, func() tea.Msg { return common.AssociateWindows{Session: session} }
			}
		case "n":
			return m, common.NewSession
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

// switchToSelected Switch to the selected session opening it if necessary
func (m SessionBrowserModel) switchToSelected() (SessionBrowserModel, tea.Cmd) {
	s := m.sessions[m.selected]
	if s.IsOpen() {
		if err := s.Switch(); err != nil {
			m.logger.Errorlogger.Printf("Error switching to session %s: %v", m.selected, err)
			return m, func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
		}
	}
	if err := m.sessions[m.selected].Open(); err != nil {
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
	return m, func() tea.Msg { return common.TmuxSessionsChanged{} }
}

// openSelected Opens the selected session. If it is already open nothing is
// done.
func (m SessionBrowserModel) openSelected() (SessionBrowserModel, tea.Cmd) {
	if s := m.sessions[m.selected]; s.IsOpen() {
		msg := fmt.Sprintf("session %s already open", s.Name)
		return m, func() tea.Msg {
			return common.OutputMsg{Err: errors.New(msg), Severity: common.Info}
		}
	}
	s := m.sessions[m.selected]
	if err := s.Open(); err != nil {
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
	m.refreshItems()
	return m, func() tea.Msg { return common.TmuxSessionsChanged{} }
}

// deleteSelected kills the session if open and proceeds to delete it from the db
func (m SessionBrowserModel) deleteSelected() (SessionBrowserModel, tea.Cmd) {
	m.killSelected()
	s := m.sessions[m.selected]
	if _, err := s.Delete(m.connector); err != nil {
		m.logger.Errorlogger.Printf("Error deleting session %s: %v", m.selected, err)
		return m, func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
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
		return m, func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
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
