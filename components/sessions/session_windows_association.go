package sessions

import (
	"github.com/GianlucaTurra/teamux/common"
	"github.com/GianlucaTurra/teamux/components/data"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type availableWindows struct {
	window   data.Window
	title    string
	desc     string
	selected bool
}

func (aw availableWindows) FilterValue() string {
	return ""
}

type SessionWindowsAssociationModel struct {
	model     list.Model
	connector data.Connector
	logger    common.Logger
	session   data.Session
	state     common.State
}

func NewSessionWindowsAssociationModel(
	connector data.Connector,
	logger common.Logger,
	session data.Session,
) SessionWindowsAssociationModel {
	var aws []list.Item
	l := list.New(aws, sessionWindowsDelegate{}, 100, 10)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(true)
	l.Styles.Title = common.HeaderStyle
	l.Title = "Asscoaitate windows to: " + session.Name
	l.Styles.PaginationStyle = common.PaginationStyle
	return SessionWindowsAssociationModel{
		model:     l,
		connector: connector,
		logger:    logger,
		session:   session,
	}
}

func (m SessionWindowsAssociationModel) Init() tea.Cmd {
	return nil
}

func (m SessionWindowsAssociationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.LoadDataMsg:
		return m, m.loadWindowsList()
	case common.UpdateDetailMsg:
		return m, func() tea.Msg { return common.NewSFocus{Session: m.session} }
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			return m, m.selectWindow()
		case "esc":
			return m, common.Browse
		case "q", "ctrl+c":
			m.state = common.Quitting
			return m, common.Quit
		}
	}
	var cmd tea.Cmd
	newModel, cmd := m.model.Update(msg)
	m.model = newModel
	return m, cmd
}

func (m SessionWindowsAssociationModel) View() string {
	if m.state == common.Quitting {
		return ""
	}
	return m.model.View()
}

// loadWindowsList() loads all windows and marks the ones already associated
func (m *SessionWindowsAssociationModel) loadWindowsList() tea.Cmd {
	var aws []list.Item
	windows, err := data.ReadAllWindows(m.connector.DB)
	if err != nil {
		m.logger.Errorlogger.Printf("Error reading windows for session_windows association: %v", err)
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	for _, w := range windows {
		notAssociated := true
		for _, aw := range m.session.Windows {
			if w.ID == aw.ID {
				notAssociated = false
				break
			}
		}
		if notAssociated {
			aws = append(aws, availableWindows{window: w, title: w.Name, selected: false})
		} else {
			aws = append(aws, availableWindows{window: w, title: w.Name, selected: true})
		}
	}
	m.model.SetItems(aws)
	return nil
}

// selectWindow() appends or removes the selected window from the many to many
// relationship
func (m *SessionWindowsAssociationModel) selectWindow() tea.Cmd {
	w := m.model.SelectedItem().(availableWindows)
	var err error
	if w.selected {
		err = m.connector.DB.Model(&m.session).Association("Windows").Delete(&w.window)
	} else {
		err = m.connector.DB.Model(&m.session).Association("Windows").Append(&w.window)
	}
	if err != nil {
		m.logger.Errorlogger.Printf("Error appending %s to session %s: %v", w.window.Name, m.session.Name, err)
		return func() tea.Msg { return common.OutputMsg{Err: err, Severity: common.Error} }
	}
	itemIndex := m.model.GlobalIndex()
	w.selected = !w.selected
	m.model.SetItem(itemIndex, w)
	return common.UpdateDetail
}
